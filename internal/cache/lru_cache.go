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
	capacity  int
	itemsList *list.List
	itemsMap  map[any]*list.Element
}

type item struct {
	key        any
	value      any
	Expiration time.Time
}

func New(capacity int) *LRUCache {
	return &LRUCache{capacity: capacity, itemsList: list.New(), itemsMap: make(map[any]*list.Element)}
}

func (c *LRUCache) Cap() int {
	return c.capacity
}

func (c *LRUCache) Len() int {
	if len(c.itemsMap) == c.itemsList.Len() {
		return len(c.itemsMap)
	}

	return -1
}

func (c *LRUCache) Add(key, value any) {
	if node, ok := c.itemsMap[key]; ok {
		c.itemsList.MoveToFront(node)
		return
	}

	if c.capacity == len(c.itemsMap) {
		delete(c.itemsMap, c.itemsList.Back().Value.(item).key)
		c.itemsList.Remove(c.itemsList.Back())
	}

	node := c.itemsList.PushFront(item{key: key, value: value})
	c.itemsMap[key] = node
}

func (c *LRUCache) Get(key any) (value any, ok bool) {
	if node, ok2 := c.itemsMap[key]; ok2 {
		nodeItem := node.Value.(item)
		if !nodeItem.Expiration.IsZero() && time.Now().After(nodeItem.Expiration) {
			delete(c.itemsMap, key)
			c.itemsList.Remove(node)

			return nil, false
		}

		c.itemsList.MoveToFront(node)
		return node.Value.(item).value, true
	}

	return nil, false
}

func (c *LRUCache) String() string {
	sb := strings.Builder{}

	for n := c.itemsList.Front(); n != nil; n = n.Next() {
		sb.WriteString(fmt.Sprintf("[%v: %v]\n", n.Value.(item).key, n.Value.(item).value))
	}

	return sb.String()
}

func (c *LRUCache) Clear() {
	for k := range c.itemsMap {
		delete(c.itemsMap, k)
	}

	for n := c.itemsList.Back(); n != nil; n = c.itemsList.Back() {
		c.itemsList.Remove(n)
	}
}

func (c *LRUCache) Remove(key any) {
	if node, ok := c.itemsMap[key]; ok {
		delete(c.itemsMap, key)
		c.itemsList.Remove(node)
	}
}

func (c *LRUCache) AddWithTTL(key, value any, ttl time.Duration) {
	if node, ok := c.itemsMap[key]; ok {
		node.Value.(*item).Expiration = time.Now().Add(ttl)
		c.itemsList.MoveToFront(node)
		return
	}

	if c.capacity == len(c.itemsMap) {
		delete(c.itemsMap, c.itemsList.Back().Value.(item).key)
		c.itemsList.Remove(c.itemsList.Back())
	}

	node := c.itemsList.PushFront(item{key: key, value: value, Expiration: time.Now().Add(ttl)})
	c.itemsMap[key] = node
}
