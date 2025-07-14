package initialize

import "fmt"

func Run() {
	fmt.Println("Initializing Event Service components...")

	LoadConfig()
	fmt.Println("Configuration loaded")

	InitMongoDB()
	fmt.Println("MongoDB connected")

	CreateMongoDBIndexes()
	fmt.Println("MongoDB indexes created")

	InitRabbitMQ()
	fmt.Println("RabbitMQ connected")

	InitConsumers()
	fmt.Println("Consumers initialized")

	fmt.Println("All components initialized successfully")
}
