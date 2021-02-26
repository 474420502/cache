package cache

import (
	"fmt"
	"log"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/474420502/gcurl"
)

func init() {
	log.Println("docker run -p 80:80 kennethreitz/httpbin. and change /etc/host 127.0.0.1 httpbin.org")
}

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
		time.Millisecond * 500,
	)

	if old != cache.Value() {
		t.Error("cache Value error")
	}

	for i := 0; i < 2; i++ {
		time.Sleep(time.Millisecond * 700) //因为更细是异步虽然触发了更新. 异步更新不算时间
		n := cache.Value()
		if old == n {
			t.Error("value should be updated", n, old)
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

func TestBlockWithCond(t *testing.T) {

	cache := New(func() interface{} {
		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
		if err != nil {
			log.Println(err)
		}
		return string(resp.Content())
	})

	cache.SetBlock(true)
	cache.SetUpdateInterval(time.Millisecond * 500)
	old := cache.Value()

	for i := 0; i < 2; i++ {
		time.Sleep(time.Millisecond * 500)
		if old == cache.Value() {
			log.Println("old not equal new value")
		}
	}

}

func TestError(t *testing.T) {
	cache := New(func() interface{} {
		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
		if err != nil {
			log.Println(err)
		}

		panic("error")

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

func TestUpdateReturnError(t *testing.T) {
	cache := New(func() interface{} {
		return fmt.Errorf("test errro")
	})

	var i = 0
	cache.SetOnError(func(err interface{}) {
		i++
		log.Println(i)
	})

	cache.Update() //第一次更新错误
	cache.SetBlock(true)
	cache.Value() // 第二次都认为需要更新

	if i != 2 {
		t.Error("SetOnError")
	}
}

func TestForce(t *testing.T) {
	wg := &sync.WaitGroup{}
	cache := New(func() interface{} {
		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
		if err != nil {
			log.Println(err)
		}
		return string(resp.Content())
	})

	cache.SetUpdateInterval(time.Millisecond * 100)

	for i := 0; i < 1000; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			var old interface{} = cache.Value()
			var count = 0
			now := time.Now()
			// log.Println(now)
			for i := 0; i < 1000; i++ {
				nvalue := cache.Value()
				if old != nvalue {
					old = nvalue
					count++
					// log.Println(old, count)
				}
				time.Sleep(time.Millisecond * 1)
			}
			predict := int(math.Round(float64(time.Since(now).Milliseconds()))) / 100
			if !(predict >= count-1 || predict <= count+1) {
				t.Error("error", predict, count)
			}
		}()
	}

	wg.Wait()
}

func TestCaseErr(t *testing.T) {
	var i = 0
	cache := New(func() interface{} {
		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
		if err != nil {
			log.Println(err)
		}

		i++
		if i%2 == 0 {
			return nil
		}

		// log.Println(string(resp.Content()))
		return string(resp.Content())
	})

	cache.SetBlock(true)

	for n := 0; n < 4; n++ {
		cache.Update()
		if v, ok := cache.Value().(string); !ok {
			t.Error("value error", v)
		}
	}

}

func TestCasePanic(t *testing.T) {
	var i = 0
	cache := New(func() interface{} {
		if i == 2 {
			panic("error xixi")
		}
		i++

		// log.Println(string(resp.Content()))
		return "nono"
	})

	cache.SetOnError(func(err interface{}) {
		if err.(string) != "error xixi" {
			t.Error("panic test is error")
		}
	})

	cache.SetBlock(true)

	for n := 0; n < 4; n++ {
		cache.Update()
		if v, ok := cache.Value().(string); !ok {
			t.Error("value error", v)
		}
	}

}
