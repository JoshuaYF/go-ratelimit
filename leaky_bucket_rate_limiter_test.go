package main

import (
	"fmt"
	"testing"
	"time"
)

var RateLimiter = NewLeakyBucketRateLimiter[int, int](5, 2)

func TestLeakyBucketRateLimiter(t *testing.T) {
	for i := 0; i < 10; i++ {
		go do(i)
	}

	time.Sleep(time.Minute)
}

func square(a int) int {
	if a == 4 {
		panic("panic test")
	}

	return a * a
}

func do(i int) {
	task := Task[int, int]{
		Invoker: square,
		Request: i,
	}
	err := RateLimiter.TryRequest(&task)
	if err != nil {
		fmt.Println(fmt.Sprintf("i = %d, err = %s", i, err))
		return
	}
	res, _panic := task.GetResult()
	fmt.Println(res, _panic)
}
