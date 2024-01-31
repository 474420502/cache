package cache

import (
	"reflect"
	"runtime"
	"sync/atomic"
)

type AsyncSymbiosis struct {
	stop chan bool
	once atomic.Bool
}

func NewAsyncSymbiosis() *AsyncSymbiosis {
	return &AsyncSymbiosis{
		stop: make(chan bool),
	}
}

func (as *AsyncSymbiosis) Stop() {
	as.stop <- true
}

func (as *AsyncSymbiosis) StartLoopWith(withTarget interface{}, loop func()) {

	wtV := reflect.ValueOf(withTarget).Type()
	if wtV.Kind() == reflect.Ptr {
		wtV = wtV.Elem()
	}
	if wtV.NumField() == 0 {
		panic("StartGoWith Struct must having Field")
	}

	if as.once.CompareAndSwap(false, true) {
		// log.Println("CompareAndSwap")
		runtime.SetFinalizer(withTarget, func(obj interface{}) {
			// log.Println("SetFinalizer")
			as.Stop()
		})

		go func() {
			defer func() {
				as.once.Store(false)
			}()

			for {
				select {
				case <-as.stop:
					// log.Println("stop")
					return
				default:
					loop()
				}
			}
		}()
	} else {
		panic("AsyncSymbiosis StartWith only allowed call one time")
	}
}
