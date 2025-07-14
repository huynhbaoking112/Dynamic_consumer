package initialize

import (
	"context"
	"fmt"
	"time"

	"event_service/global"
	"event_service/internal/common"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDB() {
	fmt.Println("Initializing MongoDB connection...")

	cfg := global.Config.MongoDB

	var uri string
	if cfg.ConnectionString != "" {
		uri = cfg.ConnectionString
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	}

	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(fmt.Errorf("%w: %v", common.ErrMongoConnection, err))
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		panic(fmt.Errorf("%w: ping failed: %v", common.ErrMongoConnection, err))
	}

	global.MongoDB = client.Database(cfg.Database)

	fmt.Printf("MongoDB connected successfully to database: %s\n", cfg.Database)
}
