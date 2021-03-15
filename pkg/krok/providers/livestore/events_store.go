package livestore

import (
	"context"

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

const (
	eventsStoreTable = "events"
)

// eventsStore is a postgres based store for eventStorer.
type eventsStore struct {
	EventsStoreDependencies
}

// EventsStoreDependencies eventsStoreStore specific dependencies.
type EventsStoreDependencies struct {
	Dependencies
	Connector *Connector
}

// NewEventsStorer creates a new eventsStore
func NewEventsStorer(deps EventsStoreDependencies) *eventsStore {
	return &eventsStore{EventsStoreDependencies: deps}
}

var _ providers.EventsStorer = &eventsStore{}

// Create an event.
func (e eventsStore) Create(ctx context.Context, event *models.Event) (*models.Event, error) {
	panic("implement me")
}

// ListEventsForRepository gets paginated list of events for a repository.
func (e eventsStore) ListEventsForRepository(ctx context.Context, repoID int, options models.ListOptions) (*[]models.Event, error) {
	panic("implement me")
}

// GetEvent retrieves details about an event.
func (e eventsStore) GetEvent(ctx context.Context, eventID int) (*models.Event, error) {
	panic("implement me")
}
