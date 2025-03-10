package main

import (
	"errors"
	"time"
)

type TokenBucketRateLimiter[Req any, Res any] struct {
	capacity int64         // The capacity of the bucket.
	rate     int64         // The rate at which tokens are added to the bucket per second.
	bucket   chan struct{} // Token Bucket.
	t        *time.Ticker  // The ticker generated by the rate variable.
}

func NewTokenBucketRateLimiter[Req any, Res any](capacity, rate int64) *TokenBucketRateLimiter[Req, Res] {
	limiter := &TokenBucketRateLimiter[Req, Res]{
		capacity: capacity,
		rate:     rate,
		bucket:   make(chan struct{}, capacity),
		t:        time.NewTicker(time.Duration(time.Second.Nanoseconds() / rate)),
	}

	for i := 0; i < int(capacity); i++ {
		limiter.bucket <- struct{}{}
	}

	go func() {
		for {
			select {
			case <-limiter.t.C:
				limiter.bucket <- struct{}{}
			}
		}
	}()

	return limiter
}

func (limiter *TokenBucketRateLimiter[Req, Res]) ResetRate(rate int64) {
	limiter.rate = rate
	limiter.t.Reset(time.Duration(time.Second.Nanoseconds() / rate))
}

// TryRequest Put the task into the leaky bucket. If it's full, discard it and return an error.
func (limiter *TokenBucketRateLimiter[Req, Res]) TryRequest(task *Task[Req, Res]) (err error) {
	task.ResChan = make(chan Res)
	task.PanicChan = make(chan any, 1)

	select {
	case <-limiter.bucket:
		go func() {
			defer func() {
				r := recover()
				if r != nil {
					task.PanicChan <- r
				}
			}()
			res := task.Invoker(task.Request)
			select {
			case task.ResChan <- res:
				// put into ResChan success
			default:
				// no listener
			}
		}()

		return nil
	default:

		return errors.New("rejected")
	}
}
