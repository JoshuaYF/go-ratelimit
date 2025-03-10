package ratelimit

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

func (t *Task[Req, Res]) asyncExec() {
	go func() {
		defer func() {
			r := recover()
			if r != nil {
				select {
				case t.PanicChan <- r:
				default:
					// no listener
				}
			}
		}()
		res := t.Invoker(t.Request)
		select {
		case t.ResChan <- res:
			// put into ResChan success
		default:
			// no listener
		}
	}()
}
