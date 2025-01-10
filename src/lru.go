package main

import (
	"container/list"
	"fmt"
)

type Cache struct {
	capacity int
	items    map[string]*list.Element
	list     *list.List
}

type Entry struct {
	key   string
	value interface{}
}

func NewLRUCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (lru *Cache) Get(key string) interface{} {
	if elem, found := lru.items[key]; found {
		lru.list.MoveToFront(elem)
		return elem.Value.(*Entry).value
	}
	return nil
}

func (lru *Cache) Put(key_name string, value interface{}) bool {
	if elem, exist := lru.items[key_name]; exist {
		elem.Value.(*Entry).value = value
		lru.list.MoveToFront(elem)
		return true
	}

	if lru.list.Len() >= lru.capacity {
		last := lru.list.Back()
		if last != nil {
			delete(lru.items, last.Value.(*Entry).key)
			lru.list.Remove(last)
		}
	}

	entry := &Entry{key: key_name, value: value}
	elem := lru.list.PushFront(entry)
	lru.items[key_name] = elem
	return true
}

func (lru *Cache) Delete(key string) bool {
	if elem, found := lru.items[key]; found {
		delete(lru.items, key)
		lru.list.Remove(elem)
		return true
	}
	return false
}

func (lru *Cache) Size() int {
	return lru.list.Len()
}

func main() {
	cache := NewLRUCache(3)

	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)

	fmt.Println(cache.Get("a")) // Output: 1, true
	cache.Put("d", 4)           // "b" will be evicted

	value := cache.Get("b")
	fmt.Println(value) // Output: false

	cache.Put("e", 5)           // "c" will be evicted
	fmt.Println(cache.Get("c")) // Output: <nil>, false
}
