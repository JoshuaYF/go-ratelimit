package main

type Invoker[Req any, Res any] func(Req) Res

type Task[Req any, Res any] struct {
	Invoker   Invoker[Req, Res]
	Request   Req
	ResChan   chan Res
	PanicChan chan any
}

func (t *Task[Req, Res]) GetResult() (result Res, _panic any) {
	select {
	case result = <-t.ResChan:
	case _panic = <-t.PanicChan:
	}

	return
}
