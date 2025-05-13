package main

import (
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/panzerstadt/go-image-pipeline/configs"
)

func get_consumer() sarama.ConsumerGroup {
	config := sarama.NewConfig()
	config.Version = configs.SaramaVersion
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true // remove immediately when i receive it
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	log.Printf("Connecting to brokers: %v\n", configs.Brokers)
	client, err := sarama.NewConsumerGroup(configs.Brokers, configs.TestConsumerGroup, config)
	if err != nil {
		log.Fatalf("unable to create kafka consumer group: %v", err)
	}

	return client
}
