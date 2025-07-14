# Event Service

Event Service is a microservice responsible for consuming events from RabbitMQ and storing activity logs in MongoDB. It follows Clean Architecture principles with clear separation of concerns.

## Architecture

The service is structured following Clean Architecture:

```
event_service/
├── cmd/
│   └── server/          # Application entry point
├── configs/             # Configuration files
├── internal/
│   ├── common/         # Constants and common errors
│   ├── consumers/      # Message consumers (equivalent to controllers)
│   ├── dto/           # Data Transfer Objects
│   ├── initialize/    # Initialization logic
│   ├── models/        # Domain models
│   ├── repo/          # Repository layer (data access)
│   └── services/      # Business logic layer
├── pkg/
│   └── setting/       # Configuration management
├── global/            # Global variables and connections
└── README.md
```

## Components

### 1. Consumer Layer (`internal/consumers/`)
- Consumes messages from RabbitMQ queues
- Handles message acknowledgment and error scenarios
- Delegates business logic to service layer

### 2. Service Layer (`internal/services/`)
- Contains business logic for processing events
- Validates and transforms event data
- Calls repository layer for data persistence

### 3. Repository Layer (`internal/repo/`)
- Handles MongoDB operations
- Implements data access interfaces
- Manages database connections and transactions

### 4. Models (`internal/models/`)
- Defines MongoDB document structures
- Contains domain entities

## Key Features

- **Event-Driven Architecture**: Consumes events from RabbitMQ Topic Exchange
- **Clean Architecture**: Modular design with clear separation of concerns
- **MongoDB Integration**: Stores activity logs in MongoDB
- **Error Handling**: Comprehensive error handling with retry logic
- **Graceful Shutdown**: Proper cleanup of connections and resources
- **Interface-Based Design**: Easy to test and extend

## Configuration

Configuration is managed through YAML files and environment variables:

```yaml
server:
  host: "localhost"
  port: 8081

mongodb:
  host: "localhost"
  port: 27017
  database: "notification"
  collection: "activity_logs"
  user: "admin"
  password: "password"

rabbitmq:
  host: "localhost"
  port: 5672
  user: "guest"
  password: "guest"
  iam_exchange: "iam_events_topic"
  activity_log_queue: "iam_activity_log_queue"
  binding_key: "#.log"
  retry_attempts: 3
  retry_delay_seconds: 5
```

## Event Processing Flow

1. **Message Consumption**: Consumer receives message from RabbitMQ queue
2. **Deserialization**: Message is deserialized into Event DTO
3. **Business Logic**: Service layer processes and validates the event
4. **Data Persistence**: Repository layer saves activity log to MongoDB
5. **Acknowledgment**: Message is acknowledged if processing succeeds

## Running the Service

```bash
# Development
go run cmd/server/main.go

# Build
go build -o event_service cmd/server/main.go

# Run binary
./event_service
```

## Dependencies

- **MongoDB Driver**: `go.mongodb.org/mongo-driver`
- **RabbitMQ AMQP**: `github.com/rabbitmq/amqp091-go`
- **Viper**: `github.com/spf13/viper` (for configuration)

## Development Status

This is the initial project setup (Step 7.10 of Phase 4). The following components are implemented:

- ✅ Project structure according to Clean Architecture
- ✅ Basic configuration management
- ✅ Interface definitions
- ✅ Model definitions
- ✅ Global variable setup
- ✅ Entry point and initialization framework

## Next Steps

- Implement MongoDB connection and operations
- Implement RabbitMQ connection and consumption
- Implement business logic in service layer
- Implement actual consumers
- Add comprehensive testing
- Add monitoring and logging 