package main

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/panzerstadt/go-image-pipeline/configs"
)

func main() {
	// get producer
	producer := get_producer()
	defer producer.Close()

	// make message
	msg := &sarama.ProducerMessage{
		Topic: configs.TestTopic,
		Value: sarama.StringEncoder("Hello World!"),
	}

	// send message
	// TODO: learn wth an offset is
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	} else {
		fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
	}
}
