package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/IBM/sarama"
	"github.com/panzerstadt/go-image-pipeline/configs"
	"github.com/panzerstadt/go-image-pipeline/shared"
)

func scanFolder(dir string) (filenames []string) {
	fmt.Println("scanning directory: " + dir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal("can't open directory: " + dir)
	}

	filenames = make([]string, 0)
	for _, entry := range entries {
		if entry.Type().IsRegular() {
			filename := entry.Name()

			if !strings.Contains(filename, configs.ProducedMarkSuffix) && !strings.Contains(filename, ".DS_Store") {
				filenames = append(filenames, filename)
			}
		}
	}

	return filenames
}

func prepareMessageForTopic(topic string, id string, dir string, filename string) *sarama.ProducerMessage {
	fullpath := path.Join(dir, filename)
	_, err := os.Stat(fullpath)
	if err != nil {
		log.Fatal(fullpath + " does not exist!")
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		// Key: sarama.StringEncoder(), // this can be used to set which partitions the message goes to (e.g. by user)
		Value: sarama.ByteEncoder(shared.MakeResizeTask(id, dir, filename)),
	}

	return msg
}

func moveToIntermediateDirAndMarkSent(dir string, filename string) {
	srcPath := path.Join(dir, filename)
	_, err := os.Stat(srcPath)
	if err != nil {
		log.Fatal(srcPath + " does not exist!")
	}

	midPath := path.Join(configs.IntermediateDir, filename)
	err = shared.FileCopy(srcPath, midPath)
	fmt.Println(srcPath + " copied to " + midPath)
	if err != nil {
		log.Fatal(err)
	}

	newPath := srcPath + configs.ProducedMarkSuffix
	os.Rename(srcPath, newPath)
	fmt.Println(srcPath + " renamed to " + newPath)
}
