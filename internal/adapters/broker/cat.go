package broker

import (
	"github.com/kvolis/tesgode/cat"
)

type CatBroker struct {
	client *cat.Cat
}

func New() *CatBroker {
	return &CatBroker{
		client: cat.New(),
	}
}

func (b *CatBroker) Connect(conn string) error {
	return b.client.Connect(conn)
}

func (b *CatBroker) Subscribe() (<-chan []byte, error) {
	catCh, err := b.client.Subscript()
	if err != nil {
		return nil, err
	}

	out := make(chan []byte)
	go func() {
		defer close(out)
		for msg := range catCh {
			out <- msg.Bytes()
		}
	}()

	return out, nil
}

func (b *CatBroker) Close() error {
	return b.client.Close()
}
