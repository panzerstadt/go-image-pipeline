package main

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/panzerstadt/go-image-pipeline/configs"
)

func create_topic() {
	config := sarama.NewConfig()
	config.Version = configs.SaramaVersion

	admin, err := sarama.NewClusterAdmin(configs.Brokers, config)
	if err != nil {
		log.Fatal("Error creating cluster admin:", err)
	}
	defer admin.Close()

	topicName := configs.TestTopic
	// redundancy for brokers
	// 2 brokers + replication factor 2 = 2 copies per partition, where
	// one leader, one follower
	detail := &sarama.TopicDetail{
		NumPartitions:     2,
		ReplicationFactor: 1, // i only gots one broker rn
	}

	err = admin.CreateTopic(topicName, detail, false)
	if err != nil {
		log.Fatalf("Error creating topic: %v", err)
	}

	log.Printf("Topic %q created.", topicName)
}

func remove_topic() {
	config := sarama.NewConfig()
	config.Version = configs.SaramaVersion

	admin, err := sarama.NewClusterAdmin(configs.Brokers, config)
	if err != nil {
		log.Fatal("Error creating cluster admin:", err)
	}
	defer admin.Close()

	topicName := configs.TestTopic
	err = admin.DeleteTopic(topicName)
	if err != nil {
		log.Fatalf("Error deleting topic: %v", err)
	}

	log.Printf("Topic %q deleted.", topicName)
}
