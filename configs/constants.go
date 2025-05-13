package configs

import "github.com/IBM/sarama"

var Brokers = []string{"localhost:9092"}
var SaramaVersion = sarama.V3_3_0_0

var TestTopic = "test-topic"

var TestConsumerGroup = "test-consumer-group"

var InputDir = "./inputs"
var IntermediateDir = "./intermediate"
var OutputDir = "./outputs"
var ProducedMarkSuffix = ".sent"
