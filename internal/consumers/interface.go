package consumers

import "context"

type Consumer interface {
	Start(ctx context.Context) error

	Stop() error

	GetName() string
}
