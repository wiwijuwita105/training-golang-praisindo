package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"log/slog"
	"os"
	"session22-latihan-cqrs/config"
	"session22-latihan-cqrs/entity"
)

var db *gorm.DB

type ConsumerRead struct{}

func (ConsumerRead) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (ConsumerRead) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h ConsumerRead) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		slog.Info("Received message", slog.Any("string(msg.Value)", string(msg.Value)))

		var userInput entity.UserInput
		_ = json.Unmarshal(msg.Value, &userInput)
		slog.Info("Unmarshall results", slog.Any("user", userInput))

		var user = entity.User{
			Name: userInput.Name,
		}
		log.Println(user)

		if db == nil {
			slog.Error("Database connection is nil")
			return errors.New("database connection is nil")
		}

		query := "INSERT INTO users (name) VALUES ($1)"
		if err := db.Exec(query, user.Name).Error; err != nil {
			log.Printf("Error deleting user: %v\n", err)
			return err
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}

func main() {
	var err error
	db, err = gorm.Open(postgres.Open(config.DBReadConn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		slog.Error("error opening write database", slog.String("dsn", config.DBReadConn), slog.Any("err", err))
		os.Exit(1)
	}

	log.Println("starting migration in db-read")
	if err := db.AutoMigrate(entity.User{}); err != nil {
		log.Println(err)
	}
	log.Println("finish migration in db-read")

	configConsumer := sarama.NewConfig()
	consumerGroup, err := sarama.NewConsumerGroup(
		[]string{config.KafkaBrokerAddress},
		config.ReadGroupID,
		configConsumer,
	)
	if err != nil {
		log.Printf("initialization error: %v", err)
		os.Exit(1)
	}

	ctx := context.Background()
	slog.InfoContext(ctx, "Start consuming from topic", slog.Any("topic", config.UserTopic))
	for {
		err = consumerGroup.Consume(
			ctx,
			[]string{config.UserTopic},
			ConsumerRead{},
		)
		if err != nil {
			slog.Error("error when call consumerGroup.Consume:", slog.Any("error", err))
		}
	}
}
