//package main
//
//import (
//	"bufio"
//	"container/list"
//	"io"
//
//	"fmt"
//	"os"
//)
//
//type Cache1 struct {
//	maxCacheSize     int64
//	currentSize      int64
//	items            map[string]*list.Element
//	doublyLinkedList *list.List
//	threshold        int64
//}
//
//type Node1 struct {
//	key         string
//	fileContent []byte
//}
//
//func LruCache(sizeinMb int64) *Cache1 {
//
//	sizeInBytes := sizeinMb << 20
//
//	return &Cache1{
//		maxCacheSize:     sizeInBytes,
//		currentSize:      0,
//		items:            make(map[string]*list.Element),
//		doublyLinkedList: list.New(),
//		threshold:        sizeInBytes * 60 / 100,
//	}
//
//}
//func (cache *Cache1) GetCache(keyName string) []byte {
//	if node, found := cache.items[keyName]; found {
//		cache.doublyLinkedList.MoveToFront(node)
//		return node.Value.(*Node1).fileContent
//	}
//	return nil
//}
//
//func (cache *Cache1) PutCache(keyName string, file *os.File) bool {
//
//	data, err := io.ReadAll(file)
//
//	if err != nil {
//		fmt.Println("Error reading file:", err)
//		return false
//	}
//
//	fileSize := int64(len(data))
//
//	if node, found := cache.items[keyName]; found {
//		cache.doublyLinkedList.MoveToFront(node)
//		fileContent, err := io.ReadAll(file)
//		if err != nil {
//			fmt.Println("Error reading file:", err)
//			return false
//		}
//		node.Value.(*Node1).fileContent = fileContent
//		return true
//	}
//
//	if cache.currentSize+fileSize >= cache.maxCacheSize {
//		isEvicted := cache.Evict()
//		if isEvicted {
//			fmt.Println("Evicting from cache")
//		}
//	}
//
//	node := &Node1{
//		key:         keyName,
//		fileContent: data,
//	}
//
//	elem := cache.doublyLinkedList.PushFront(node)
//	cache.items[keyName] = elem
//	cache.currentSize += fileSize
//	return true
//}
//
//func (cache *Cache1) Evict() bool {
//
//	if cache.currentSize <= cache.threshold {
//		return false // No eviction needed, the cache is already at or below 60% usage
//	}
//
//	for cache.currentSize > cache.threshold {
//		back := cache.doublyLinkedList.Back()
//		if back == nil {
//			break // No more items to evict
//		}
//
//		node := back.Value.(*Node1)
//		cache.currentSize = cache.currentSize - int64(len(node.fileContent))
//		delete(cache.items, node.key)
//		cache.doublyLinkedList.Remove(back)
//		fmt.Println("Evicting from cache" + node.key)
//	}
//	return true
//}
//
//func main() {
//	cache := LruCache(1024)
//
//	fmt.Print("Enter the file path: ")
//	reader := bufio.NewReader(os.Stdin)
//	filePath, _ := reader.ReadString('\n')
//	filePath = filePath[:len(filePath)-1]
//
//	file, err := os.Open(filePath)
//	if err != nil {
//		fmt.Println("Error opening file:", err)
//		return
//	}
//	defer file.Close()
//
//	cache.PutCache("File1", file)
//
//	content := cache.GetCache("file1")
//
//	if content != nil {
//		fmt.Println("File content retrieved from cache:")
//		fmt.Println(string(content)) // Print the content as a string (assuming text file)
//	} else {
//		fmt.Println("File not found in cache.")
//	}
//}

package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
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

	// If the file is already in cache, update its content and move it to the front
	if node, found := cache.items[keyName]; found {
		cache.doublyLinkedList.MoveToFront(node)
		node.Value.(*Node1).fileContent = data
		return true
	}

	// Check if there is enough space for the new file
	if cache.currentSize+fileSize >= cache.maxCacheSize {
		isEvicted := cache.Evict()
		if isEvicted {
			fmt.Println("Evicting from cache")
		}
	}

	// Create a new node and add it to the cache
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
		return false // No eviction needed
	}

	// Evict files until the cache size is below the threshold
	for cache.currentSize > cache.threshold {
		back := cache.doublyLinkedList.Back()
		if back == nil {
			break // No more items to evict
		}

		node := back.Value.(*Node1)
		cache.currentSize -= int64(len(node.fileContent))
		delete(cache.items, node.key)
		cache.doublyLinkedList.Remove(back)
		fmt.Println("Evicted from cache:", node.key)
	}
	return true
}

func main() {
	cache := LruCache(1024)

	// Get the file path from user input
	fmt.Print("Enter the file path: ")
	reader := bufio.NewReader(os.Stdin)
	filePath, _ := reader.ReadString('\n')
	filePath = strings.TrimSpace(filePath) // Trim the new line character

	startTime := time.Now()

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Add the file to the cache
	cache.PutCache("File1", file)

	// Retrieve the file content from the cache
	content := cache.GetCache("File1")

	if content != nil {
		fmt.Println("File content retrieved from cache:")
		fmt.Println(string(content))
		// Print the content as a string (assuming text file)

		elapsedPut := time.Since(startTime)

		fmt.Println("Time elapsed:", elapsedPut)
	} else {
		fmt.Println("File not found in cache.")
	}

}
