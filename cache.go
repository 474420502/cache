package cache

import (
	"sync"
)

type node[KEY comparable, VALUE any] struct {
	key   KEY
	value VALUE
	prev  *node[KEY, VALUE]
	next  *node[KEY, VALUE]
}

type CacheLRU[KEY comparable, VALUE any] struct {
	mu sync.Mutex

	capacity  int // must >= 32
	batchSize int // size = limitCount >> 3

	store map[KEY]*node[KEY, VALUE]

	zero VALUE

	linkHeader *node[KEY, VALUE]
	linkTail   *node[KEY, VALUE]
}

func NewCacheLRU[KEY comparable, VALUE any](capacity int) *CacheLRU[KEY, VALUE] {
	if capacity < 32 {
		panic("limitSize must >= 32")
	}
	lm := &CacheLRU[KEY, VALUE]{
		capacity:  capacity,
		batchSize: capacity >> 3,
		store:     make(map[KEY]*node[KEY, VALUE]),
	}
	return lm
}

func (lm *CacheLRU[KEY, VALUE]) Get(key KEY) (VALUE, bool) {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	n, ok := lm.store[key]
	if ok {
		lm.moveNodeToHeader(n)
		return n.value, ok
	}
	return lm.zero, ok
}

func (lm *CacheLRU[KEY, VALUE]) Values() (result []VALUE) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	cur := lm.linkHeader

	for cur != nil {
		result = append(result, cur.value)
		cur = cur.next
	}

	return
}

func (lm *CacheLRU[KEY, VALUE]) Remove(key KEY) (val VALUE, ok bool) {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	if n, ok := lm.store[key]; ok {
		lm.remove(n)
		delete(lm.store, key)
		return n.value, true
	}
	return lm.zero, false
}

func (lm *CacheLRU[KEY, VALUE]) Clear() {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lm.store = make(map[KEY]*node[KEY, VALUE])
	lm.linkHeader = nil
	lm.linkTail = nil
}

func (lm *CacheLRU[KEY, VALUE]) Size() int {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	return len(lm.store)
}

func (lm *CacheLRU[KEY, VALUE]) Set(key KEY, value VALUE) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	if len(lm.store) == 0 {

		n := &node[KEY, VALUE]{
			key:   key,
			value: value,
			prev:  nil,
			next:  nil,
		}

		lm.linkHeader = n
		lm.linkTail = n
		lm.store[key] = n

		lm.removeOverLimit()

		return
	}

	n, ok := lm.store[key]
	if ok {
		n.value = value
		lm.moveNodeToHeader(n)
		return
	}

	n = &node[KEY, VALUE]{
		key:   key,
		value: value,
		prev:  nil,
		next:  lm.linkHeader,
	}

	lm.linkHeader.prev = n

	lm.store[key] = n
	lm.linkHeader = n // 换头

	lm.removeOverLimit()
}

func (lm *CacheLRU[KEY, VALUE]) moveNodeToHeader(n *node[KEY, VALUE]) {
	if lm.linkHeader == n {
		return
	}

	if lm.linkTail == n {
		lm.linkTail = n.prev
		lm.linkTail.next = nil

		n.next = lm.linkHeader
		lm.linkHeader.prev = n
		n.prev = nil

		lm.linkHeader = n
		return
	}

	prev := n.prev // 一定不为nil

	prev.next = n.next
	n.next.prev = prev

	n.next = lm.linkHeader
	lm.linkHeader.prev = n
	n.prev = nil
	lm.linkHeader = n
}

func (lm *CacheLRU[KEY, VALUE]) removeOverLimit() {
	if len(lm.store) < lm.capacity {
		return
	}

	cur := lm.linkTail

	for i := 0; i < lm.batchSize; i++ {
		delete(lm.store, cur.key)
		cur = cur.prev
	}

	lm.linkTail = cur
	lm.linkTail.next = nil
}

func (lm *CacheLRU[KEY, VALUE]) remove(n *node[KEY, VALUE]) {
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		lm.linkHeader = n.next
	}

	if n.next != nil {
		n.next.prev = n.prev
	} else {
		lm.linkTail = n.prev
	}
}
