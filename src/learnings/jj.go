package learnings

import (
	"container/list"
	"fmt"
	"io"
	"os"
)

type Cache1 struct {
	maxCacheSize     int64
	currentSize      int64
	items            map[string]*list.Element
	doublyLinkedList *list.List
	threshold        int64
}

type Node1 struct {
	key         string
	fileContent []byte
}

func LruCache(sizeinMb int64) *Cache1 {

	sizeInBytes := sizeinMb << 20

	return &Cache1{
		maxCacheSize:     sizeInBytes,
		currentSize:      0,
		items:            make(map[string]*list.Element),
		doublyLinkedList: list.New(),
		threshold:        sizeInBytes * 60 / 100,
	}

}
func (cache *Cache1) GetCache(keyName string) []byte {
	if node, found := cache.items[keyName]; found {
		cache.doublyLinkedList.MoveToFront(node)
		return node.Value.(*Node1).fileContent
	}
	return nil
}

func (cache *Cache1) PutCache(keyName string, file *os.File) bool {

	data, err := io.ReadAll(file)

	if err != nil {
		fmt.Println("Error reading file:", err)
		return false
	}

	fileSize := int64(len(data))

	if node, found := cache.items[keyName]; found {
		cache.doublyLinkedList.MoveToFront(node)
		fileContent, err := io.ReadAll(file)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return false
		}
		node.Value.(*Node1).fileContent = fileContent
		return true
	}

	if cache.currentSize+fileSize >= cache.maxCacheSize {
		isEvicted := cache.Evict()
		if isEvicted {
			fmt.Println("Evicting from cache")
		}
	}

	node := &Node1{
		key:         keyName,
		fileContent: data,
	}

	elem := cache.doublyLinkedList.PushFront(node)
	cache.items[keyName] = elem
	cache.currentSize += fileSize
	return true
}

func (cache *Cache1) Evict() bool {

	if cache.currentSize <= cache.threshold {
		return false // No eviction needed, the cache is already at or below 60% usage
	}

	for cache.currentSize > cache.threshold {
		back := cache.doublyLinkedList.Back()
		if back == nil {
			break // No more items to evict
		}

		node := back.Value.(*Node1)
		cache.currentSize = cache.currentSize - int64(len(node.fileContent))
		delete(cache.items, node.key)
		cache.doublyLinkedList.Remove(back)
		fmt.Println("Evicting from cache" + node.key)
	}
	return true
}
