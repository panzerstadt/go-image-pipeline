package main

import (
	"log"

	"github.com/panzerstadt/go-image-pipeline/pb"
	"google.golang.org/protobuf/proto"
)

func sendResizeTask() []byte {
	task := &pb.ResizeTask{
		ImageId:     "123",
		Path:        "/source/test.JPG",
		Resize:      true,
		Progressive: true,
	}
	data, err := proto.Marshal(task)
	if err != nil {
		log.Fatal("can't open file")
	}
	return data
}
