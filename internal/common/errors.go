package common

import "errors"

var (
	// MongoDB errors
	ErrMongoConnection = errors.New("failed to connect to MongoDB")
	ErrMongoInsert     = errors.New("failed to insert document to MongoDB")
	ErrMongoQuery      = errors.New("failed to query MongoDB")

	// RabbitMQ errors
	ErrRabbitConnection = errors.New("failed to connect to RabbitMQ")
	ErrRabbitChannel    = errors.New("failed to create RabbitMQ channel")
	ErrRabbitConsume    = errors.New("failed to consume from RabbitMQ")
	ErrRabbitPublish    = errors.New("failed to publish to RabbitMQ")
	ErrRabbitAck        = errors.New("failed to acknowledge message")
	ErrRabbitNack       = errors.New("failed to negative acknowledge message")

	// Event processing errors
	ErrEventDeserialization = errors.New("failed to deserialize event")
	ErrEventProcessing      = errors.New("failed to process event")
	ErrEventValidation      = errors.New("event validation failed")

	// Configuration errors
	ErrConfigLoad       = errors.New("failed to load configuration")
	ErrConfigValidation = errors.New("configuration validation failed")
)
