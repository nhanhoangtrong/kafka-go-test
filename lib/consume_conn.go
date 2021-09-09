package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	topic := "test1"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		panic(err)
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := conn.ReadBatch(10e3, 1e6) // 10KB to 1MB

	// Passing no signals to Notify means that
	// all signals will be sent to the channel.

	b := make([]byte, 10e3)
	i := 0
	for {
		i++
		_, err := batch.Read(b)
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Println(string(b))
	}
	fmt.Printf("Run %d\n", i)

	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}
	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
}
