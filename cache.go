// Package cache 缓存 异步更新
package cache

import (
	"log"
	"runtime"
	"sync/atomic"
	"time"
)

type Cache interface {
	SetShare(share interface{})
	SetOnUpdateError(func(err interface{}))
	Destroy()
	Value() interface{}
	GetUpdate() time.Time
}

// UpdateMehtod 更新方法
type UpdateMehtod func(share interface{}) interface{}

// New 创建一个Cache对象
func New(interval time.Duration, u UpdateMehtod) Cache {

	c := &CacheInterval{
		updateMehtod: u,
		interval:     interval,
		onUpdateError: func(err interface{}) {
			log.Println(err)
		},
	}

	go func() {
		runtime.Gosched()
		for {
			if atomic.LoadInt32(&c.isDestroy) == 1 {
				break
			}
			c.update()
			time.Sleep(c.interval)
		}
	}()

	return c
}
