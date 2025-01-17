package main

import (
	"container/list"
	"fmt"
)

type Cache struct {
	bucketSize       int
	items            map[string]*list.Element // consists of key_name : node [node in d-l-l]
	doublyLinkedList *list.List
}

type Node struct {
	key   string
	value interface{} // using interface as our value could be anything.
}

func LRUCache(bucketSize int) *Cache {
	return &Cache{
		bucketSize:       bucketSize,
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
			fmt.Println("Failed to invalidate ... couldn't put object to cache")
			return false
		}
	}

	entry := &Node{key: keyName, value: value}
	elem := cache.doublyLinkedList.PushFront(entry)
	cache.items[keyName] = elem
	fmt.Printf("object of key name :: %v , has been added to cache successfully\n", keyName)
	return true
}

func (cache *Cache) Invalidate() bool {
	last := cache.doublyLinkedList.Back()
	if last != nil {
		return cache.Delete(last.Value.(*Node).key)
	}

	return false
}

func (cache *Cache) Delete(keyName string) bool {
	if node, found := cache.items[keyName]; found {
		delete(cache.items, keyName)
		cache.doublyLinkedList.Remove(node)
		fmt.Printf("object of key name :: %v , has been removed from cache successfully\n", keyName)
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
	fmt.Println(cache.Get("e")) // Output: <nil>, false
}
