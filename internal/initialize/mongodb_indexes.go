package initialize

import (
	"context"
	"fmt"
	"time"

	"event_service/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateMongoDBIndexes() {
	fmt.Println("Creating MongoDB indexes...")

	cfg := global.Config.MongoDB
	collection := global.MongoDB.Collection(cfg.ActivityLogCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{"topic", 1},
				{"processedAt", -1},
			},
			Options: options.Index().SetName("topic_processedAt_idx"),
		},
		{
			Keys: bson.D{
				{"payload.userId", 1},
				{"processedAt", -1},
			},
			Options: options.Index().SetName("userId_processedAt_idx"),
		},
		{
			Keys: bson.D{
				{"eventId", 1},
			},
			Options: options.Index().SetName("eventId_idx").SetUnique(true),
		},
		{
			Keys: bson.D{
				{"sourceService", 1},
				{"processedAt", -1},
			},
			Options: options.Index().SetName("sourceService_processedAt_idx"),
		},
		{
			Keys: bson.D{
				{"timestamp", -1},
			},
			Options: options.Index().SetName("timestamp_idx"),
		},
		{
			Keys: bson.D{
				{"processedAt", -1},
			},
			Options: options.Index().SetName("processedAt_idx"),
		},
	}

	names, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		fmt.Printf("Warning: Failed to create some indexes: %v\n", err)
		return
	}

	fmt.Printf("MongoDB indexes created successfully: %v\n", names)
}
