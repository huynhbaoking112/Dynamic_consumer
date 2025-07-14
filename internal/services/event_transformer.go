package services

import (
	"fmt"

	"event_service/internal/dto"
)

type EventTransformer struct{}

func NewEventTransformer() *EventTransformer {
	return &EventTransformer{}
}

// ValidateEventStructure validates that the event has the required structure
func (t *EventTransformer) ValidateEventStructure(event *dto.GenericEvent) error {
	if event == nil {
		return fmt.Errorf("event is nil")
	}

	// Check required fields
	if event.EventID == "" {
		return fmt.Errorf("eventId is required")
	}

	if event.Topic == "" {
		return fmt.Errorf("topic is required")
	}

	if event.SourceService == "" {
		return fmt.Errorf("sourceService is required")
	}

	if event.Timestamp == "" {
		return fmt.Errorf("timestamp is required")
	}

	return nil
}
