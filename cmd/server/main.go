package main

import (
	"fmt"

	"event_service/internal/initialize"
)

func main() {
	fmt.Println("Starting Event Service...")

	// Initialize all components
	initialize.Run()

	fmt.Println("Event Service started successfully")

	// Wait for consumers to finish (they handle graceful shutdown internally)
	initialize.WaitForConsumers()

	fmt.Println("Event Service shutdown completed")
}
