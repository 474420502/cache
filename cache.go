// Package cache 缓存 异步更新
package cache

import (
	"log"
	"runtime"
	"sync"
	"time"
)

// Cache 缓存
type Cache struct {
	share interface{}

	isBlock   bool
	isDestroy bool

	updateMehtod UpdateMehtod

	onUpdateError func(err interface{})

	LastUpdate time.Time
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
		onUpdateError: func(err interface{}) {
			log.Println(err)
		},
	}

	go func() {

		for {

			c.vLock.Lock()
			c.update()
			if c.isDestroy {
				c.vLock.Unlock()
				break
			} else {
				c.vLock.Unlock()
			}

			time.Sleep(c.interval)
		}
	}()

	runtime.Gosched()
	return c
}

// New 创建一个Cache对象, 异步更新必须调用Destroy 销毁
func NewWithBlock(interval time.Duration, u UpdateMehtod) *Cache {
	c := &Cache{
		updateMehtod: u,
		isBlock:      true,
		interval:     interval,
		onUpdateError: func(err interface{}) {
			log.Println(err)
		},
	}

	return c
}

// NewWithEvery 创建一个Cache对象. 必须每次都触发UpdateMethod可以自定以条件
func NewWithEvery(u UpdateMehtod) *Cache {
	c := &Cache{
		updateMehtod: u,
		isBlock:      true,
		onUpdateError: func(err interface{}) {
			log.Println(err)
		},
	}

	return c
}

func (cache *Cache) SetShare(share interface{}) {
	cache.vLock.Lock()
	defer cache.vLock.Unlock()
	cache.share = share
}

// SetOnUpdateError 默认false
func (cache *Cache) SetOnUpdateError(errFunc func(err interface{})) {
	cache.vLock.Lock()
	defer cache.vLock.Unlock()
	cache.onUpdateError = errFunc
}

// Destroy 异步更新必须调用Destroy, 销毁对象
func (cache *Cache) Destroy() {
	cache.vLock.Lock()
	defer cache.vLock.Unlock()
	cache.isBlock = true
}

// update 主动更新 没锁
func (cache *Cache) update() {

	defer func() {
		cache.LastUpdate = time.Now()
		if err := recover(); err != nil {
			cache.onUpdateError(err)
		}
	}()

	value := cache.updateMehtod(cache.share)
	if value == nil {
		return
	}

	if err, ok := value.(error); ok {
		cache.onUpdateError(err)
		return
	}

	cache.value = value
}

// Update 主动更新
func (cache *Cache) Update() {
	cache.vLock.Lock()
	defer cache.vLock.Unlock()
	cache.update()
}

// Value 获取缓存的值
func (cache *Cache) Value() interface{} {
	cache.vLock.Lock()
	defer cache.vLock.Unlock()

	if cache.isBlock && time.Since(cache.LastUpdate) >= cache.interval {
		cache.update()
	}

	return cache.value
}
