package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

func main() {
	topic := "test1"
	partition := 0

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// signal.Notify(sigs)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     topic,
		Partition: partition,
		MinBytes:  10e2, // 1KB
		MaxBytes:  10e6, // 10MB
	})

	go func() {
		s := <-sigs
		log.Printf("receive signal %v, gracefully close\n", s)
		reader.Close()
	}()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln("error read", err)
		}
		fmt.Printf("msg: %v\n", string(msg.Value))
	}
}
