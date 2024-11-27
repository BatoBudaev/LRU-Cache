package main

import (
	"LRU-Cache/internal/cache"
	"fmt"
)

func main() {
	c := cache.New(3)
	fmt.Printf("Cap: %d Len:%d\n", c.Cap(), c.Len())
	c.Add("asd", 10)
	c.Add("asd", 10)
	fmt.Printf("Cap: %d Len:%d\n", c.Cap(), c.Len())
	c.Add(12, 1)
	c.Add("qwe", "asd")
	c.Add("ffff", "123")
	c.Add(12, 1)
	fmt.Printf("Cap: %d Len:%d\n", c.Cap(), c.Len())
	//c.Clear()
	fmt.Println("________")
	fmt.Println(c)
	c.Remove("ffff")
	fmt.Printf("Cap: %d Len:%d\n", c.Cap(), c.Len())
	fmt.Println(c)
	c.Add(1, 2)
	fmt.Println(c)
	fmt.Printf("Cap: %d Len:%d\n", c.Cap(), c.Len())
}
