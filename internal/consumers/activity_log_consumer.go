package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"event_service/global"
	"event_service/internal/common"
	"event_service/internal/dto"
	"event_service/internal/services"

	"github.com/rabbitmq/amqp091-go"
)

type activityLogConsumer struct {
	name        string
	logService  services.LogService
	channel     *amqp091.Channel
	queue       string
	isRunning   bool
	stopChannel chan bool
}

func NewActivityLogConsumer() Consumer {
	return &activityLogConsumer{
		name:        "ActivityLogConsumer",
		logService:  services.NewLogService(),
		stopChannel: make(chan bool),
	}
}

func (c *activityLogConsumer) GetName() string {
	return c.name
}

func (c *activityLogConsumer) Start(ctx context.Context) error {
	fmt.Printf("Starting %s...\n", c.name)

	c.channel = global.ActivityLogRabbitCh
	c.queue = global.Config.RabbitMQ.ActivityLogQueue

	if c.channel == nil {
		return fmt.Errorf("RabbitMQ channel is not initialized")
	}

	// Set QoS to process one message at a time
	err := c.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %v", err)
	}

	messages, err := c.channel.Consume(
		c.queue, // queue
		c.name,  // consumer
		false,   // auto-ack (we'll ack manually)
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return fmt.Errorf("%w: %v", common.ErrRabbitConsume, err)
	}

	c.isRunning = true
	fmt.Printf("%s started successfully\n", c.name)

	go c.processMessages(ctx, messages)

	return nil
}

func (c *activityLogConsumer) Stop() error {
	fmt.Printf("Stopping %s...\n", c.name)

	if !c.isRunning {
		return nil
	}

	c.isRunning = false
	close(c.stopChannel)

	// Cancel the consumer
	if c.channel != nil {
		err := c.channel.Cancel(c.name, false)
		if err != nil {
			fmt.Printf("Error canceling consumer %s: %v\n", c.name, err)
		}
	}

	fmt.Printf("%s stopped\n", c.name)
	return nil
}

func (c *activityLogConsumer) processMessages(ctx context.Context, messages <-chan amqp091.Delivery) {
	for {
		select {
		case <-c.stopChannel:
			fmt.Printf("%s message processing stopped\n", c.name)
			return

		case message, ok := <-messages:
			if !ok {
				fmt.Printf("%s message channel closed\n", c.name)
				return
			}

			c.handleMessage(ctx, message)
		}
	}
}

func (c *activityLogConsumer) handleMessage(ctx context.Context, message amqp091.Delivery) {
	fmt.Printf("Received message with routing key: %s\n", message.RoutingKey)

	processCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err := c.Handle(processCtx, message.Body)
	if err != nil {
		fmt.Printf("Error processing message: %v\n", err)

		if c.shouldRetry(message) {
			c.retryMessage(message, err)
		} else {
			c.rejectMessage(message, err)
		}
		return
	}

	// Acknowledge successful processing
	err = message.Ack(false)
	if err != nil {
		fmt.Printf("Error acknowledging message: %v\n", err)
	} else {
		fmt.Printf("Message processed and acknowledged successfully\n")
	}
}

func (c *activityLogConsumer) Handle(ctx context.Context, body []byte) error {
	var event dto.GenericEvent
	err := json.Unmarshal(body, &event)
	if err != nil {
		return fmt.Errorf("%w: %v", common.ErrEventDeserialization, err)
	}

	fmt.Printf("Processing event: %s with topic: %s\n", event.EventID, event.Topic)

	err = c.logService.ProcessEvent(ctx, &event)
	if err != nil {
		return fmt.Errorf("%w: %v", common.ErrEventProcessing, err)
	}

	return nil
}

func (c *activityLogConsumer) shouldRetry(message amqp091.Delivery) bool {
	retryCount := c.getRetryCount(message)
	maxRetries := global.Config.RabbitMQ.RetryAttempts

	return retryCount < maxRetries
}

func (c *activityLogConsumer) getRetryCount(message amqp091.Delivery) int {
	if message.Headers == nil {
		return 0
	}

	if retryCount, exists := message.Headers["x-retry-count"]; exists {
		if count, ok := retryCount.(int32); ok {
			return int(count)
		}
	}

	return 0
}

func (c *activityLogConsumer) retryMessage(message amqp091.Delivery, processingError error) {
	retryCount := c.getRetryCount(message) + 1
	retryDelay := global.Config.RabbitMQ.RetryDelaySeconds

	fmt.Printf("Retrying message (attempt %d): %v\n", retryCount, processingError)

	headers := make(amqp091.Table)
	if message.Headers != nil {
		headers = message.Headers
	}
	headers["x-retry-count"] = int32(retryCount)
	headers["x-original-routing-key"] = message.RoutingKey
	headers["x-retry-reason"] = processingError.Error()

	err := c.channel.Publish(
		global.Config.RabbitMQ.IAMExchange, // exchange
		message.RoutingKey,                 // routing key
		false,                              // mandatory
		false,                              // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        message.Body,
			Headers:     headers,
			Expiration:  fmt.Sprintf("%d000", retryDelay), // milliseconds
		},
	)

	if err != nil {
		fmt.Printf("Error publishing retry message: %v\n", err)
		c.rejectMessage(message, err)
		return
	}

	// Ack the original message
	err = message.Ack(false)
	if err != nil {
		fmt.Printf("Error acknowledging retry message: %v\n", err)
	}
}

func (c *activityLogConsumer) rejectMessage(message amqp091.Delivery, processingError error) {
	fmt.Printf("Rejecting message after max retries: %v\n", processingError)

	// Log the failed message for manual investigation
	c.logFailedMessage(message, processingError)

	// Reject the message (send to dead letter queue if configured)
	err := message.Reject(false) // false = don't requeue
	if err != nil {
		fmt.Printf("Error rejecting message: %v\n", err)
	}
}

func (c *activityLogConsumer) logFailedMessage(message amqp091.Delivery, processingError error) {
	fmt.Printf("FAILED MESSAGE LOG:\n")
	fmt.Printf("  Routing Key: %s\n", message.RoutingKey)
	fmt.Printf("  Error: %v\n", processingError)
	fmt.Printf("  Body: %s\n", string(message.Body))
	fmt.Printf("  Headers: %+v\n", message.Headers)
	fmt.Printf("  Retry Count: %d\n", c.getRetryCount(message))
}
