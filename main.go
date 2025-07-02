package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	brokerAddress := "localhost:9092"
	topic := "hello-kafka"

	// Создаём топик (если не существует)
	conn, err := kafka.Dial("tcp", brokerAddress)
	if err != nil {
		log.Fatalf("failed to connect to kafka: %v", err)
	}
	defer conn.Close()
	_ = conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})

	// Producer: пишем сообщение
	writer := kafka.Writer{
		Addr:     kafka.TCP(brokerAddress),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	err = writer.WriteMessages(context.Background(),
		kafka.Message{Value: []byte("hello from kafka-go")},
	)
	if err != nil {
		log.Fatalf("failed to write message: %v", err)
	}
	writer.Close()

	// Consumer: читаем сообщение
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{brokerAddress},
		Topic:     topic,
		Partition: 0,
		MinBytes:  1,
		MaxBytes:  10e6,
	})
	defer reader.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	msg, err := reader.ReadMessage(ctx)
	if err != nil {
		log.Fatalf("failed to read message: %v", err)
	}
	if string(msg.Value) == "hello from kafka-go" {
		fmt.Println("Kafka is working!")
	} else {
		fmt.Println("Kafka responded, but message was:", string(msg.Value))
	}
}