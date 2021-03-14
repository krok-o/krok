package providers

import (
	"context"

	"github.com/krok-o/krok/pkg/models"
)

// Executor manages runs regarding events for repositories.
type Executor interface {
	// CreateRun creates a run for an event.
	// The created run will have to be saved somehow. This is up to the implementation.
	// It MUST use the Event's ID as identification because that's what defines/holds
	// the currently running commands. The loose coupling between a run and the commands
	// is the event. So to cancel a Run, the user will provide the Event's ID.
	CreateRun(ctx context.Context, event *models.Event) error
	// CancelRun will cancel a run and mark all commands as cancelled.
	// The ID here is the ID of the event corresponding to this run.
	CancelRun(ctx context.Context, id int) error
}
