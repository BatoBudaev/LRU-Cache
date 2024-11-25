package main

import (
	"container/list"
	"time"
)

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

type ICache interface {
	Cap() int
	Len() int
	Clear() // удаляет все ключи
	Add(key, value any)
	AddWithTTL(key, value any, ttl time.Duration) // добавляет ключ со сроком жизни ttl
	Get(key any) (value any, ok bool)
	Remove(key any)
}

func main() {
	cache := New(10)

}
