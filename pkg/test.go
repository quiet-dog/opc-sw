package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	go func() {
		select {
		case <-ctx.Done():
			{
				fmt.Println("context done1")
			}
		}

	}()

	go func() {
		select {
		case <-ctx.Done():
			{
				fmt.Println("context done2")
			}
		}

	}()

	time.Sleep(10 * time.Second)
}
