package cache

import (
	"log"
	"testing"
	"time"

	"github.com/474420502/gcurl"
)

func TestCase1(t *testing.T) {
	cache := New(func() interface{} {
		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
		if err != nil {
			log.Println(err)
		}

		// log.Println(string(resp.Content()))
		return string(resp.Content())
	})

	old := cache.Value()
	if old == nil {
		t.Error("value is not nil")
	}

	cache.SetUpdateInterval(
		time.Second,
	)

	if old != cache.Value() {
		t.Error("cache Value error")
	}

	for i := 0; i < 2; i++ {
		time.Sleep(time.Second)
		if old == cache.Value() {
			t.Error("value should be updated")
		}
		old = cache.Value()
	}

}

func TestCase2(t *testing.T) {
	cache := New(func() interface{} {
		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
		if err != nil {
			log.Println(err)
		}
		return string(resp.Content())
	})
	var i = 0
	cache.SetUpdateCond(func() bool {
		return i == 2
	})

	old := cache.Value()
	for i = 0; i < 10; i++ {
		time.Sleep(time.Millisecond * 10 * time.Duration(i))
		n := cache.Value()
		if old != n {
			if i <= 2 {
				t.Error("cond is error", i, old, n)
			}
			break
		}

		if i == 9 {
			t.Error("cond is error")
		}
	}

	time.Sleep(time.Millisecond * 100)
}

func TestCaseBlock(t *testing.T) {
	var isFirst bool = false
	cache := New(func() interface{} {
		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
		if err != nil {
			log.Println(err)
		}

		if isFirst {
			time.Sleep(time.Second)
		}
		isFirst = true

		return string(resp.Content())
	})

	cache.SetBlock(true)
	old := cache.Value()
	now := time.Now()
	for i := 0; i < 2; i++ {
		cache.Update()
		if old == cache.Value() {
			t.Error("error")
		}
	}
	if time.Since(now) <= time.Second*2 {
		t.Error("block is error")
	}
}

func TestError(t *testing.T) {
	cache := New(func() interface{} {
		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
		if err != nil {
			log.Println(err)
		}

		panic("erro")

		return string(resp.Content())
	})
	var i = 0
	cache.SetOnError(func(err interface{}) {
		i = 1
	})

	cache.Value()

	if i != 1 {
		t.Error("onError is error")
	}
}
