package main

import (
	""
	"bufio"
	"fmt"
	"os"
)

func main() {
	cache := learnings.LruCache(1024)

	fmt.Print("Enter the file path: ")
	reader := bufio.NewReader(os.Stdin)
	filePath, _ := reader.ReadString('\n')
	filePath = filePath[:len(filePath)-1]

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	cache.PutCache("File1", file)

	content := cache.GetCache("file1")

	if content != nil {
		fmt.Println("File content retrieved from cache:")
		fmt.Println(string(content)) // Print the content as a string (assuming text file)
	} else {
		fmt.Println("File not found in cache.")
	}
}
