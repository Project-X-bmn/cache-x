package src

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

func (c *Cache) Get(key string) (interface{}, bool) {
	if elem, found := c.items[key]; found {
		c.list.MoveToFront(elem)
		return elem.Value.(*Entry).value, true
	}
	return nil, false
}

func (c *Cache) Put(key string, value interface{}) {
	if elem, found := c.items[key]; found {
		elem.Value.(*Entry).value = value
		c.list.MoveToFront(elem)
		return
	}

	if c.list.Len() >= c.capacity {
		last := c.list.Back()
		if last != nil {
			delete(c.items, last.Value.(*Entry).key)
			c.list.Remove(last)
		}
	}

	entry := &Entry{key: key, value: value}
	elem := c.list.PushFront(entry)
	c.items[key] = elem
}

func (c *Cache) Delete(key string) {
	if elem, found := c.items[key]; found {
		delete(c.items, key)
		c.list.Remove(elem)
	}
}

func (c *Cache) Size() int {
	return c.list.Len()
}

func main() {
	cache := NewLRUCache(3)

	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)

	fmt.Println(cache.Get("a")) // Output: 1, true
	cache.Put("d", 4)           // "b" will be evicted

	_, found := cache.Get("b")
	fmt.Println(found) // Output: false

	cache.Put("e", 5)           // "c" will be evicted
	fmt.Println(cache.Get("c")) // Output: <nil>, false
}
