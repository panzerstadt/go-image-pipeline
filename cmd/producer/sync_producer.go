package main

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/panzerstadt/go-image-pipeline/configs"
)

// https://www.tencentcloud.com/document/product/597/60360
func get_producer() sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Version = sarama.V3_3_0_0
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(configs.Brokers, config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}

	return producer
}
