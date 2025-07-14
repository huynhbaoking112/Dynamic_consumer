package dto

// GenericEvent represents a generic event structure from RabbitMQ
type GenericEvent struct {
	EventID       string                 `json:"eventId"`
	Topic         string                 `json:"topic"`
	SourceService string                 `json:"sourceService"`
	Timestamp     string                 `json:"timestamp"`
	Payload       map[string]interface{} `json:"payload"`
}
