package providers

import (
	"context"

	"github.com/krok-o/krok/pkg/models"
)

// PlatformStorer handles operations for platforms.
type PlatformStorer interface {
	Create(ctx context.Context, p *models.Platform) (*models.Platform, error)
	Get(ctx context.Context, id int) (*models.Platform, error)
	GetByName(ctx context.Context, name string) (*models.Platform, error)
	Delete(ctx context.Context, id int) error
	// List will list all platforms. If enabled is defined, will use that to filter
	// platforms. If enabled is nil, we will list all platforms regardless if they
	// are enabled or not.
	List(ctx context.Context, enabled *bool) ([]*models.Platform, error)
	// Update will allow enabling and disabling a platform.
	Update(ctx context.Context, p *models.Platform) (*models.Platform, error)
}
