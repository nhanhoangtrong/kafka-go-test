package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	kafka "github.com/segmentio/kafka-go"
)

func getKafkaReader(kafkaURL, topic string, partition int) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		Partition: partition,
		// GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e2, // 1KB
		MaxBytes: 10e6, // 10MB
	})
}

type UUIDObject struct {
	Key  int    `json:"key"`
	Uuid string `json:"uuid"`
}

// Stage 1: Consume 5 arrays of bytes
func consume(values [5][]byte) <-chan []byte {
	out := make(chan []byte)
	go func() {
		for i := 0; i < 5; i++ {
			out <- values[i]
		}
		close(out)
	}()
	return out
}

// Stage 2: parse array of bytes
func parse(in <-chan []byte) <-chan UUIDObject {
	out := make(chan UUIDObject)
	go func() {
		for value := range in {
			var obj UUIDObject
			// Parse
			err := json.Unmarshal(value, &obj)
			if err == nil {
				out <- obj
			}
		}
		close(out)
	}()
	return out
}

// Stage 3: print parsed value
func print(in <-chan UUIDObject) {
	for obj := range in {
		fmt.Printf("Key: %d, UUID: %s\n", obj.Key, obj.Uuid)
	}
}

func main() {
	// get kafka reader
	kafkaURL := "localhost:9092"
	topic := "test1"
	// groupID := "default"
	partition := 0

	reader := getKafkaReader(kafkaURL, topic, partition)
	defer reader.Close()

	fmt.Println("start comsuming... !!")
	for {
		var values [5][]byte
		for i := 0; i < 5; i++ {
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Fatalln("error read", err)
			}
			values[i] = msg.Value
		}
		out := consume(values)

		// Distribute the parse stage across some goroutines
		c1 := parse(out)
		c2 := parse(out)

		// Finall, merge and print
		print(merge(c1, c2))
	}
}
