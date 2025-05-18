package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/IBM/sarama"
	"github.com/panzerstadt/go-image-pipeline/cmd/consumer/handlers"
	"github.com/panzerstadt/go-image-pipeline/configs"
)

func main() {
	imageProcessorConsumer := make_consumer_for_consumer_group(configs.ConsumerGroupImageProcessor)
	defer imageProcessorConsumer.Close()

	// shard for now. TODO: start my own consumer separately
	notificationsConsumer := make_consumer_for_consumer_group(configs.ConsumerGroupNotifications)
	defer notificationsConsumer.Close()

	analyticsConsumer := make_consumer_for_consumer_group(configs.ConsumerGroupAnalytics)
	defer analyticsConsumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var wg sync.WaitGroup
	startConsumer(&wg, ctx, notificationsConsumer, configs.TopicImageJobs, &handlers.NotificationsProcessor{})
	startConsumer(&wg, ctx, imageProcessorConsumer, configs.TopicImageJobs, &handlers.ImageProcessor{})
	startConsumer(&wg, ctx, analyticsConsumer, configs.TopicImageJobs, &handlers.AnalyticsProcessor{})

	<-signals
	log.Println("shutting down")
	cancel()
	wg.Wait()
}

func startConsumer(wg *sync.WaitGroup, ctx context.Context, consumer sarama.ConsumerGroup, topic string, handler sarama.ConsumerGroupHandler) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			err := consumer.Consume(ctx, []string{topic}, handler)
			if err != nil {
				log.Printf("consumer error: %v", err)
			}
			if ctx.Err() != nil {
				return
			}

		}
	}()
}
