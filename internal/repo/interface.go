package repo

import (
	"context"
	"event_service/internal/models"
)

type ActivityLogRepository interface {
	Create(ctx context.Context, log *models.ActivityLog) error
}
