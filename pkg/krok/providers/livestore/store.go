package livestore

import (
	"time"

	"github.com/rs/zerolog"

	"github.com/krok-o/krok/pkg/krok/providers"
)

const timeoutForTransactions = 1 * time.Minute

// Config has the configuration options for the store
type Config struct {
	Hostname string
	Database string
	Username string
	Password string
}

// Dependencies defines the dependencies of this command store
type Dependencies struct {
	Logger    zerolog.Logger
	Converter providers.EnvironmentConverter
}
