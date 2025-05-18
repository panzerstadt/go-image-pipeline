package main

import (
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/panzerstadt/go-image-pipeline/configs"
)

func make_consumer_for_consumer_group(consumerGroup string) sarama.ConsumerGroup {
	config := sarama.NewConfig()
	config.Version = configs.SaramaVersion
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true // remove immediately when i receive it
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	log.Printf("Connecting to consumerGroup: '%s' for brokers: '%v'\n", consumerGroup, configs.Brokers)
	client, err := sarama.NewConsumerGroup(configs.Brokers, consumerGroup, config)
	if err != nil {
		log.Fatalf("unable to create kafka consumer group: %v", err)
	}

	return client
}
