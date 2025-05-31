package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/panzerstadt/go-image-pipeline/shared"
)

var user = "panzerstadt"
var remote_host = "homelab-photostore"
var remote_scp_path = "/home/panzerstadt/public"
var public_images_path = "/Users/tliqun/Public/images"

type SyncProcessor struct{}

func (h *SyncProcessor) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *SyncProcessor) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// this should be on a separate topic, the sync topic
// because this deletes and recopies
func (h *SyncProcessor) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		time.Sleep(shared.RandDuration(1, 3))
		fmt.Printf("NOTIFICATIONS: Received message: key=%s, value=%s, partition=%d, offset=%d\n",
			string(msg.Key), strings.ReplaceAll(strings.TrimSpace(string(msg.Value)), "\n", " "), msg.Partition, msg.Offset)

		task := shared.ReceiveSyncTask(msg)
		path := task.Dir
		if info, err := os.Stat(task.Dir); err != nil || !info.IsDir() {
			fmt.Printf("Directory %s is invalid or does not exist. using default path\n", task.Dir)
			path = public_images_path
		}

		paths, err := os.ReadDir(path)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		fmt.Printf("%v directories will be synced to image server.\n", len(paths))
		for _, fileOrDir := range paths {
			fmt.Println(fileOrDir)
		}
		// note that this requires ssh keys already installed on the remote host
		// cmd: ssh-copy-id username@remote_host
		var scp_cmd = fmt.Sprintf("rsync -avz --delete %s %s@%s:%s", path, user, remote_host, remote_scp_path)

		fmt.Println("running this command: ", scp_cmd)
		cmd := exec.Command("bash", "-c", scp_cmd)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(output))

		// fmt.Println("NOTIFICATIONS: sleeping for 3 seconds...")
		time.Sleep(time.Second * 3)
		sess.MarkMessage(msg, "")
	}
	return nil
}
