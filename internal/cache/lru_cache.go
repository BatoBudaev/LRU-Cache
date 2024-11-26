package cache

import (
	"container/list"
	"fmt"
	"strings"
	"time"
)

type ICache interface {
	Cap() int
	Len() int
	Clear() // удаляет все ключи
	Add(key, value any)
	AddWithTTL(key, value any, ttl time.Duration) // добавляет ключ со сроком жизни ttl
	Get(key any) (value any, ok bool)
	Remove(key any)
}

type LRUCache struct {
	capacity int
	itemList *list.List
	itemsMap map[any]*list.Element
}

type item struct {
	key   any
	value any
	ttl   time.Duration
}

func New(capacity int) *LRUCache {
	return &LRUCache{capacity: capacity, itemList: list.New(), itemsMap: make(map[any]*list.Element)}
}

func (c *LRUCache) Cap() int {
	return c.capacity
}

func (c *LRUCache) Len() int {
	if len(c.itemsMap) == c.itemList.Len() {
		return len(c.itemsMap)
	}

	return -1
}

func (c *LRUCache) Add(key, value any) {
	if node, ok := c.itemsMap[key]; ok {
		c.itemList.MoveToFront(node)
		return
	}

	if c.capacity == len(c.itemsMap) {
		delete(c.itemsMap, c.itemList.Back().Value.(item).key)
		c.itemList.Remove(c.itemList.Back())
	}

	node := c.itemList.PushFront(item{key: key, value: value})
	c.itemsMap[key] = node
}

func (c *LRUCache) String() string {
	sb := strings.Builder{}

	for n := c.itemList.Front(); n != nil; n = n.Next() {
		sb.WriteString(fmt.Sprintf("[%v: %v]\n", n.Value.(item).key, n.Value.(item).value))
	}

	return sb.String()
}
