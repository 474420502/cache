package cache

import (
	"log"
	"sync"
	"testing"
	"time"

	"github.com/474420502/gcurl"
)

func TestBlockCase1(t *testing.T) {
	cache := NewBlockCache(time.Millisecond*50, func(share interface{}) interface{} {
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

func TestBlockCase2(t *testing.T) {
	cache := NewBlockCache(time.Millisecond*50, func(share interface{}) interface{} {
		if share == nil {
			// log.Println("share is nil", share)
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

	cache.Value()

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

// func TestBlockBlockWithCond(t *testing.T) {

// }

func TestBlockError(t *testing.T) {
	cache := NewBlockCache(time.Millisecond*50, func(share interface{}) interface{} {
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

func TestBlockMulti(t *testing.T) {

	defer func() {
		if err := recover(); err != nil {
			t.Error(err)
		}
	}()

	var i = 0

	cache := NewBlockCache(time.Millisecond*50, func(share interface{}) interface{} {
		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
		if err != nil {
			log.Println(err)
		}
		i++
		if i%10 == 0 {
			panic("error")
		}
		return string(resp.Content())
	})

	cache.SetOnUpdateError(func(err interface{}) {
		i = 1
	})
	wg := &sync.WaitGroup{}
	for n := 0; n < 50; n++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			for x := 0; x < 100; x++ {
				cache.Value()
				cache.GetUpdate()

				time.Sleep(time.Millisecond)
			}
		}(wg)
	}

	wg.Wait()
}

func TestBlockTime(t *testing.T) {

	defer func() {
		if err := recover(); err != nil {
			t.Error(err)
		}
	}()

	var i = 0

	cache := NewBlockCache(time.Millisecond*1, func(share interface{}) interface{} {
		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
		if err != nil {
			log.Println(err)
		}
		i++
		if i%10 == 0 {
			panic("error")
		}
		time.Sleep(time.Millisecond * 10)
		return string(resp.Content())
	})

	cache.SetOnUpdateError(func(err interface{}) {
		i = 1
	})

	now := time.Now()

	wg := &sync.WaitGroup{}
	for n := 0; n < 10; n++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			for x := 0; x < 20; x++ {
				cache.Value()
				cache.GetUpdate()
				time.Sleep(time.Millisecond)
			}
		}(wg)
	}

	wg.Wait()

	if time.Since(now) <= time.Millisecond*200 {
		t.Error("block time is error", time.Since(now))
	}

}
