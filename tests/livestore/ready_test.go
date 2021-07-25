package livestore

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/krok-o/krok/pkg/krok/providers/environment"
	"github.com/krok-o/krok/pkg/krok/providers/livestore"
	"github.com/krok-o/krok/pkg/krok/providers/ready"
	"github.com/krok-o/krok/tests/dbaccess"
)

func TestReadynessTest(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	env := environment.NewDockerConverter(environment.Dependencies{Logger: logger})
	connector := livestore.NewDatabaseConnector(livestore.Config{
		Hostname: hostname,
		Database: dbaccess.Db,
		Username: dbaccess.Username,
		Password: dbaccess.Password,
	}, livestore.Dependencies{
		Logger:    logger,
		Converter: env,
	})
	ready := ready.NewReadyCheckProvider(ready.Dependencies{
		Logger:    logger,
		Connector: connector,
	})
	assert.True(t, ready.Ready(context.Background()))
}
