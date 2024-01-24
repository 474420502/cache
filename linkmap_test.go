package cache

import (
	"log"
	"testing"
)

func TestLinkMap(t *testing.T) {
	lm := newLinkMap[int, int](32)

	for i := 0; i < 60; i++ {
		lm.Set(i, i)
	}

	lm.Set(3, 3)
	lm.Set(2, 2)
	lm.Set(4, 4)
	lm.Get(2)
	lm.Get(1)

	log.Println(lm.Values(), len(lm.Values()), len(lm.store))
}
