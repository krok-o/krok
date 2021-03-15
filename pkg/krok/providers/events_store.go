package providers

import (
	"context"

	"github.com/krok-o/krok/pkg/models"
)

// EventsStorer will store events.
type EventsStorer interface {
	Create(ctx context.Context, event *models.Event) (*models.Event, error)
	ListEventsForRepository(ctx context.Context, repoID int, options models.ListOptions) (*[]models.Event, error)
	GetEvent(ctx context.Context, eventID int) (*models.Event, error)
	Update(ctx context.Context, event *models.Event) (*models.Event, error)
}
