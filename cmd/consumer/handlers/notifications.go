package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/panzerstadt/go-image-pipeline/configs"
	"github.com/panzerstadt/go-image-pipeline/shared"
)

type NotificationsProcessor struct{}

func (h *NotificationsProcessor) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *NotificationsProcessor) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *NotificationsProcessor) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		time.Sleep(shared.RandDuration(1, 3))
		fmt.Printf("NOTIFICATIONS: Received message: key=%s, value=%s, partition=%d, offset=%d\n", string(msg.Key), strings.ReplaceAll(strings.TrimSpace(string(msg.Value)), "\n", " "), msg.Partition, msg.Offset)
		task := shared.ReceiveResizeTask(msg)
		filename := task.Path

		notification := fmt.Sprintf("We have started processing the filename  %v\n", filename)
		shared.FileAppend(configs.NotificationsFile, []byte(notification))

		// fmt.Println("NOTIFICATIONS: sleeping for 3 seconds...")
		time.Sleep(time.Second * 3)
		sess.MarkMessage(msg, "")
	}
	return nil
}
