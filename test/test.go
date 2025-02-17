package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	a := sync.Map{}
	b := sync.Map{}
	c := sync.Map{}
	d := make(chan int)
	a.Store(d, nil)
	b.Store(d, nil)
	c.Store(d, nil)
	go func() {
		for {
			time.Sleep(1 * time.Second)
			d <- 1
		}
	}()

	go func() {
		for {
			a.Range(func(key, value interface{}) bool {
				select {
				case x, ok := <-key.(chan int):
					if ok {
						fmt.Println("aaa", x)
					}
				}
				return true
			})
		}
	}()

	go func() {
		for {
			b.Range(func(key, value interface{}) bool {
				select {
				case x, ok := <-key.(chan int):
					if ok {
						fmt.Println("bbb", x)
					}
				}
				return true
			})
		}
	}()

	go func() {
		for {
			c.Range(func(key, value interface{}) bool {
				select {
				case x, ok := <-key.(chan int):
					if ok {
						fmt.Println("ccc", x)
					}
				}
				return true
			})
		}
	}()
	for {
	}
}
