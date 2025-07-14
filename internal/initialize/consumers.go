package initialize

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"event_service/global"
	"event_service/internal/consumers"
)

var ConsumerManager *consumers.ConsumerManager

func InitConsumers() {
	fmt.Println("Initializing consumers...")

	// Create consumer manager
	ConsumerManager = consumers.NewConsumerManager()

	// Register consumers
	activityLogConsumer := consumers.NewActivityLogConsumer()
	ConsumerManager.RegisterConsumer(activityLogConsumer)

	// Start all consumers
	err := ConsumerManager.StartAll()
	if err != nil {
		panic(fmt.Errorf("failed to start consumers: %v", err))
	}

	// Setup graceful shutdown
	setupGracefulShutdown()

	fmt.Printf("Consumers initialized successfully: %v\n", ConsumerManager.GetConsumerNames())
}

func setupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("Received shutdown signal, stopping consumers...")

		err := ConsumerManager.StopAll()
		if err != nil {
			fmt.Printf("Error stopping consumers: %v\n", err)
		}

		// Close RabbitMQ connections
		if global.ActivityLogRabbitCh != nil {
			global.ActivityLogRabbitCh.Close()
		}
		if global.RabbitMQ != nil {
			global.RabbitMQ.Close()
		}

		fmt.Println("Graceful shutdown completed")
		os.Exit(0)
	}()
}

func WaitForConsumers() {
	if ConsumerManager != nil {
		fmt.Println("Waiting for consumers...")
		ConsumerManager.Wait()
	}
}
