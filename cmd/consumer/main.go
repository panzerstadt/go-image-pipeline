package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/panzerstadt/go-image-pipeline/configs"
)

func main() {
	consumer := get_consumer()
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			err := consumer.Consume(ctx, []string{configs.TestTopic}, &consumerHandler{})
			if err != nil {
				log.Printf("consume error: %v", err)
			}
			if ctx.Err() != nil {
				return
			}

		}
	}()

	<-signals
	log.Println("shutting down")
	cancel()
	wg.Wait()
}

type consumerHandler struct{}

func (h *consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("Received message: key=%s, value=%s, partition=%d, offset=%d\n", string(msg.Key), string(msg.Value), msg.Partition, msg.Offset)
		task := receive(msg)
		filename := task.Path

		// 2. resize image and make progressive jpeg
		//      - err: non-retryable
		middlePath := path.Join(configs.IntermediateDir, filename)
		_, err := os.Stat(middlePath)
		if err != nil {
			fmt.Println("file not found at: ", middlePath)
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
		fmt.Println("imagemagick process output: " + string(output))
		// 2b. call llm to guess and fill in jpeg metadata
		//      - err: retryable
		//      - stub: wait 5 seconds
		time.Sleep(time.Second * 5)
		// 3. save output to /outputs
		//      - err: non-retryable
		// save(configs.OutputDir+filename, []byte{})
		sess.MarkMessage(msg, "")
	}
	return nil
}
