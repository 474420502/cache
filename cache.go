// Package cache 缓存 异步更新
package cache

import (
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Cache 缓存
type Cache struct {
	slock sync.Mutex
	share interface{}

	isDestroy int32

	updateMehtod UpdateMehtod

	onUpdateError func(err interface{})

	lastUpdate time.Time
	interval   time.Duration

	vLock sync.Mutex
	value interface{}
}

// UpdateMehtod 更新方法
type UpdateMehtod func(share interface{}) interface{}

// New 创建一个Cache对象
func New(interval time.Duration, u UpdateMehtod) *Cache {

	c := &Cache{
		updateMehtod: u,
		interval:     interval,
		onUpdateError: func(err interface{}) {
			log.Println(err)
		},
	}

	c.update()
	go func() {
		for {
			time.Sleep(c.interval)
			if atomic.LoadInt32(&c.isDestroy) == 1 {
				break
			}
			c.update()
		}
	}()

	runtime.Gosched()
	return c
}

func (cache *Cache) SetShare(share interface{}) {
	cache.slock.Lock()
	defer cache.slock.Unlock()
	cache.share = share
}

// SetOnUpdateError 默认false
func (cache *Cache) SetOnUpdateError(errFunc func(err interface{})) {
	cache.slock.Lock()
	defer cache.slock.Unlock()
	cache.onUpdateError = errFunc
}

// Destroy 异步更新必须调用Destroy, 销毁对象
func (cache *Cache) Destroy() {
	atomic.StoreInt32(&cache.isDestroy, 1)
}

// update 主动更新 没锁
func (cache *Cache) update() {
	cache.slock.Lock()
	defer cache.slock.Unlock()

	defer func() {
		cache.lastUpdate = time.Now()
		if err := recover(); err != nil {
			cache.onUpdateError(err)
		}
	}()

	func() {

		value := cache.updateMehtod(cache.share)
		if value == nil {
			return
		}

		if err, ok := value.(error); ok {
			cache.onUpdateError(err)
			return
		}

		cache.vLock.Lock()
		cache.value = value
		cache.vLock.Unlock()
	}()
}

// Value 获取缓存的值
func (cache *Cache) Value() interface{} {
	cache.vLock.Lock()
	defer cache.vLock.Unlock()
	return cache.value
}

// Value 获取缓存的值
func (cache *Cache) GetUpdate() time.Time {
	cache.slock.Lock()
	defer cache.slock.Unlock()
	return cache.lastUpdate
}
