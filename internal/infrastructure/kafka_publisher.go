package infrastructure

import (
	"context"
	"encoding/json"
	"event-system/internal/domain"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	brokers  []string
	writers  map[string]*kafka.Writer // кэш writers для топиков
	registry *EventRegistry
}

// NewKafkaPublisher создает publisher с подключением к Kafka cluster
func NewKafkaPublisher(brokers []string, registry *EventRegistry) *KafkaPublisher {
	return &KafkaPublisher{
		brokers:  brokers,
		writers:  make(map[string]*kafka.Writer),
		registry: registry,
	}
}

func (kp *KafkaPublisher) Publish(event *domain.Event) error {
	topic, _, err := kp.registry.ResolveChannel(event.Type)
	if err != nil {
		return fmt.Errorf("failed to resolve channel for event type %s: %w", event.Type, err)
	}

	writer, err := kp.getOrCreateWriter(topic)
	if err != nil {
		return fmt.Errorf("failed to get writer for topic %s: %w", topic, err)
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := kafka.Message{
		Key:   nil,
		Value: data,
		Time:  event.Timestamp,
		Headers: []kafka.Header{
			{Key: "event-type", Value: []byte(event.Type)},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := writer.WriteMessages(ctx, msg); err != nil {
		log.Printf("kafka publish error to topic %s: %v", topic, err)
		return fmt.Errorf("failed to publish to kafka topic %s: %w", topic, err)
	}

	log.Printf("✅ Event %s published to Kafka topic: %s", event.Type, topic)
	return nil
}

// getOrCreateWriter получает существующий writer или создает новый для топика
func (kp *KafkaPublisher) getOrCreateWriter(topic string) (*kafka.Writer, error) {
	// Проверяем кэш
	if writer, exists := kp.writers[topic]; exists {
		return writer, nil
	}

	// Создаем новый writer для топика
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  kp.brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{}, // Распределение по партициям
		// Настройки производительности
		BatchSize:    1, // Для тестирования отправляем сразу
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: 1,
		Async:        false,
	})

	// Сохраняем в кэш
	kp.writers[topic] = writer

	log.Printf("📝 Created new Kafka writer for topic: %s", topic)
	return writer, nil
}

// Close закрывает все writers
func (kp *KafkaPublisher) Close() error {
	for topic, writer := range kp.writers {
		if err := writer.Close(); err != nil {
			log.Printf("error closing writer for topic %s: %v", topic, err)
		}
	}
	return nil
}

// EnsureTopicsExist создает топики если их нет (опционально, для development)
func (kp *KafkaPublisher) EnsureTopicsExist(topics []string) error {
	conn, err := kafka.Dial("tcp", kp.brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	for _, topic := range topics {
		topicConfigs := []kafka.TopicConfig{
			{
				Topic:             topic,
				NumPartitions:     1, // Для начала одна партиция
				ReplicationFactor: 1, // Для development одна реплика
			},
		}

		err = conn.CreateTopics(topicConfigs...)
		if err != nil {
			log.Printf("topic %s already exists or error: %v", topic, err)
		} else {
			log.Printf("✅ Created Kafka topic: %s", topic)
		}
	}
	return nil
}
