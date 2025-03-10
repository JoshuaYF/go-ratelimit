package ratelimit

import (
	"fmt"
	"testing"
	"time"
)

var TokenRateLimiter = NewTokenBucketRateLimiter[int, int](5, 2)

func TestTokenBucketRateLimiter(t *testing.T) {
	for i := 0; i < 10; i++ {
		go do0(i)
	}

	time.Sleep(time.Minute)
}

func do0(i int) {
	task := Task[int, int]{
		Invoker: square,
		Request: i,
	}
	err := TokenRateLimiter.TryRequest(&task)
	if err != nil {
		fmt.Println(fmt.Sprintf("i = %d, err = %s", i, err))
		return
	}
	res, _panic := task.GetResult()
	fmt.Println(res, _panic)
}
