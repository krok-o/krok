package providers

import (
	"context"

	"github.com/krok-o/krok/pkg/models"
)

// CommandRunStorer will store and update individual command run
// details and progress.
type CommandRunStorer interface {
	CreateRun(ctx context.Context, run *models.CommandRun) (*models.CommandRun, error)
	UpdateRunStatus(ctx context.Context, id int, status string, outcome string) error
	Get(ctx context.Context, id int) (*models.CommandRun, error)
}
