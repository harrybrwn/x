package collections

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestHashSet(t *testing.T) {
	s := NewSyncHashSet[string]()
	if s.Has("a") {
		t.Error("should not contain this key")
	}
	if s.Del("a") {
		t.Error("should not contain this key")
	}
	s.Put("a")
	if !s.Has("a") {
		t.Error("should contain this key")
	}
	if !s.Del("a") {
		t.Error("should contain this key")
	}
	if s.LoadPut("b") {
		t.Error("should not contain this key")
	}
	if !s.Has("b") {
		t.Error("should contain this key")
	}
}

// Run this with -race
func TestSyncHashSet(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	var wg sync.WaitGroup
	s := NewSyncHashSet[string]()
	wg.Go(func() {
		for range 256 {
			s.Put("a")
		}
	})
	var counter int64
	wg.Go(func() {
		for range 256 {
			ok := s.Has("a")
			if ok {
				atomic.AddInt64(&counter, 1)
			}
		}
	})
	wg.Go(func() {
		for range 256 {
			s.Del("a")
		}
	})
	wg.Wait()
	if counter == 0 {
		t.Error("should not be zero")
	}
}
