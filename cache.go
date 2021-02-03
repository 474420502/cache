// Package cache 缓存 异步更新
package cache

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// Cache 缓存
type Cache struct {
	updateCond *updateCond

	isBlock bool

	isUpdating   int32
	updateMehtod UpdateMehtod
	onError      func(err interface{})

	valueLock sync.Mutex
	value     interface{}
}

// UpdateMehtod 更新方法
type UpdateMehtod func() interface{}

// New 创建一个Cache对象
func New(u UpdateMehtod) *Cache {
	c := &Cache{
		updateMehtod: u,
		isBlock:      false,
		onError: func(err interface{}) {
			log.Println(err)
		},
	}

	return c
}

// SetOnError 默认false
func (cache *Cache) SetOnError(errFunc func(err interface{})) {
	cache.valueLock.Lock()
	defer cache.valueLock.Unlock()
	cache.onError = errFunc
}

// SetUpdateMethod 设置cache更新方法
func (cache *Cache) SetUpdateMethod(method UpdateMehtod) {
	cache.valueLock.Lock()
	defer cache.valueLock.Unlock()
	cache.updateMehtod = method
}

// SetUpdateCond 设置cache更新的条件. 时间间隔更新失效
func (cache *Cache) SetUpdateCond(method func() bool) {
	cache.valueLock.Lock()
	defer cache.valueLock.Unlock()
	cache.updateCond = &updateCond{
		cond: method,
	}
}

// SetUpdateInterval 设置cache更新的条件. 时间间隔更新. SetUpdateCond会失效. Cond也能完成所有更新方式
func (cache *Cache) SetUpdateInterval(interval time.Duration) {
	cache.valueLock.Lock()
	defer cache.valueLock.Unlock()
	cache.updateCond = &updateCond{
		interval: interval,
	}
}

// SetBlock 默认false
func (cache *Cache) SetBlock(is bool) {
	cache.valueLock.Lock()
	defer cache.valueLock.Unlock()
	cache.isBlock = is
}

// Update 主动更新
func (cache *Cache) Update() {

	cache.valueLock.Lock()
	defer cache.valueLock.Unlock()

	if cache.value == nil {
		cache.value = cache.updateMehtod()
		if cache.value == nil {
			log.Println("warning data first update is nil")
		}
		return
	}

	if cache.isBlock {
		v := cache.updateMehtod()
		switch value := v.(type) {
		case nil:
		case error:
			cache.onError(value)
		default:
			cache.value = v
		}

		return
	}

	switch {
	case cache.updateCond.cond != nil:
		if cache.updateCond.cond() {
			cache.asyncUpdating()
		}
	case time.Since(cache.updateCond.updateAt) >= cache.updateCond.interval:
		cache.asyncUpdating()
	}
}

func (cache *Cache) Value() interface{} {

	defer func() {
		if err := recover(); err != nil {
			cache.onError(err)
		}
	}()

	if cache.updateCond != nil {
		cache.Update()
		return cache.value
	}

	cache.valueLock.Lock()
	defer cache.valueLock.Unlock()
	if cache.value == nil {
		cache.value = cache.updateMehtod()
		if cache.value == nil {
			log.Println("warning data first update is nil")
		}
	}
	return cache.value
}

func (cache *Cache) asyncUpdating() {
	if atomic.CompareAndSwapInt32(&cache.isUpdating, 0, 1) {
		go func() {
			defer atomic.StoreInt32(&cache.isUpdating, 0)
			defer func() {
				if err := recover(); err != nil {
					log.Println("recover:", err)
				}
			}()

			v := cache.updateMehtod()

			switch value := v.(type) {
			case nil:
				return
			case error:
				cache.onError(value)
			default:
				cache.valueLock.Lock()
				cache.value = v
				cache.valueLock.Unlock()
			}
		}()
	}
}

type updateCond struct {
	cond     func() bool
	updateAt time.Time
	interval time.Duration
}
