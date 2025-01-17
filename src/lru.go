package main

import (
	"container/list"
	"fmt"
)

type Cache struct {
	capacity         int
	items            map[string]*list.Element // consists of key_name : node [node in d-l-l]
	doublyLinkedList *list.List
}

type Node struct {
	key   string
	value interface{} // using interface as our value could be anything.
}

func LRUCache(capacity int) *Cache {
	return &Cache{
		capacity:         capacity,
		items:            make(map[string]*list.Element),
		doublyLinkedList: list.New(),
	}
}

func (cache *Cache) Get(keyName string) interface{} {

	if node, found := cache.items[keyName]; found {
		cache.doublyLinkedList.MoveToFront(node)
		return node.Value.(*Node).value
	}
	return nil
}

func (cache *Cache) Put(keyName string, value interface{}) bool {
	if node, exist := cache.items[keyName]; exist {
		node.Value.(*Node).value = value
		cache.doublyLinkedList.MoveToFront(node)
		return true
	}

	if cache.doublyLinkedList.Len() >= cache.capacity {
		last := cache.doublyLinkedList.Back()
		if last != nil {
			//delete(cache.items, last.Value.(*Node).key)
			bool1 := cache.Delete(last.Value.(*Node).key)
			if bool1 == true {
				cache.doublyLinkedList.Remove(last)
			} else {
				fmt.Println("Eviction missed ... Error occurs")
			}

		}
	}

	entry := &Node{key: keyName, value: value}
	elem := cache.doublyLinkedList.PushFront(entry)
	cache.items[keyName] = elem
	return true
}

func (cache *Cache) Delete(keyName string) bool {
	if node, found := cache.items[keyName]; found {
		delete(cache.items, keyName)
		cache.doublyLinkedList.Remove(node)
		return true
	}
	return false
}

func (cache *Cache) Size() int {
	return cache.doublyLinkedList.Len()
}

func main() {
	cache := LRUCache(3)

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
