package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/panzerstadt/go-image-pipeline/configs"
	"github.com/panzerstadt/go-image-pipeline/shared"
)

type ImageProcessor struct{}

func (h *ImageProcessor) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ImageProcessor) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ImageProcessor) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		start := time.Now()

		time.Sleep(shared.RandDuration(1, 3))
		fmt.Printf("JOB: Received message: key=%s, value=%s, partition=%d, offset=%d\n", string(msg.Key), strings.ReplaceAll(strings.TrimSpace(string(msg.Value)), "\n", " "), msg.Partition, msg.Offset)
		task := shared.ReceiveResizeTask(msg)
		filename := task.Path

		// 2. resize image and make progressive jpeg
		//      - err: non-retryable
		middlePath := path.Join(configs.IntermediateDir, filename)
		_, err := os.Stat(middlePath)
		if err != nil {
			fmt.Println("JOB: file not found at: ", middlePath)
			sess.MarkMessage(msg, "file not found at: "+middlePath)
			break
		}
		outPath := path.Join(configs.OutputDir, filename)
		// 	/opt/homebrew/bin/convert -strip -interlace Plane -quality 80 -resize 2000x2000 $f $o
		// todo: fully in go?
		cmd := exec.Command("/opt/homebrew/bin/convert", "-strip", "-interlace", "Plane", "-quality", "80", "-resize", "2000x2000", middlePath, outPath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			sess.MarkMessage(msg, "image processing failed for: "+middlePath)
			break
		}
		fmt.Println("JOB: imagemagick process output: " + strings.TrimSpace(string(output)))
		// 2b. call llm to guess and fill in jpeg metadata
		//      - err: retryable
		//      - stub: wait 5 seconds
		// fmt.Println("JOB: sleeping for 5 seconds...")
		time.Sleep(time.Second * 5)
		// 3. save output to /outputs
		//      - err: non-retryable
		// save(configs.OutputDir+filename, []byte{})
		duration := time.Since(start)
		fmt.Printf("JOB: image processing for %s took %d seconds\n", filename, duration/time.Second)

		sess.MarkMessage(msg, fmt.Sprintf("%v", duration))
	}
	return nil
}
