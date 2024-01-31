package cache

import (
	"log"
	"testing"
	"time"
)

type ABC struct {
	X int
}

func TestCase(t *testing.T) {

	ac := NewAsyncCache[*ABC](func(ctx *AsyncCacheContext[*ABC]) (*ABC, error) {

		uv := 1
		if ctx.current != nil {
			uv = ctx.current.X + 1
		}

		abc := &ABC{X: uv}

		time.Sleep(time.Second)
		log.Println("ABC created")
		return abc, nil
	})

	if err := ac.WaitFirstUpdate(time.Second * 5); err != nil {
		t.Errorf(err.Error())
	}
	v := ac.Get()
	if v == nil {
		t.Errorf("v = nil?")
		return
	}

	if v.X != 1 {
		t.Errorf("v.X != 1? first update")
		return
	}

	// time.Sleep(time.Second * 2)
	for i := 0; i < 2; i++ {
		time.Sleep(time.Millisecond * 1100)
		ac.GetSync(func(current *ABC) {
			if current.X <= 1 {
				t.Error("current.X <= 2?")
			}
			log.Println(current)
		})
	}

}
