package shared

import (
	"log"

	"github.com/panzerstadt/go-image-pipeline/pb"
	"google.golang.org/protobuf/proto"
)

func MakeResizeTask(id string, dir string, path string) []byte {
	task := &pb.ResizeTask{
		ImageId:     id,
		Path:        path,
		Dir:         dir,
		Resize:      true,
		Progressive: true,
	}
	data, err := proto.Marshal(task)
	if err != nil {
		log.Fatal("can't open file")
	}
	return data
}
