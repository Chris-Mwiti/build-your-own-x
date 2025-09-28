package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"
)

var ERR_CTX_DONE = errors.New("context done.")

//used to retry functionalities for a number of counts
func RetryFn[T any](parentCtx context.Context, count int, initialBackOff time.Duration, fn func(context.Context) (T, error)) (T, error){
	if count <= 0 {
		count = 1
	}

	ctx, cancel := context.WithTimeout(parentCtx, 10 * time.Second)
	defer cancel()

	var lastError error
	backoff := initialBackOff

	for i := 1; i <= count; i++{
		select {
		case <-ctx.Done():
			return *new(T), ctx.Err()
		default:
			//we will delay after the initial request due to cpu resource intensive nature
			if i > 1{
				//jittering effect
				jitter := time.Duration(rand.Int63n(int64(backoff)))
				log.Printf("Retrying attempt %d:%d after failure: sleeping %v", i, count, (backoff+jitter))
				time.Sleep(backoff + jitter)
				backoff *= 2
			}

			res, err := fn(ctx)	
			if err == nil {
				return res, nil 
			}
		 
			lastError = err
			log.Printf("Attempt failed %d:%d failed %v\n", i, count, err)
		}
	}

	return *new(T), fmt.Errorf("Retry failed: %v\n", lastError) 
}


