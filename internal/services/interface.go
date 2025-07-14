package services

import (
	"context"
	"event_service/internal/dto"
)

type LogService interface {
	ProcessEvent(ctx context.Context, event *dto.GenericEvent) error
}
