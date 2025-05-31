package shared

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/panzerstadt/go-image-pipeline/pb"
	"google.golang.org/protobuf/proto"
)

func ReceiveResizeTask(msg *sarama.ConsumerMessage) *pb.ResizeTask {
	var task pb.ResizeTask
	err := proto.Unmarshal(msg.Value, &task)
	if err != nil {
		log.Fatal("error unmarshaling resize task")
	}

	return &task
}

func ReceiveSyncTask(msg *sarama.ConsumerMessage) *pb.SyncTask {
	var task pb.SyncTask
	err := proto.Unmarshal(msg.Value, &task)
	if err != nil {
		log.Fatal("error unmarshaling sync task")
	}
	return &task
}
