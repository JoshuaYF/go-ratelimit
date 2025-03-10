# go-ratelimit
A rate limiting tool based on Leaky Bucket and Token Bucket, implemented in Golang.

Requires Golang 1.18 or higher, with support for generics.

### Install:
```shell
go get github.com/JoshuaYF/go-ratelimit@latest
```

### Usage:

```go
var RateLimter = NewLeakyBucketRateLimter[int, int](5, 2)

func TestLeakyBucketRateLimter(t *testing.T) {
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
	err := RateLimter.TryRequest(&task)
	if err != nil {
		fmt.Println(fmt.Sprintf("i = %d, err = %s", i, err))
		return
	}
	res, _panic := task.GetResult()
	fmt.Println(res, _panic)
}
```
