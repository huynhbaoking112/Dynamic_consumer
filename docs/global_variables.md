# Global Variables Management - Event Service

## 1. Overview

Event Service sử dụng **Global Variables Pattern** để quản lý shared resources như database connections, configuration, và RabbitMQ channels. Pattern này giúp đảm bảo single instance của các resources quan trọng và dễ dàng access từ mọi nơi trong application.

## 2. Global Variables Structure

### 2.1 Current Global Variables

```go
// global/global.go
package global

import (
    "event_service/pkg/setting"
    
    "github.com/rabbitmq/amqp091-go"
    "go.mongodb.org/mongo-driver/mongo"
)

var (
    // Configuration
    Config *setting.Config

    // Database connections
    MongoDB              *mongo.Database      // MongoDB database connection
    RabbitMQ             *amqp091.Connection  // RabbitMQ connection
    ActivityLogRabbitCh  *amqp091.Channel     // RabbitMQ channel for activity logs
)
```

### 2.2 Variable Categories

1. **Configuration**: Application settings loaded from config files
2. **Database Connections**: MongoDB client và database instances
3. **Message Broker**: RabbitMQ connections và channels
4. **Future Extensions**: Logger, cache, metrics collectors

## 3. Initialization Pattern

### 3.1 Initialization Sequence

```
LoadConfig() → InitMongoDB() → InitRabbitMQ() → InitConsumers()
     ↓              ↓              ↓              ↓
Set Config → Set MongoDB → Set RabbitCh → Use Globals
```

### 3.2 Configuration Loading

```go
func LoadConfig() {
    // Load from YAML file or environment variables
    config := &setting.Config{...}
    
    // Set global configuration
    global.Config = config
}
```

### 3.3 Database Initialization

```go
func InitMongoDB() {
    cfg := global.Config.MongoDB
    
    // Create MongoDB connection
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        panic(err)
    }
    
    // Set global database instance
    global.MongoDB = client.Database(cfg.Database)
}
```

### 3.4 RabbitMQ Initialization

```go
func InitRabbitMQ() {
    cfg := global.Config.RabbitMQ
    
    // Create connection
    conn, err := amqp091.Dial(connStr)
    if err != nil {
        panic(err)
    }
    
    // Create channel
    ch, err := conn.Channel()
    if err != nil {
        panic(err)
    }
    
    // Set global variables
    global.RabbitMQ = conn
    global.ActivityLogRabbitCh = ch
}
```

## 4. Usage Patterns

### 4.1 Repository Layer Usage

```go
func NewActivityLogRepository() ActivityLogRepository {
    cfg := global.Config.MongoDB
    collection := global.MongoDB.Collection(cfg.ActivityLogCollection)
    return &activityLogRepository{
        collection: collection,
    }
}
```

### 4.2 Consumer Layer Usage

```go
func (c *activityLogConsumer) Start(ctx context.Context) error {
    c.channel = global.ActivityLogRabbitCh
    c.queue = global.Config.RabbitMQ.ActivityLogQueue
    
    // Use the global channel for consuming
    messages, err := c.channel.Consume(...)
}
```

### 4.3 Service Layer Usage

```go
func (s *logService) ProcessEvent(ctx context.Context, event *dto.GenericEvent) error {
    // Configuration access
    retryAttempts := global.Config.RabbitMQ.RetryAttempts
    
    // Business logic using config
    if retryCount > retryAttempts {
        return errors.New("max retries exceeded")
    }
}
```

## 5. Best Practices

### 5.1 Initialization Guidelines

1. **Initialize in Order**: Dependencies must be initialized before dependents
2. **Fail Fast**: Panic on initialization errors to prevent partial startup
3. **Single Assignment**: Set global variables only once during initialization
4. **Validation**: Validate configuration before setting globals
5. **Cleanup**: Properly close connections during shutdown

### 5.2 Usage Guidelines

1. **Read-Only Access**: Treat global variables as read-only after initialization
2. **No Mutation**: Don't modify global variables during runtime
3. **Null Checks**: Always check for nil before using global variables
4. **Error Handling**: Handle connection errors gracefully
5. **Testing**: Use dependency injection in tests instead of globals

### 5.3 Naming Conventions

1. **Descriptive Names**: Use clear, descriptive variable names
2. **Consistent Prefixes**: Group related variables with consistent prefixes
3. **Singular vs Plural**: Use singular for single instances, plural for collections
4. **Avoid Abbreviations**: Use full words for clarity

## 6. Adding New Global Variables

### 6.1 Step-by-Step Process

1. **Define in global.go**
   ```go
   var (
       // Existing variables...
       NewService *NewServiceClient
   )
   ```

2. **Add to Configuration**
   ```go
   type Config struct {
       // Existing fields...
       NewService NewServiceConfig `mapstructure:"newservice"`
   }
   ```

3. **Create Initialization Function**
   ```go
   func InitNewService() {
       cfg := global.Config.NewService
       
       client := &NewServiceClient{
           Host: cfg.Host,
           Port: cfg.Port,
       }
       
       if err := client.Connect(); err != nil {
           panic(err)
       }
       
       global.NewService = client
   }
   ```

4. **Update Run Sequence**
   ```go
   func Run() {
       LoadConfig()
       InitMongoDB()
       InitRabbitMQ()
       InitNewService()  // Add here
       InitConsumers()
   }
   ```

### 6.2 Configuration Pattern

```go
type NewServiceConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Timeout  int    `mapstructure:"timeout"`
    Enabled  bool   `mapstructure:"enabled"`
}
```

## 7. Error Handling

### 7.1 Initialization Errors

```go
func InitMongoDB() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("MongoDB initialization failed: %v\n", r)
            os.Exit(1)
        }
    }()
    
    // Initialization logic...
}
```

### 7.2 Runtime Error Handling

```go
func (r *repository) Create(ctx context.Context, doc interface{}) error {
    if global.MongoDB == nil {
        return errors.New("MongoDB not initialized")
    }
    
    collection := global.MongoDB.Collection("collection_name")
    // Database operations...
}
```

## 8. Testing Considerations

### 8.1 Test Setup

```go
func setupTestGlobals() {
    // Set test configuration
    global.Config = &setting.Config{
        MongoDB: setting.MongoDB{
            Database: "test_db",
        },
    }
    
    // Set test database
    global.MongoDB = testClient.Database("test_db")
}
```

### 8.2 Dependency Injection for Tests

```go
type Repository struct {
    db *mongo.Database
}

func NewRepository(db *mongo.Database) *Repository {
    if db == nil {
        db = global.MongoDB  // Use global as fallback
    }
    return &Repository{db: db}
}

// In tests
func TestRepository(t *testing.T) {
    testDB := setupTestDB()
    repo := NewRepository(testDB)  // Inject test DB
    // Test logic...
}
```

## 9. Cleanup and Shutdown

### 9.1 Graceful Shutdown

```go
func setupGracefulShutdown() {
    c := make(chan os.Signal, 1)
    signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-c
        fmt.Println("Shutting down...")
        
        // Close RabbitMQ connections
        if global.ActivityLogRabbitCh != nil {
            global.ActivityLogRabbitCh.Close()
        }
        if global.RabbitMQ != nil {
            global.RabbitMQ.Close()
        }
        
        // Close MongoDB connections
        if global.MongoDB != nil {
            global.MongoDB.Client().Disconnect(context.Background())
        }
        
        os.Exit(0)
    }()
}
```

### 9.2 Resource Cleanup Order

```
1. Stop consumers
2. Close RabbitMQ channels
3. Close RabbitMQ connections
4. Close MongoDB connections
5. Exit application
```

## 10. Common Pitfalls

### 10.1 Issues to Avoid

1. **Race Conditions**: Don't access globals before initialization
2. **Nil Pointer Dereference**: Always check for nil
3. **Partial Initialization**: Ensure all dependencies are ready
4. **Memory Leaks**: Properly close connections during shutdown
5. **Test Pollution**: Reset globals between tests

### 10.2 Anti-Patterns

```go
// DON'T: Modify globals during runtime
func SomeFunction() {
    global.Config.MongoDB.Database = "new_db"  // BAD
}

// DON'T: Use globals without nil checks
func AnotherFunction() {
    collection := global.MongoDB.Collection("test")  // BAD if MongoDB is nil
}

// DON'T: Initialize globals in multiple places
func InitSomething() {
    global.Config = &setting.Config{}  // BAD if already initialized
}
```
