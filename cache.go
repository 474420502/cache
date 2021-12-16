// Package cache 缓存 异步更新
package cache

import (
	"log"
	"runtime"
	"time"
)

type Cache interface {
	SetShare(share interface{})
	SetOnUpdateError(func(err interface{}))
	Value() interface{}
	ForceUpdate()
	GetUpdate() time.Time
}

// UpdateMehtod 更新方法
type UpdateMehtod func(share interface{}) interface{}

// New 创建一个Cache对象timeupdate
func New(interval time.Duration, updateMethod UpdateMehtod) Cache {

	cbackup := &cacheIntervalBackup{
		isDestroy:     make(chan struct{}),
		isForceUpdate: make(chan struct{}),
		interval:      time.NewTicker(interval),
		updateMehtod:  updateMethod,
		onUpdateError: func(err interface{}) {
			log.Println(err)
		},
	}

	cache := &CacheInterval{
		cache: cbackup,
	}

	runtime.SetFinalizer(cache, func(obj *CacheInterval) {
		obj.cache.Destroy()
	})

	go func() {
		for {
			select {
			case <-cbackup.isDestroy:
				return
			case <-cbackup.isForceUpdate:
				func() {
					defer func() { cbackup.isForceUpdate <- struct{}{} }()
					cbackup.update()
					<-cbackup.interval.C
				}()
			case <-cbackup.interval.C:
				cbackup.update()
			}
		}
	}()

	cbackup.ForceUpdate()
	return cache
}
