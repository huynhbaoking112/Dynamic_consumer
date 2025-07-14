package consumers

import (
	"context"
	"fmt"
	"sync"
)

type ConsumerManager struct {
	consumers []Consumer
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewConsumerManager() *ConsumerManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &ConsumerManager{
		consumers: make([]Consumer, 0),
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (cm *ConsumerManager) RegisterConsumer(consumer Consumer) {
	cm.consumers = append(cm.consumers, consumer)
	fmt.Printf("Registered consumer: %s\n", consumer.GetName())
}

func (cm *ConsumerManager) StartAll() error {
	fmt.Printf("Starting %d consumers...\n", len(cm.consumers))

	for _, consumer := range cm.consumers {
		cm.wg.Add(1)
		go cm.startConsumer(consumer)
	}

	fmt.Println("All consumers started")
	return nil
}

func (cm *ConsumerManager) StopAll() error {
	fmt.Println("Stopping all consumers...")

	// Cancel the context to signal stop
	cm.cancel()

	// Stop each consumer
	for _, consumer := range cm.consumers {
		err := consumer.Stop()
		if err != nil {
			fmt.Printf("Error stopping consumer %s: %v\n", consumer.GetName(), err)
		}
	}

	// Wait for all goroutines to finish
	cm.wg.Wait()

	fmt.Println("All consumers stopped")
	return nil
}

func (cm *ConsumerManager) Wait() {
	cm.wg.Wait()
}

func (cm *ConsumerManager) startConsumer(consumer Consumer) {
	defer cm.wg.Done()

	err := consumer.Start(cm.ctx)
	if err != nil {
		fmt.Printf("Error starting consumer %s: %v\n", consumer.GetName(), err)
		return
	}

	// Wait for context cancellation
	<-cm.ctx.Done()
	fmt.Printf("Consumer %s context cancelled\n", consumer.GetName())
}

func (cm *ConsumerManager) GetConsumerCount() int {
	return len(cm.consumers)
}

func (cm *ConsumerManager) GetConsumerNames() []string {
	names := make([]string, len(cm.consumers))
	for i, consumer := range cm.consumers {
		names[i] = consumer.GetName()
	}
	return names
}
