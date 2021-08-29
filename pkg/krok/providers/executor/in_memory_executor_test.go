package executor

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/docker/docker/client"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/krok-o/krok/pkg/krok/providers/mocks"
	"github.com/krok-o/krok/pkg/models"
)

func TestCreateRun(t *testing.T) {
	if _, err := client.NewClientWithOpts(client.FromEnv); err != nil {
		t.Skip("Could not run test. This test requires Docker to be accessible.")
	}
	logger := zerolog.New(os.Stderr)
	mcr := &mocks.CommandRunStorer{}
	mcr.On("CreateRun", mock.Anything, &models.CommandRun{
		EventID:     1,
		CommandName: "test-command",
		Status:      "created",
		Outcome:     "",
		CreateAt:    time.Date(1981, 1, 1, 1, 1, 1, 1, time.UTC),
	}).Return(&models.CommandRun{
		ID:          1,
		EventID:     1,
		CommandName: "test-command",
		Status:      "created",
		Outcome:     "",
		CreateAt:    time.Date(1981, 1, 1, 1, 1, 1, 1, time.UTC),
	}, nil)
	mcs := &mocks.CommandStorer{}
	mcs.On("IsPlatformSupported", mock.Anything, 1, 1).Return(true, nil)
	mcs.On("ListSettings", mock.Anything, 1).Return(nil, nil)
	mt := &mocks.Clock{}
	mt.On("Now").Return(time.Date(1981, 1, 1, 1, 1, 1, 1, time.UTC))
	ime := NewInMemoryExecutor(Config{
		DefaultMaximumCommandRuntime: 10,
	}, Dependencies{
		Logger:        logger,
		CommandRuns:   mcr,
		CommandStorer: mcs,
		Clock:         mt,
	})
	err := ime.CreateRun(context.Background(), &models.Event{
		ID:           1,
		EventID:      "id",
		CreateAt:     time.Date(1981, 1, 1, 1, 1, 1, 1, time.UTC),
		RepositoryID: 1,
		Payload:      "{}",
		VCS:          1,
		EventType:    "push",
	}, []*models.Command{
		{
			Name: "test-command",
			ID:   1,
			Repositories: []*models.Repository{
				{
					Name: "test-repo",
					ID:   1,
					URL:  "https://github.com/Skarlso/test",
					VCS:  1,
					Auth: &models.Auth{Secret: "secret"},
				},
			},
			Image:   "skarlso/slack-notification:v0.0.5",
			Enabled: true,
			Platforms: []models.Platform{
				{
					ID:   1,
					Name: "github",
				},
			},
		},
	})
	assert.NoError(t, err)
	// todo: wait here for the finishing of containers and such.
	// also add a Skip if there is no docker socket available.
}
