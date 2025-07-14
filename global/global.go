package global

import (
	"event_service/pkg/setting"

	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Config *setting.Config

	MongoDB             *mongo.Database     // MongoDB database connection
	RabbitMQ            *amqp091.Connection // RabbitMQ connection
	ActivityLogRabbitCh *amqp091.Channel    // RabbitMQ channel for activity logs
)
