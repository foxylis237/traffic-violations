package service

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kvolis/tesgode/internal/config"
	"github.com/kvolis/tesgode/internal/ports"
	"github.com/kvolis/tesgode/internal/retry"
	"github.com/kvolis/tesgode/models"
)

type Violation struct {
	LicenseNum string `json:"licenseNum"`
	UnixTime   int    `json:"unixTime"`
}

type Service struct {
	broker    ports.Broker
	storage   ports.Storage
	cfg       config.Config
	wg        sync.WaitGroup
	processed atomic.Int64
}

func New(broker ports.Broker, storage ports.Storage, cfg config.Config) *Service {
	return &Service{
		broker:  broker,
		storage: storage,
		cfg:     cfg,
	}
}

func (s *Service) Run(ctx context.Context) error {
	msgCh, err := s.broker.Subscribe()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for range s.cfg.Workers {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			for {
				select {
				case msg, ok := <-msgCh:
					if !ok {
						return
					}
					s.process(msg)
					if s.processed.Add(1) >= int64(s.cfg.Limit) {
						cancel()
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	<-ctx.Done()
	log.Println("shutting down...")
	s.broker.Close()
	s.wg.Wait()
	log.Printf("done, processed %d messages", s.processed.Load())
	return nil
}

func (s *Service) process(raw []byte) {
	var p models.Passage
	if err := json.Unmarshal(raw, &p); err != nil {
		log.Printf("parse error: %v", err)
		return
	}

	if len(p.Track) == 0 {
		return
	}

	last := p.Track[0]
	for _, pt := range p.Track[1:] {
		if pt.T > last.T {
			last = pt
		}
	}

	if last.T%60 < 45 {
		if s.cfg.LogLevel == config.LogDebug {
			log.Printf("no violation: %s unix=%d", p.LicenseNum, last.T)
		}
		return
	}

	value, _ := json.Marshal(Violation{
		LicenseNum: p.LicenseNum,
		UnixTime:   last.T,
	})

	err := retry.Do(3, 100*time.Millisecond, func() error {
		_, err := s.storage.Save(p.LicenseNum, value)
		return err
	})
	if err != nil {
		log.Printf("save failed: %s: %v", p.LicenseNum, err)
		return
	}

	log.Printf("violation saved: %s unix=%d", p.LicenseNum, last.T)
}
