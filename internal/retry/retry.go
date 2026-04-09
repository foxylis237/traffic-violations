package retry

import (
	"time"
)

func Do(attempts int, delay time.Duration, fn func() error) error {
	var err error
	for range attempts {
		err = fn()
		if err == nil {
			return nil
		}
		time.Sleep(delay)
	}
	return err
}
