package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/panzerstadt/go-image-pipeline/shared"
)

type AnalyticsProcessor struct{}

func (h *AnalyticsProcessor) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *AnalyticsProcessor) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *AnalyticsProcessor) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		time.Sleep(shared.RandDuration(1, 3))
		fmt.Printf("ANALYTICS: Received message: key=%s, value=%s, partition=%d, offset=%d\n", string(msg.Key), strings.ReplaceAll(strings.TrimSpace(string(msg.Value)), "\n", " "), msg.Partition, msg.Offset)
		task := shared.ReceiveResizeTask(msg)
		filename := task.Path

		fmt.Printf("ANALYTICS: analytics for %s\n", filename)
		// notification := fmt.Sprintf("We have started processing the filename  %v", filename)
		// shared.FileAppend(configs.NotificationsFile, []byte(notification))

		// fmt.Println("sleeping for 3 seconds...")
		time.Sleep(time.Second * 3)
		sess.MarkMessage(msg, "")
	}
	return nil
}
