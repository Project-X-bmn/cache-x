package learnings

import (
	"container/list"
	"fmt"
)

const maxCacheSize = 1 << 30 // 1GB in bytes

type Cache struct {
	currentSize      int64 // Current total memory usage of the cache (in bytes)
	items            map[string]*list.Element
	doublyLinkedList *list.List
}

type Node struct {
	key   string
	value interface{} // Store file content or other types of data
	size  int64       // Size of the value in bytes
}

func LRUCache() *Cache {
	return &Cache{
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
	// Get the size of the value (for example, file size in bytes)
	nodeSize := int64(0)
	switch v := value.(type) {
	case []byte:
		nodeSize = int64(len(v)) // Size of byte array (file content)
	default:
		return false
	}

	// If the node exists, update it
	if node, exist := cache.items[keyName]; exist {
		// Update size
		cache.currentSize -= node.Value.(*Node).size
		node.Value.(*Node).value = value
		node.Value.(*Node).size = nodeSize
		cache.currentSize += nodeSize
		cache.doublyLinkedList.MoveToFront(node)
		return true
	}

	// If the cache exceeds the max capacity, evict
	for cache.currentSize+nodeSize > maxCacheSize {
		cache.Invalidate()
	}

	// Add the new node
	entry := &Node{key: keyName, value: value, size: nodeSize}
	elem := cache.doublyLinkedList.PushFront(entry)
	cache.items[keyName] = elem
	cache.currentSize += nodeSize
	return true
}

func (cache *Cache) Invalidate() {
	last := cache.doublyLinkedList.Back()
	if last != nil {
		cache.Delete(last.Value.(*Node).key)
	}
}

func (cache *Cache) Delete(keyName string) bool {
	if node, found := cache.items[keyName]; found {
		cache.currentSize -= node.Value.(*Node).size
		delete(cache.items, keyName)
		cache.doublyLinkedList.Remove(node)
		return true
	}
	return false
}

func (cache *Cache) Size() int64 {
	return cache.currentSize
}

func main() {
	cache := LRUCache()

	// Store a 100MB file (for example purposes)
	filePath := "example.txt"
	data := make([]byte, 100<<20) // 100MB of dummy data
	cache.Put(filePath, data)
	fmt.Println("Cache size:", cache.Size(), "bytes")

	// Simulate storing a 1GB file (if you want)
	largeFile := make([]byte, 1<<30) // 1GB of dummy data
	cache.Put("large_file.txt", largeFile)
	fmt.Println("Cache size after adding 1GB file:", cache.Size(), "bytes")

	// Test eviction when the cache exceeds 1GB
	// You can add multiple files to see eviction behavior
}
