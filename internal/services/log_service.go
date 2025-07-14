package services

import (
	"context"
	"fmt"
	"time"

	"event_service/internal/common"
	"event_service/internal/dto"
	"event_service/internal/models"
	"event_service/internal/repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type logService struct {
	activityLogRepo repo.ActivityLogRepository
	transformer     *EventTransformer
}

func NewLogService() LogService {
	return &logService{
		activityLogRepo: repo.NewActivityLogRepository(),
		transformer:     NewEventTransformer(),
	}
}

func (s *logService) ProcessEvent(ctx context.Context, event *dto.GenericEvent) error {

	if err := s.transformer.ValidateEventStructure(event); err != nil {
		return fmt.Errorf("%w: %v", common.ErrEventValidation, err)
	}

	timestamp, err := s.parseTimestamp(event.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp: %v", err)
	}

	activityLog := &models.ActivityLog{
		ID:            primitive.NilObjectID,
		EventID:       event.EventID,
		Topic:         event.Topic,
		SourceService: event.SourceService,
		Timestamp:     timestamp,
		Payload:       event.Payload,
	}

	if err := s.activityLogRepo.Create(ctx, activityLog); err != nil {
		return fmt.Errorf("%w: %v", common.ErrEventProcessing, err)
	}

	fmt.Printf("Event %s processed successfully and saved to database\n", event.EventID)
	return nil
}

func (s *logService) parseTimestamp(timestampStr string) (time.Time, error) {
	// Try different timestamp formats
	formats := []string{
		time.RFC3339,     // 2006-01-02T15:04:05Z07:00
		time.RFC3339Nano, // 2006-01-02T15:04:05.999999999Z07:00
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timestampStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse timestamp: %s", timestampStr)
}
