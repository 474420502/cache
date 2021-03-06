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

func TestCaseBlock(t *testing.T) {
	cache := NewWithBlock(time.Duration(0), func(share interface{}) interface{} {
		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
		if err != nil {
			log.Println(err)
		}

		return string(resp.Content())
	})
	old := cache.Value()
	if old == cache.Value() {
		t.Error("block is error")
	}
}

func TestCase3(t *testing.T) {
	cache := New(time.Second*2, func(share interface{}) interface{} {
		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
		if err != nil {
			log.Println(err)
		}

		return string(resp.Content())
	})
	defer cache.Destroy()

	old := cache.Value()
	cache.Update()
	n := cache.Value()
	if old == n {
		t.Error("value should be updated", n, old)
	}

}

func TestCase4(t *testing.T) {
	cache := NewWithEvery(func(share interface{}) interface{} {
		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
		if err != nil {
			log.Println(err)
		}

		return string(resp.Content())
	})
	defer cache.Destroy()

	old := cache.Value()
	cache.Update()
	n := cache.Value()
	if old == n {
		t.Error("value should be updated", n, old)
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

	time.Sleep(time.Millisecond * 50)

	cache.Value()
	if i != 1 {
		t.Error("onError is error")
	}
}

// func TestUpdateReturnError(t *testing.T) {
// 	cache := New(func() interface{} {
// 		return fmt.Errorf("test errro")
// 	})

// 	var i = 0
// 	cache.SetOnError(func(err interface{}) {
// 		i++
// 		log.Println(i)
// 	})

// 	cache.Update() //第一次更新错误
// 	cache.SetBlock(true)
// 	cache.Value() // 第二次都认为需要更新

// 	if i != 2 {
// 		t.Error("SetOnError")
// 	}
// }

// func TestForce(t *testing.T) {
// 	wg := &sync.WaitGroup{}
// 	cache := New(func() interface{} {
// 		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		return string(resp.Content())
// 	})

// 	cache.SetUpdateInterval(time.Millisecond * 100)

// 	for i := 0; i < 1000; i++ {
// 		go func() {
// 			wg.Add(1)
// 			defer wg.Done()
// 			var old interface{} = cache.Value()
// 			var count = 0
// 			now := time.Now()

// 			for i := 0; i < 1000; i++ {
// 				nvalue := cache.Value()
// 				if old != nvalue {
// 					old = nvalue
// 					count++

// 				}
// 				time.Sleep(time.Millisecond * 1)
// 			}
// 			predict := int(math.Round(float64(time.Since(now).Milliseconds()))) / 100
// 			if !(predict >= count-1 || predict <= count+1) {
// 				t.Error("error", predict, count)
// 			}
// 		}()
// 	}

// 	wg.Wait()
// }

// func TestCaseErr(t *testing.T) {
// 	var i = 0
// 	cache := New(func() interface{} {
// 		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		i++
// 		if i%2 == 0 {
// 			return nil
// 		}

// 		// log.Println(string(resp.Content()))
// 		return string(resp.Content())
// 	})

// 	cache.SetBlock(true)

// 	for n := 0; n < 4; n++ {
// 		cache.Update()
// 		if v, ok := cache.Value().(string); !ok {
// 			t.Error("value error", v)
// 		}
// 	}

// }

// func TestCasePanic(t *testing.T) {
// 	var i = 0
// 	cache := New(func() interface{} {
// 		if i == 2 {
// 			panic("error xixi")
// 		}
// 		i++

// 		return "nono"
// 	})

// 	cache.SetOnError(func(err interface{}) {
// 		if err.(string) != "error xixi" {
// 			t.Error("panic test is error")
// 		}
// 	})

// 	cache.SetBlock(true)

// 	for n := 0; n < 4; n++ {
// 		cache.Update()
// 		if v, ok := cache.Value().(string); !ok {
// 			t.Error("value error", v)
// 		}
// 	}

// }

// func TestCondUpdate(t *testing.T) {

// 	cache := New(func() interface{} {
// 		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		return string(resp.Content())
// 	})

// 	ts := time.Now()
// 	endts := ts.Add(time.Millisecond * 50)
// 	cache.SetUpdateCond(func() bool {
// 		return !ts.Before(endts)
// 	})

// 	cache.SetBlock(true) // 如果不设置阻塞, 第一次Value会因为网络不更新

// 	old := cache.Value()

// 	for i := 0; i < 1; i++ {
// 		time.Sleep(time.Millisecond * 50)
// 		ts = time.Now()
// 		if old == cache.Value() {
// 			t.Error("old should not equal new value")
// 		}
// 	}

// }
