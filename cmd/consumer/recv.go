package main

import (
	"fmt"

	"github.com/panzerstadt/go-image-pipeline/pb"
)

func process() {
	var task pb.ResizeTask
	// err := proto.Unmarshal(msg.Value, &task)
	// if err != nil {
	// 	log.Fatal("error unmarshaling the protobuf")
	// }

	fmt.Println(task)
}
