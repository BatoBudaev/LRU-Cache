package cache

import (
	"container/list"
	"fmt"
	"strings"
	"sync"
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
	mu        sync.RWMutex
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
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.capacity
}

func (c *LRUCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.removeExpired()

	if len(c.itemsMap) == c.itemsList.Len() {
		return len(c.itemsMap)
	}

	return -1
}

func (c *LRUCache) Add(key, value any) {
	c.addItem(key, value, 0, false)
}

func (c *LRUCache) AddWithTTL(key, value any, ttl time.Duration) {
	c.addItem(key, value, ttl, true)
}

func (c *LRUCache) addItem(key, value any, ttl time.Duration, hasTTL bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, ok := c.itemsMap[key]; ok {
		cacheItem := node.Value.(*item)
		cacheItem.value = value
		if hasTTL {
			cacheItem.Expiration = time.Now().Add(ttl)
		}

		c.itemsList.MoveToFront(node)
		return
	}

	if c.capacity == len(c.itemsMap) {
		backNode := c.itemsList.Back()
		if backNode != nil {
			delete(c.itemsMap, backNode.Value.(*item).key)
			c.itemsList.Remove(backNode)
		}
	}

	var expiration time.Time
	if hasTTL && ttl != 0 {
		expiration = time.Now().Add(ttl)
	}

	node := c.itemsList.PushFront(&item{key: key, value: value, Expiration: expiration})
	c.itemsMap[key] = node
}

func (c *LRUCache) Get(key any) (value any, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.removeExpired()
	if node, ok2 := c.itemsMap[key]; ok2 {
		c.itemsList.MoveToFront(node)

		return node.Value.(*item).value, true
	}

	return nil, false
}

func (c *LRUCache) String() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.removeExpired()
	sb := strings.Builder{}

	for n := c.itemsList.Front(); n != nil; n = n.Next() {
		cacheItem := n.Value.(*item)
		sb.WriteString(fmt.Sprintf("[%v: %v]\n", cacheItem.key, cacheItem.value))
	}

	return sb.String()
}

func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k := range c.itemsMap {
		delete(c.itemsMap, k)
	}

	for n := c.itemsList.Back(); n != nil; n = c.itemsList.Back() {
		c.itemsList.Remove(n)
	}
}

func (c *LRUCache) Remove(key any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, ok := c.itemsMap[key]; ok {
		delete(c.itemsMap, key)
		c.itemsList.Remove(node)
	}
}

func (c *LRUCache) removeExpired() {
	now := time.Now()
	for node := c.itemsList.Back(); node != nil; {
		cacheItem := node.Value.(*item)

		if !cacheItem.Expiration.IsZero() && now.After(cacheItem.Expiration) {
			prevNode := node.Prev()
			delete(c.itemsMap, cacheItem.key)
			c.itemsList.Remove(node)
			node = prevNode
		} else {
			node = node.Prev()
		}
	}
}
