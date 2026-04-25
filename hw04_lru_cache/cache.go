package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.RWMutex //будем использовать для всех трёх полей кэша
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key //нужно только для удаления старого значения из map
	value any
}

func (l *lruCache) Set(key Key, value interface{}) (isInserted bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if item, ok := l.items[key]; ok { //comma ok: true, если ключ в мапе нашелся, и false, если нет
		ci, ok := item.Value.(*cacheItem)
		if !ok {
			panic("cache item is not a cacheItem!!!")
		}
		ci.value = value
		l.queue.MoveToFront(item)
		return true
	}

	if l.queue.Len() >= l.capacity {
		lastItem := l.queue.Back()
		delete(l.items, lastItem.Value.(*cacheItem).key)
		l.queue.Remove(lastItem)
	}
	newItem := l.queue.PushFront(&cacheItem{key: key, value: value})
	l.items[key] = newItem
	return false
}

func (l *lruCache) Get(key Key) (any, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	item := l.items[key]
	if item != nil {
		l.queue.MoveToFront(item)
		ci, ok := item.Value.(*cacheItem)
		if !ok {
			panic("cache item is not a cacheItem!!!")
		}
		return ci.value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.items = make(map[Key]*ListItem)
	l.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
