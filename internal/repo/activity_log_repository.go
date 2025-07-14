package repo

import (
	"context"
	"fmt"
	"time"

	"event_service/global"
	"event_service/internal/common"
	"event_service/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type activityLogRepository struct {
	collection *mongo.Collection
}

func NewActivityLogRepository() ActivityLogRepository {
	cfg := global.Config.MongoDB
	collection := global.MongoDB.Collection(cfg.ActivityLogCollection)
	return &activityLogRepository{
		collection: collection,
	}
}

func (r *activityLogRepository) Create(ctx context.Context, log *models.ActivityLog) error {
	log.ProcessedAt = time.Now()
	log.Version = 1

	result, err := r.collection.InsertOne(ctx, log)
	if err != nil {
		return fmt.Errorf("%w: %v", common.ErrMongoInsert, err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		log.ID = oid
	}

	return nil
}
