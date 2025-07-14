# Consumer Architecture Design - Event Service

## 1. Overview

Event Service sử dụng **Consumer Pattern** để xử lý messages từ RabbitMQ một cách bất đồng bộ và có khả năng mở rộng cao. Kiến trúc consumer được thiết kế theo **Clean Architecture** principles với clear separation of concerns.

## 2. Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   RabbitMQ      │───▶│   Consumer      │───▶│   Service       │
│   (Message      │    │   Layer         │    │   Layer         │
│    Broker)      │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │                        │
                              ▼                        ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   Consumer      │    │   Repository    │
                       │   Manager       │    │   Layer         │
                       └─────────────────┘    └─────────────────┘
                                                       │
                                                       ▼
                                              ┌─────────────────┐
                                              │   MongoDB       │
                                              │   (Storage)     │
                                              └─────────────────┘
```

## 3. Core Components

### 3.1 Consumer Interface

```go
type Consumer interface {
    Start(ctx context.Context) error
    Stop() error
    GetName() string
}
```

**Responsibilities:**
- Định nghĩa contract cho tất cả consumers
- Lifecycle management (Start/Stop)
- Identity management (GetName)

### 3.2 Message Handler Interface

```go
type MessageHandler interface {
    Handle(ctx context.Context, body []byte) error
}
```

**Responsibilities:**
- Xử lý message content
- Deserialization và validation
- Delegate business logic to service layer

### 3.3 Consumer Manager

```go
type ConsumerManager struct {
    consumers []Consumer
    wg        sync.WaitGroup
    ctx       context.Context
    cancel    context.CancelFunc
}
```

**Responsibilities:**
- Quản lý lifecycle của tất cả consumers
- Concurrent startup và graceful shutdown
- Error handling và recovery

## 4. Consumer Implementation Pattern

### 4.1 Standard Consumer Structure

```go
type {ConsumerName}Consumer struct {
    name        string
    service     {Service}Interface
    channel     *amqp091.Channel
    queue       string
    isRunning   bool
    stopChannel chan bool
}
```

### 4.2 Constructor Pattern

```go
func New{ConsumerName}Consumer() Consumer {
    return &{consumerName}Consumer{
        name:        "{ConsumerName}Consumer",
        service:     services.New{Service}Service(),
        stopChannel: make(chan bool),
    }
}
```

### 4.3 Implementation Methods

#### Start Method
```go
func (c *{consumerName}Consumer) Start(ctx context.Context) error {
    // 1. Initialize RabbitMQ connection
    // 2. Set QoS settings
    // 3. Start consuming messages
    // 4. Launch message processing goroutine
    // 5. Return immediately (non-blocking)
}
```

#### Stop Method
```go
func (c *{consumerName}Consumer) Stop() error {
    // 1. Set running flag to false
    // 2. Close stop channel
    // 3. Cancel RabbitMQ consumer
    // 4. Wait for graceful shutdown
}
```

#### Handle Method
```go
func (c *{consumerName}Consumer) Handle(ctx context.Context, body []byte) error {
    // 1. Deserialize message
    // 2. Validate message structure
    // 3. Delegate to service layer
    // 4. Return error for retry logic
}
```

## 5. Message Processing Pipeline

### 5.1 Processing Flow

```
RabbitMQ Message → Consumer.Start() → Message Channel → 
handleMessage() → Handle() → Service.Process() → 
Repository.Save() → Database
```

### 5.2 Error Handling Flow

```
Processing Error → shouldRetry() → 
├─ True  → retryMessage() → Publish with delay → ACK original
└─ False → rejectMessage() → Log failure → REJECT to DLQ
```

### 5.3 Success Flow

```
Successful Processing → message.Ack(false) → Remove from queue
```

## 6. Retry Strategy

### 6.1 Retry Headers

```go
headers["x-retry-count"] = int32(retryCount)
headers["x-original-routing-key"] = message.RoutingKey
headers["x-retry-reason"] = processingError.Error()
```

### 6.2 Retry Logic

```go
func (c *consumer) shouldRetry(message amqp091.Delivery) bool {
    retryCount := c.getRetryCount(message)
    maxRetries := global.Config.RabbitMQ.RetryAttempts
    return retryCount < maxRetries
}
```

### 6.3 Exponential Backoff

- **First Retry**: Immediate
- **Subsequent Retries**: Configurable delay (default: 5 seconds)
- **Max Retries**: Configurable (default: 3 attempts)

## 7. Configuration Management

### 7.1 RabbitMQ Configuration

```yaml
rabbitmq:
  host: "localhost"
  port: 5672
  user: "guest"
  password: "guest"
  iam_exchange: "iam_events_topic"
  activity_log_queue: "iam_activity_log_queue"
  activity_log_binding_key: "#.log"
  retry_attempts: 3
  retry_delay_seconds: 5
```

### 7.2 Consumer-Specific Configuration

```go
type ConsumerConfig struct {
    QueueName    string
    BindingKey   string
    PrefetchCount int
    RetryAttempts int
    RetryDelay    int
}
```

## 8. Error Handling Patterns

### 8.1 Error Categories

1. **Transient Errors** (Retry-able):
   - Network timeouts
   - Database connection issues
   - Temporary service unavailability

2. **Permanent Errors** (Non-retry-able):
   - Invalid message format
   - Business rule violations
   - Authentication failures

### 8.2 Error Handling Strategy

```go
func (c *consumer) handleMessage(ctx context.Context, message amqp091.Delivery) {
    processCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    err := c.Handle(processCtx, message.Body)
    if err != nil {
        if c.shouldRetry(message) {
            c.retryMessage(message, err)
        } else {
            c.rejectMessage(message, err)
        }
        return
    }

    message.Ack(false)
}
```

## 9. Graceful Shutdown

### 9.1 Shutdown Sequence

```
Signal → ConsumerManager.StopAll() → 
Consumer.Stop() → Channel.Cancel() → 
WaitGroup.Wait() → Connection.Close()
```

### 9.2 Implementation

```go
func setupGracefulShutdown() {
    c := make(chan os.Signal, 1)
    signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-c
        fmt.Println("Received shutdown signal...")
        
        err := ConsumerManager.StopAll()
        if err != nil {
            fmt.Printf("Error stopping consumers: %v\n", err)
        }

        // Cleanup resources
        if global.ActivityLogRabbitCh != nil {
            global.ActivityLogRabbitCh.Close()
        }
        if global.RabbitMQ != nil {
            global.RabbitMQ.Close()
        }

        os.Exit(0)
    }()
}
```

## 10. Adding New Consumers

### 10.1 Step-by-Step Guide

1. **Define Consumer Struct**
   ```go
   type newFeatureConsumer struct {
       name        string
       service     services.NewFeatureService
       channel     *amqp091.Channel
       queue       string
       isRunning   bool
       stopChannel chan bool
   }
   ```

2. **Implement Consumer Interface**
   ```go
   func (c *newFeatureConsumer) Start(ctx context.Context) error { ... }
   func (c *newFeatureConsumer) Stop() error { ... }
   func (c *newFeatureConsumer) GetName() string { ... }
   ```

3. **Implement MessageHandler Interface**
   ```go
   func (c *newFeatureConsumer) Handle(ctx context.Context, body []byte) error { ... }
   ```

4. **Create Constructor**
   ```go
   func NewNewFeatureConsumer() Consumer {
       return &newFeatureConsumer{
           name:        "NewFeatureConsumer",
           service:     services.NewNewFeatureService(),
           stopChannel: make(chan bool),
       }
   }
   ```

5. **Register in ConsumerManager**
   ```go
   func InitConsumers() {
       ConsumerManager = consumers.NewConsumerManager()
       
       // Existing consumers
       activityLogConsumer := consumers.NewActivityLogConsumer()
       ConsumerManager.RegisterConsumer(activityLogConsumer)
       
       // New consumer
       newFeatureConsumer := consumers.NewNewFeatureConsumer()
       ConsumerManager.RegisterConsumer(newFeatureConsumer)
       
       ConsumerManager.StartAll()
   }
   ```

### 10.2 Configuration Updates

1. **Add Queue Configuration**
   ```yaml
   rabbitmq:
     new_feature_queue: "iam_new_feature_queue"
     new_feature_binding_key: "*.new_feature.*"
   ```

2. **Update Config Struct**
   ```go
   type RabbitMQ struct {
       // Existing fields...
       NewFeatureQueue      string `mapstructure:"new_feature_queue"`
       NewFeatureBindingKey string `mapstructure:"new_feature_binding_key"`
   }
   ```

3. **Update RabbitMQ Initialization**
   ```go
   func InitRabbitMQ() {
       // Declare new queue
       q, err := ch.QueueDeclare(
           cfg.NewFeatureQueue,
           true,  // durable
           false, // delete when unused
           false, // exclusive
           false, // no-wait
           nil,   // arguments
       )
       
       // Bind queue to exchange
       err = ch.QueueBind(
           q.Name,
           cfg.NewFeatureBindingKey,
           cfg.IAMExchange,
           false,
           nil,
       )
   }
   ```

## 11. Best Practices

### 11.1 Message Processing

1. **Always use timeouts** for message processing
2. **Implement idempotency** in business logic
3. **Use structured logging** with correlation IDs
4. **Validate message format** before processing
5. **Handle partial failures** gracefully

### 11.2 Error Handling

1. **Categorize errors** (transient vs permanent)
2. **Log failed messages** with full context
3. **Use dead letter queues** for poison messages
4. **Implement circuit breakers** for external dependencies
5. **Monitor retry patterns** for system health

### 11.3 Performance

1. **Set appropriate QoS** (prefetch count)
2. **Use connection pooling** when needed
3. **Implement backpressure** mechanisms
4. **Monitor memory usage** and garbage collection
5. **Optimize database operations** (batching, indexing)


Khi implement consumer mới, hãy follow các patterns và best practices đã được thiết lập để đảm bảo consistency và reliability của hệ thống. 