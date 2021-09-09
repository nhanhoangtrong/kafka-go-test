package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		defer cancel()
		time.Sleep(3 * time.Second)
		fmt.Println("Cancelled")
	}()

	go doSomething(ctx, 3)
	go doAnother(ctx, 8)

	fmt.Scanln()
}

func doSomething(ctx context.Context, timeSleep time.Duration) {
	fmt.Println("start doSomething")
	select {
	case <-time.After(timeSleep):
		fmt.Println("Finished doSomething")
	case <-ctx.Done():
		fmt.Println("Cancelled doSomething")
	}
}

func doAnother(ctx context.Context, n int) {
	for i := 0; i < n; i++ {
		if ctx.Err() != nil {
			fmt.Println("Leaving doAnother")
			return
		}

		time.Sleep(1 * time.Second)
		fmt.Println("Iteration", i)
	}
}
