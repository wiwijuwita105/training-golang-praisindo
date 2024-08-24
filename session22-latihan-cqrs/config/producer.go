package config

import (
	"fmt"
	"github.com/IBM/sarama"
)

func SetupProducer() (sarama.SyncProducer, error) {
	configSarama := sarama.NewConfig()
	configSarama.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{KafkaBrokerAddress}, configSarama)
	if err != nil {
		return nil, fmt.Errorf("failed to setup producer: %w", err)
	}
	return producer, nil
}
