package ready

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"

	"github.com/krok-o/krok/pkg/krok/providers/livestore"
)

// Dependencies defines the dependencies for the plugin provider.
type Dependencies struct {
	Logger    zerolog.Logger
	Connector *livestore.Connector
}

// Checker checks if the database connection is up and running.
type Checker struct {
	Dependencies
}

// NewReadyCheckProvider returns ready if the database connection is established.
func NewReadyCheckProvider(deps Dependencies) *Checker {
	return &Checker{Dependencies: deps}
}

// Ready checks if krok is ready to serve requests.
func (c *Checker) Ready(ctx context.Context) bool {
	if err := c.Connector.ExecuteWithTransaction(ctx, c.Logger, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, "select 1"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return false
	}
	return true
}
