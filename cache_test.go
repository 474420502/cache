package cache

import (
	"log"
	"testing"
	"time"

	"github.com/474420502/gcurl"
)

func init() {
	log.Println("docker run -p 80:80 kennethreitz/httpbin. and change /etc/host 127.0.0.1 httpbin.org")
}

func TestCase1(t *testing.T) {
	cache := New(time.Millisecond*50, func(share interface{}) interface{} {
		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
		if err != nil {
			log.Println(err)
		}

		return string(resp.Content())
	})
	defer cache.Destroy()

	old := cache.Value()
	for i := 0; i < 2; i++ {
		time.Sleep(time.Millisecond * 70) //因为更细是异步虽然触发了更新. 异步更新不算时间
		n := cache.Value()
		if old == n {
			t.Error("value should be updated", n, old)
		}
		old = cache.Value()
	}

}

func TestCase2(t *testing.T) {
	cache := New(time.Millisecond*50, func(share interface{}) interface{} {
		if share == nil {
			log.Println("share is nil", share)
			time.Sleep(time.Millisecond * 50)
			return nil
		}

		if share != 1 {
			t.Error("share is not 1")
		}

		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
		if err != nil {
			log.Println(err)
		}

		return string(resp.Content())
	})
	defer cache.Destroy()

	cache.SetShare(1)
	old := cache.Value()
	if old != nil {
		t.Error("old not nil")
	}
	time.Sleep(time.Millisecond * 100)
	if cache.Value() == nil {
		t.Error("uid get error")
	}
}

// func TestBlockWithCond(t *testing.T) {

// }

func TestError(t *testing.T) {
	cache := New(time.Millisecond*50, func(share interface{}) interface{} {
		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
		if err != nil {
			log.Println(err)
		}
		panic("error")
		return string(resp.Content())
	})

	var i = 0
	cache.SetOnUpdateError(func(err interface{}) {
		i = 1
	})

	time.Sleep(time.Millisecond * 100)

	cache.Value()
	if i != 1 {
		t.Error("onError is error")
	}
}
