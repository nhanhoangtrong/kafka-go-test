package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	c := make(chan string)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	go doSomething(ctx, 5, c)

	select {
	case <-ctx.Done():
		fmt.Printf("Context cancelled, %v\n", ctx.Err())

	case res := <-c:
		fmt.Printf("Receive: %v\n", res)
	}
}

func doSomething(ctx context.Context, timeSleep time.Duration, ch chan string) {
	fmt.Println("Sleeping...")
	time.Sleep(timeSleep * time.Second)
	fmt.Println("Waking up...!")
	ch <- "Done"
}
