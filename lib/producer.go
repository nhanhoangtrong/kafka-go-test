package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
)

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func main() {
	kafkaURL := "localhost:9092"
	topic := "test1"
	writer := newKafkaWriter(kafkaURL, topic)
	defer writer.Close()

	fmt.Println("start producing...!!")

	for i := 0; ; i++ {
		key := fmt.Sprintf("Key-%d", i)
		msg := kafka.Message{
			Key:   []byte(key),
			Value: []byte(fmt.Sprintf("{\"key\": %d, \"uuid\": \"%v\"}", i, uuid.New())),
		}

		err := writer.WriteMessages(context.Background(), msg)
		if err != nil {
			fmt.Println("write error", err)
		} else {
			fmt.Println("produced", key, string(msg.Value))
		}

		time.Sleep(1 * time.Second)
	}
}
