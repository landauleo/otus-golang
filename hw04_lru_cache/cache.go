package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key //нужно только для удаления старого значения из map
	value interface{}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	item := l.items[key]
	if item != nil {
		//type assertion МАТЬ ЕГО
		//подменяем значение
		item.Value.(*cacheItem).value = value
		l.queue.MoveToFront(item)
		return true
	}

	if l.queue.Len() > l.capacity {
		lastItem := l.queue.Back()
		delete(l.items, lastItem.Value.(*cacheItem).key)
		l.queue.Remove(l.queue.Back())
	}
	newItem := l.queue.PushFront(&cacheItem{key: key, value: value})
	l.items[key] = newItem
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	item := l.items[key]
	if item != nil {
		l.queue.MoveToFront(item)
		return item.Value.(*cacheItem).value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
