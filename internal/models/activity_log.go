package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActivityLog struct {
	ID            primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	EventID       string                 `bson:"eventId" json:"eventId"`
	Topic         string                 `bson:"topic" json:"topic"`
	SourceService string                 `bson:"sourceService" json:"sourceService"`
	Timestamp     time.Time              `bson:"timestamp" json:"timestamp"`
	Payload       map[string]interface{} `bson:"payload" json:"payload"`
	ProcessedAt   time.Time              `bson:"processedAt" json:"processedAt"`
	Version       int                    `bson:"version" json:"version"`
}
