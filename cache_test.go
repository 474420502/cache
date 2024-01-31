package cache

import (
	"strconv"
	"testing"
)

func TestLinkMap(t *testing.T) {
	// lm := NewCacheLRU[int, int](32)

	// for i := 0; i < 60; i++ {
	// 	lm.Set(i, i)
	// }

	// lm.Set(3, 3)
	// lm.Set(2, 2)
	// lm.Set(4, 4)
	// lm.Get(2)
	// lm.Get(1)

	// log.Println(lm.Values(), len(lm.Values()), len(lm.store))
}

func TestCacheLRU(t *testing.T) {
	cache := NewCacheLRU[int, string](32)

	// 添加数据
	cache.Set(1, "one")
	cache.Set(2, "two")
	cache.Set(3, "three")
	cache.Set(4, "four")
	cache.Set(5, "five")

	// 获取数据
	value, ok := cache.Get(3)
	if !ok || value != "three" {
		t.Errorf("Get failed. Expected value: three, Got value: %s", value)
	}

	// 更新数据
	cache.Set(3, "new_three")
	value, ok = cache.Get(3)
	if !ok || value != "new_three" {
		t.Errorf("Update failed. Expected value: new_three, Got value: %s", value)
	}

	// 移除数据
	cache.Remove(4)
	_, ok = cache.Get(4)
	if ok {
		t.Errorf("Remove failed. Key 4 should have been removed")
	}

	// 清空缓存
	cache.Clear()
	size := cache.Size()
	if size != 0 || cache.linkHeader != nil || cache.linkTail != nil {
		t.Errorf("Clear failed. Expected size: 0, Got size: %d", size)
	}

	// 添加超过容量的数据，检查缓存淘汰
	for i := 0; i < 33; i++ {
		cache.Set(i, strconv.Itoa(i))
	}

	_, ok = cache.Get(1)
	if ok {
		t.Errorf("LRU eviction failed. Key 1 should have been evicted")
	}

	// 获取所有值
	values := cache.Values()
	if cache.capacity+1-cache.batchSize != len(values) {
		t.Errorf("eviction error? %d - %d = %d?", cache.capacity+1, cache.batchSize, len(values))
	}
	// expectedValues := []string{"two", "three", "four", "five", "six"}
	// if !stringSlicesEqual(values, expectedValues) {
	// 	t.Errorf("Values mismatch. Expected: %v, Got: %v", expectedValues, values)
	// }
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
