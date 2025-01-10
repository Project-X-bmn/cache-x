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

type Node struct {
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

func (lru *Cache) Get(keyName string) interface{} {
	if item, found := lru.items[keyName]; found {
		lru.list.MoveToFront(item)
		return item.Value.(*Node).value
	}
	return nil
}

func (lru *Cache) Put(keyName string, value interface{}) bool {
	if item, exist := lru.items[keyName]; exist {
		item.Value.(*Node).value = value
		lru.list.MoveToFront(item)
		return true
	}

	if lru.list.Len() >= lru.capacity {
		last := lru.list.Back()
		if last != nil {
			delete(lru.items, last.Value.(*Node).key)
			lru.list.Remove(last)
		}
	}

	entry := &Node{key: keyName, value: value}
	elem := lru.list.PushFront(entry)
	lru.items[keyName] = elem
	return true
}

func (lru *Cache) Delete(keyName string) bool {
	if item, found := lru.items[keyName]; found {
		delete(lru.items, keyName)
		lru.list.Remove(item)
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
