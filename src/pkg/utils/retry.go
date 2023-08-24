package utils

import (
	"log"
	"math/rand"
	"time"

	"github.com/rotisserie/eris"
)

// Retry executes the provided function with Retry logic.
// in case of failures, the `fu` is retried `maxRetries` number of time,
// and wait for (`retryIntervalSec` + random `jitterLimitSec`) duration between each Retry.
func Retry(maxRetries int, retryIntervalSec, jitterLimitSec int, fu func() error) error {
	// Create a backoff function that adds jitter to the delay duration.
	backoff := func(attempt int) time.Duration {
		return (time.Second * time.Duration(retryIntervalSec)) + (time.Second * time.Duration(jitterLimitSec) * time.Duration(rand.Intn(100)) / 100)
	}

	var err error
	for i := 0; i <= maxRetries; i++ {
		err = fu()
		if err == nil {
			// Success, no need to retry
			return nil
		}

		if i == maxRetries {
			break // skip delay after last attempt
		}

		t := backoff(i)
		log.Printf("retrying in %s, err=%s", t, err)
		time.Sleep(t) // Sleep before retrying
	}

	return eris.Wrapf(err, "failed after %d retries", maxRetries)
}
