package utils

import "time"

func Retry(attempts int, sleep time.Duration, fn func() error) error {
	err := fn()
	if err == nil {
		return nil
	}

	for i := 0; i < attempts-1; i++ {
		time.Sleep(sleep)
		err = fn()
		if err == nil {
			return nil
		}
		sleep *= 2
	}

	return err
}
