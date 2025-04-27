package main

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/panzerstadt/go-image-pipeline/pb"
	"google.golang.org/protobuf/proto"
)

func receive(msg *sarama.ConsumerMessage) pb.ResizeTask {
	var task pb.ResizeTask
	err := proto.Unmarshal(msg.Value, &task)
	if err != nil {
		log.Fatal("error unmarshaling the protobuf")
	}

	return task
}
