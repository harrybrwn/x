package collections

import "sync"

type Set[T comparable] interface {
	Put(T)
	// Del deletes an item from the set. If the item was present in the set and
	// successfully deleted then it returns true, otherwise if the item was not
	// in the set it returns false.
	Del(T) bool
	Has(T) bool
}

type HashSet[T comparable] map[T]struct{}

func (hs HashSet[T]) Put(v T) { hs[v] = struct{}{} }

func (hs HashSet[T]) Del(v T) bool {
	_, ok := hs[v]
	delete(hs, v)
	return ok
}

func (hs HashSet[T]) Has(v T) bool {
	_, ok := hs[v]
	return ok
}

// LoadPut checks if the set has the given value and returns true or false while
// adding the value to the set.
func (hs HashSet[T]) LoadPut(v T) bool {
	_, ok := hs[v]
	hs[v] = struct{}{}
	return ok
}

func NewSyncHashSet[T comparable]() *SyncHashSet[T] {
	return &SyncHashSet[T]{s: make(HashSet[T])}
}

type SyncHashSet[T comparable] struct {
	mu sync.RWMutex
	s  HashSet[T]
}

func (hs *SyncHashSet[T]) Put(v T) {
	hs.mu.Lock()
	hs.s.Put(v)
	hs.mu.Unlock()
}

func (hs *SyncHashSet[T]) Del(v T) bool {
	hs.mu.Lock()
	ok := hs.s.Del(v)
	hs.mu.Unlock()
	return ok
}

func (hs *SyncHashSet[T]) Has(v T) bool {
	hs.mu.RLock()
	ok := hs.s.Has(v)
	hs.mu.RUnlock()
	return ok
}

func (hs *SyncHashSet[T]) LoadPut(v T) bool {
	hs.mu.Lock()
	ok := hs.s.LoadPut(v)
	hs.mu.Unlock()
	return ok
}
