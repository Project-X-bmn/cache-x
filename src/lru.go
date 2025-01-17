package main

import (
	"container/list"
	"fmt"
)

type Cache struct {
	bucketSize       int
	threshold        int
	items            map[string]*list.Element // consists of key_name : node [node in d-l-l]
	doublyLinkedList *list.List
}

type Node struct {
	key   string
	value interface{} // using interface as our value could be anything.
}

func LRUCache(capacity int, size int) *Cache {
	return &Cache{
		bucketSize:       capacity,
		threshold:        size,
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

	if cache.doublyLinkedList.Len() >= cache.bucketSize {
		if status := cache.Invalidate(); !status {
			return false
		}

	}
	entry := &Node{key: keyName, value: value}
	elem := cache.doublyLinkedList.PushFront(entry)
	cache.items[keyName] = elem
	return true

}

func (cache *Cache) Invalidate() bool {

	flag := false
	for i := 0; i < cache.threshold/cache.bucketSize; i++ {
		last := cache.doublyLinkedList.Back()
		if last != nil {
			flag = cache.Delete(last.Value.(*Node).key)
			flag = true
		}

		if flag == false {
			return false
		}
	}
	return flag

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
	cache := LRUCache(10, 30)

	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)
	cache.Put("d", 1)
	cache.Put("e", 2)
	cache.Put("f", 3)
	cache.Put("g", 1)
	cache.Put("h", 2)
	cache.Put("i", 3)

	fmt.Println(cache.doublyLinkedList.Len())

	cache.Put("Nik", 25)
	fmt.Println(cache.doublyLinkedList.Len())

	cache.Put("bar", 21)
	fmt.Println(cache.doublyLinkedList.Len())

	//fmt.Println(cache.Get("a")) // Output: 1, true
	//cache.Put("d", 4)           // "b" will be evicted
	//
	//value := cache.Get("b")
	//fmt.Println(value) // Output: false
	//
	//cache.Put("e", 5)           // "c" will be evicted
	//fmt.Println(cache.Get("c")) // Output: <nil>, false

}
