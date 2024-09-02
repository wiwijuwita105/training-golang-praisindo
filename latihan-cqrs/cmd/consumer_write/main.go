package main

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"latihan-cqrs/entity"
	"log"
	"log/slog"
	"sync"
)

type UserData map[string][]entity.UserInput

type UserDataStore struct {
	data UserData
	mu   sync.RWMutex
}

func (ns *UserDataStore) Add(key string, user entity.UserInput) error {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	ns.data[key] = append(ns.data[key], user)
}

type Consumer struct {
	store *UserDataStore
}

func (*Consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (*Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (consumer *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		userID := string(msg.Key)
		var user entity.UserInput
		err := json.Unmarshal(msg.Value, &user)
		if err != nil {
			log.Printf("failed to unmarshal notification: %v", err)
			continue
		}
		slog.Info("Consuming notification and adding it to storage", slog.Any("notification", notification))
		consumer.store.Add(userID, user)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func setupConsumerGroup(ctx context.Context, store *UserDataStore) {
	config := sarama.NewConfig()

	consumerGroup, err := sarama.NewConsumerGroup([]string{entity.KafkaBrokerAddress}, entity.WriteGroupID, config)
	if err != nil {
		log.Printf("initialization error: %v", err)
	}
	defer consumerGroup.Close()

	consumer := &Consumer{
		store: store,
	}

	for {
		err = consumerGroup.Consume(ctx, []string{entity.DefaultTopic}, consumer)
		if err != nil {
			log.Printf("error from consumer: %v", err)
		}
		if ctx.Err() != nil {
			return
		}
	}
}

func main() {
	store := &UserDataStore{
		data: make(UserData),
	}

	ctx, cancel := context.WithCancel(context.Background())
	go setupConsumerGroup(ctx, store)
	defer cancel()
}
