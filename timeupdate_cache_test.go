package cache

import (
	"log"
)

func init() {
	log.Println("docker run -p 80:80 kennethreitz/httpbin. and set /etc/host 127.0.0.1 httpbin.org")
}

// func TestCase1(t *testing.T) {
// 	cache := New(time.Millisecond*50, func(share interface{}) interface{} {
// 		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		return string(resp.Content())
// 	})

// 	old := cache.Value()
// 	for i := 0; i < 2; i++ {
// 		time.Sleep(time.Millisecond * 70) //因为更细是异步虽然触发了更新. 异步更新不算时间
// 		n := cache.Value()
// 		if old == n {
// 			t.Error("value should be updated", n, old)
// 		}
// 		old = cache.Value()
// 		old = cache.Value()
// 		if old != n {
// 			t.Error("value should be updated", n, old)
// 		}
// 	}

// }

// func TestCase2(t *testing.T) {
// 	cache := New(time.Millisecond*50, func(share interface{}) interface{} {
// 		if share == nil {
// 			log.Println("share is nil", share)
// 			time.Sleep(time.Millisecond * 50)
// 			return nil
// 		}

// 		if share != 1 {
// 			t.Error("share is not 1")
// 		}

// 		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		return string(resp.Content())
// 	})

// 	cache.SetShare(1)
// 	old := cache.Value()
// 	if old != nil {
// 		t.Error("old not nil")
// 	}
// 	time.Sleep(time.Millisecond * 100)
// 	if cache.Value() == nil {
// 		t.Error("uid get error")
// 	}
// }

// // func TestBlockWithCond(t *testing.T) {

// // }

// func TestError(t *testing.T) {
// 	cache := New(time.Millisecond*50, func(share interface{}) interface{} {
// 		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		panic("error")
// 		return string(resp.Content())
// 	})

// 	var i = 0
// 	cache.SetOnUpdateError(func(err interface{}) {
// 		i = 1
// 	})

// 	time.Sleep(time.Millisecond * 100)

// 	cache.Value()
// 	if i != 1 {
// 		t.Error("onError is error")
// 	}
// }

// func TestMulti(t *testing.T) {

// 	defer func() {
// 		if err := recover(); err != nil {
// 			t.Error(err)
// 		}
// 	}()

// 	var i = 0

// 	cache := New(time.Millisecond*50, func(share interface{}) interface{} {
// 		resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		i++
// 		if i%10 == 0 {
// 			panic("error")
// 		}
// 		return string(resp.Content())
// 	})

// 	cache.SetOnUpdateError(func(err interface{}) {
// 		i = 1
// 	})
// 	wg := &sync.WaitGroup{}
// 	for n := 0; n < 100; n++ {
// 		wg.Add(1)
// 		go func(wg *sync.WaitGroup) {
// 			defer wg.Done()
// 			for x := 0; x < 400; x++ {
// 				cache.Value()
// 				cache.GetUpdate()
// 				time.Sleep(time.Millisecond)
// 			}
// 		}(wg)
// 	}

// 	wg.Wait()
// }

// func TestRuntimeCleanCache(t *testing.T) {
// 	func() {

// 		cache := New(time.Millisecond*50, func(share interface{}) interface{} {
// 			resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
// 			if err != nil {
// 				log.Println(err)
// 			}
// 			return string(resp.Content())
// 		})

// 		var result string = cache.Value().(string)
// 		for i := 0; i < 3; i++ {
// 			cvalue := cache.Value().(string)
// 			if result != cvalue {
// 				t.Error(result, cvalue)
// 			}
// 		}

// 		for i := 0; i < 3; i++ {
// 			// time.Sleep(time.Millisecond * 30)
// 			cache.ForceUpdate()
// 			cvalue := cache.Value().(string)
// 			if result == cvalue {
// 				t.Error(result, cvalue)
// 			}
// 			result = cvalue
// 		}

// 	}()

// 	runtime.GC()
// 	time.Sleep(time.Second * 2)
// }

// func TestForceUpdate(t *testing.T) {

// }

// func TestShareProblem(t *testing.T) {
// 	wg := &sync.WaitGroup{}
// 	for i := 0; i < 10000000; i++ {
// 		wg.Add(1)
// 		go func(wg *sync.WaitGroup) {
// 			defer wg.Done()
// 			c := createShareCache()
// 			if c.Value() == nil {
// 				return
// 			}
// 		}(wg)
// 	}

// 	wg.Wait()
// }

// func createShareCache() Cache {
// 	cache := New(time.Second, func(share interface{}) interface{} {
// 		s := share.(string)
// 		return s
// 	})
// 	cache.SetShare("1")
// 	return cache
// }
