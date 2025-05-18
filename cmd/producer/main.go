package main

import (
	"fmt"
	"log"
	"time"

	"github.com/panzerstadt/go-image-pipeline/configs"
)

// for now its ok to do double work
func main() {
	// get producer
	producer := get_producer()
	defer producer.Close()

	// 1. cron to find files in the folder
	for {
		filenames := scanFolder(configs.InputDir)
		// 2. loop through all files in folder and prepare message
		for id, filename := range filenames {
			msg := prepareMessageForTopic(configs.TopicImageJobs, fmt.Sprintf("%v", id), configs.InputDir, filename)
			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				log.Printf("Failed to send message: %v", err)
				if err.Error() == "kafka server: Request was for a topic or partition that does not exist on this broker" {
					panic("kafka setup error")
				}
				if err.Error() == "circuit breaker is open" {
					panic("something's wrong with the kafka setup")
				}
			} else {
				fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
			}
			moveToIntermediateDirAndMarkSent(configs.InputDir, filename)
		}

		time.Sleep(time.Second * 10)
	}
}

/**
and offset is a marker that the consumer has processed that message (consumed)
this decouples "consuming" from "deleting" because we can then delete the messages asynchronously,
which also means we can redo the messages if something goes wrong outside of the scope of the
event stream (e.g. 3 days later we noticed there were bugs that require reprocessing)
*/
