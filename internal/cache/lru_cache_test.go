package cache

import (
	"sync"
	"testing"
	"time"
)

func TestLRUCache_Add(t *testing.T) {
	cache := New(3)

	cache.Add("key1", "value1")
	cache.Add("key2", "value2")
	cache.Add("key3", "value3")

	cache.Add("key2", "value22") // Замена по ключу
	cache.Add("key4", "value4")  // Вытеснение

	if cache.Len() != 3 {
		t.Errorf("Неверная длина: %d != 3", cache.Len())
	}

	if cache.Cap() != 3 {
		t.Errorf("Неверная вместимость: %d != 3", cache.Cap())
	}

	if val, ok := cache.Get("key2"); !ok || val != "value22" {
		t.Errorf("Неверное значение value22 != %v", val)
	}

	if _, ok := cache.Get("key1"); ok {
		t.Error("Ключ не удалился")
	}
}

func TestLRUCache_AddWithTTL(t *testing.T) {
	cache := New(3)

	cache.AddWithTTL("key1", "value1", 1*time.Second)
	cache.AddWithTTL("key2", "value2", 0) // 0 - без ttl
	cache.Add("key3", "value3")

	time.Sleep(2 * time.Second)

	if cache.Len() != 2 {
		t.Errorf("Неверная длина: %d != 2", cache.Len())
	}

	if _, ok := cache.Get("key1"); ok {
		t.Error("Ключ не удалился")
	}
}

func TestLRUCache_Clear(t *testing.T) {
	cache := New(3)

	cache.Add("key1", "value1")
	cache.Add("key2", "value2")
	cache.Add("key3", "value3")

	cache.Clear()

	if cache.Len() != 0 {
		t.Errorf("Неверная длина: %d != 0", cache.Len())
	}

	if cache.Cap() != 3 {
		t.Errorf("Неверная вместимость: %d != 3", cache.Cap())
	}
}

func TestLRUCache_Remove(t *testing.T) {
	cache := New(3)

	cache.Add("key1", "value1")
	cache.Add("key2", "value2")
	cache.Add("key3", "value3")

	cache.Remove("key2")

	if cache.Len() != 2 {
		t.Errorf("Неверная длина: %d != 2", cache.Len())
	}

	if _, ok := cache.Get("key2"); ok {
		t.Error("Ключ не удалился")
	}

	cache.Remove("key1")
	cache.Remove("key3")

	if cache.Len() != 0 {
		t.Errorf("Неверная длина: %d != 0", cache.Len())
	}
}

func TestLRUCache_Goroutines(t *testing.T) {
	cache := New(10)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cache.Add(i, i*i)
		}(i)
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cache.Get(i)
		}(i)
	}

	wg.Wait()

	if cache.Len() != 10 {
		t.Errorf("Неверная длина: %d != 0", cache.Len())
	}
}
