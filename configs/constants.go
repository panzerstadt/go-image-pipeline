package configs

import "github.com/IBM/sarama"

var Brokers = []string{"localhost:9092"}
var SaramaVersion = sarama.V3_3_0_0

var TopicImageJobs = "image-jobs"

var ConsumerGroupImageProcessor = "consumer-image-processor"    // does the actual image processing
var ConsumerGroupNotifications = "consumer-image-notifications" // pings me (sends me email maybe)
var ConsumerGroupAnalytics = "consumer-analytics"               // counts number of images done, how long it took tc

var InputDir = "./inputs"
var IntermediateDir = "./intermediate"
var OutputDir = "./outputs"
var ProducedMarkSuffix = ".sent"
var NotificationsFile = "notifications.txt"
