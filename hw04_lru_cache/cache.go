package hw04lrucache

import "sync"

type Key string

// Элемент кэша хранит в себе ключ, по которому он лежит в словаре, и само значение.
type cacheElem struct {
	key   Key
	value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mutex    sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	// Блокировка для потокобезопасности
	cache.mutex.Lock()

	// Снятие блокировки
	defer cache.mutex.Unlock()

	// Элемент есть в кэше - обновить значение, переместить в начало очереди, вернуть true
	if elem, exist := cache.items[key]; exist {
		elem.Value.(*cacheElem).value = value
		cache.queue.MoveToFront(elem)
		return true
	}

	// Элемент отсутствует в словаре - добавить в словарь,  поместить в начало очереди.
	// если размер очереди больше ёмкости кэша - удалить последний элемент из очереди и его значение из словаря

	// Сперва удалить (это по большей части касается словоря, чтобы не выделять место из-за превышения емкости)
	if cache.queue.Len() >= cache.capacity {
		last := cache.queue.Back()
		cache.queue.Remove(last)
		delete(cache.items, last.Value.(*cacheElem).key)
	}

	// Потом добавить
	newElem := cache.queue.PushFront(&cacheElem{key, value})
	cache.items[key] = newElem

	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	// Блокировка для потокобезопасности
	cache.mutex.Lock()

	// Снятие блокировки
	defer cache.mutex.Unlock()

	elem, exist := cache.items[key]

	// Если элемента нет в словаре, то вернуть nil и false
	if !exist {
		return nil, false
	}

	// Если элемент присутствует в словаре, то переместить элемент в начало очереди и вернуть его значение и true
	cache.queue.MoveToFront(elem)
	return elem.Value.(*cacheElem).value, true
}

func (cache *lruCache) Clear() {
	// Блокировка для потокобезопасности
	cache.mutex.Lock()

	// Снятие блокировки
	defer cache.mutex.Unlock()

	// Создание новой очереди и словаря (удаление будет делать GC)
	cache.queue = NewList()
	cache.items = make(map[Key]*ListItem, cache.capacity)
}
