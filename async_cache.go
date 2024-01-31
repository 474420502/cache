// Package cache 缓存 异步更新
package cache

import (
	"fmt"
	"sync"
	"time"
)

type AsyncCacheContext[SHARE any] struct {
	params  any
	current SHARE
	err     any
}

func (ctx *AsyncCacheContext[SHARE]) Params() any {
	return ctx.params
}

func (ctx *AsyncCacheContext[SHARE]) GetCurrent() SHARE {
	return ctx.current
}

func (ctx *AsyncCacheContext[SHARE]) Error() any {
	return ctx.err
}

// AsyncCache 双缓冲
type AsyncCache[SHARE any] struct {
	current SHARE

	params        any
	updateHandler func(ctx *AsyncCacheContext[SHARE]) (SHARE, error)
	errorHandler  func(ctx *AsyncCacheContext[SHARE])

	firstOnce   sync.Once
	firstUpdate chan bool

	lock sync.Mutex
}

func NewAsyncCache[SHARE any](updateHandler func(ctx *AsyncCacheContext[SHARE]) (SHARE, error)) *AsyncCache[SHARE] {
	ac := &AsyncCache[SHARE]{
		updateHandler: updateHandler,
		firstUpdate:   make(chan bool),
	}

	NewAsyncSymbiosis().StartLoopWith(ac, func() {
		ac.updateProcess()
	})

	return ac
}

func (cache *AsyncCache[SHARE]) Get() SHARE {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	current := cache.current
	return current
}

func (cache *AsyncCache[SHARE]) GetSync(doWithGet func(current SHARE)) {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	doWithGet(cache.current)
}

func (cache *AsyncCache[SHARE]) SetParams(params any) {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	cache.params = params
}

func (cache *AsyncCache[SHARE]) SetError(handler func(ctx *AsyncCacheContext[SHARE])) {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	cache.errorHandler = handler
}

func (cache *AsyncCache[SHARE]) WaitFirstUpdate(waitTime time.Duration) error {

	select {
	case <-cache.firstUpdate:
		return nil
	case <-time.After(waitTime):
		return fmt.Errorf("timeout waiting for first update")
	}

}

func (cache *AsyncCache[SHARE]) updateProcess() {

	defer func() {
		if err := recover(); err != nil {
			cache.errorHandler(&AsyncCacheContext[SHARE]{err: err})
		}
	}()

	cache.lock.Lock()
	params := cache.params
	current := cache.current
	cache.lock.Unlock()

	readyValue, err := cache.updateHandler(&AsyncCacheContext[SHARE]{
		params:  params,
		current: current,
	})

	if err != nil {
		panic(err)
	}

	cache.lock.Lock()
	defer cache.lock.Unlock()
	cache.current = readyValue
	cache.firstOnce.Do(func() {
		cache.firstUpdate <- true
	})
}
