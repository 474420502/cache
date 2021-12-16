package example

import (
	"log"
	"testing"
	"time"

	"github.com/474420502/cache"
	"github.com/474420502/gcurl"
)

func TestCacheTimeInterval(t *testing.T) {
	// 实例
	cache := cache.New(
		time.Millisecond*50, // 每 50 millisecond 更新一次
		func(share interface{}) interface{} { // 更新的方法
			resp, err := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
			if err != nil {
				log.Println(err)
			}

			return string(resp.Content()) // 返回更新的时间. 这个数据会在不更新期间缓存
		})

	old := cache.Value() // 旧值
	for i := 0; i < 2; i++ {
		time.Sleep(time.Millisecond * 70) // 70 Millisecond 查询一次值
		n := cache.Value()
		if old == n {
			t.Error("value should be updated", n, old)
		}
		old = cache.Value()
		if old != n { // 如果在 70 Millisecond 内值不一样就报错
			t.Error("value should be updated", n, old)
		}
	}

	// 初始化后只需要获取
	log.Println(cache.Value()) // 该值为 cache.UpdateMehtod 的值. 可用于缓存 一些远程更新数据 但是不需要频繁更新.
}
