package executor

import (
	"context"
	"errors"
	"os/exec"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

// Config defines configuration for this provider
type Config struct {
	NodePath                     string
	DefaultMaximumCommandRuntime int
}

// Dependencies defines dependencies for this provider
type Dependencies struct {
	Logger      zerolog.Logger
	CommandRuns providers.CommandRunStorer
}

// InMemoryExecuter defines an Executor which runs commands
// alongside Krok. It saves runs in a map and constantly updates it.
// Cancelling will go over all processes belonging to that run
// and kill them.
type InMemoryExecuter struct {
	Config
	Dependencies

	// List of os processes which represent one command each.
	runs     map[int][]*exec.Cmd
	runsLock sync.RWMutex
}

// NewInMemoryExecuter creates a new InMemoryExecuter which will hold all runs in its memory.
// In case of a crash, human intervention will be required.
// TODO: Later, save runs in db with the process id to cancel so Krok can pick up runs again.
func NewInMemoryExecuter(cfg Config, deps Dependencies) *InMemoryExecuter {
	m := make(map[int][]*exec.Cmd)
	return &InMemoryExecuter{
		Config:       cfg,
		Dependencies: deps,
		runs:         m,
	}
}

// CreateRun creates a run for an event.
func (ime *InMemoryExecuter) CreateRun(ctx context.Context, event *models.Event) error {
	log := ime.Logger.
		With().
		Int("event_id", event.ID).
		Int("repository_id", event.RepositoryID).
		Int("commands", len(event.Commands)).
		Logger()

	log.Info().Msg("Starting run")
	commands := make([]*exec.Cmd, 0)
	// Start these here with the runner go routine
	for _, c := range event.Commands {
		// TODO: find a way to define the command parameters.
		var err error
		commandRun := &models.CommandRun{
			EventID:  event.ID,
			Status:   "created",
			Outcome:  "",
			CreateAt: time.Now(),
		}
		commandRun, err = ime.CommandRuns.CreateRun(ctx, commandRun)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to create run for command")
			return err
		}
		cmd := exec.Command(ime.NodePath, c.Location)
		commands = append(commands, cmd)
		go ime.runCommand(ctx, cmd, commandRun.ID, []byte(event.Payload))
	}
	ime.runsLock.Lock()
	ime.runs[event.ID] = commands
	ime.runsLock.Unlock()
	return nil
}

// runCommand takes a single command and executes it, waiting for it to finish,
// or time out. Either way, it will update the corresponding command row.
func (ime *InMemoryExecuter) runCommand(ctx context.Context, cmd *exec.Cmd, commandID int, payload []byte) {
	done := make(chan error, 1)
	update := func(status string, outcome string) {
		if err := ime.CommandRuns.UpdateRunStatus(ctx, commandID, status, outcome); err != nil {
			ime.Logger.Debug().Err(err).Msg("Updating status of command failed.")
		}
	}
	if err := cmd.Start(); err != nil {
		update("failed", err.Error())
		ime.Logger.Debug().Err(err).Msg("Failed to start command.")
		return
	}
	in, err := cmd.StdinPipe()
	if err != nil {
		update("failed", err.Error())
		ime.Logger.Debug().Err(err).Msg("Failed to run command.")
		return
	}
	if _, err := in.Write(payload); err != nil {
		update("failed", err.Error())
		ime.Logger.Debug().Err(err).Msg("Failed to send payload to command.")
		return
	}

	go func() {
		done <- cmd.Wait()
	}()

	for {
		select {
		case err := <-done:
			if err != nil {
				update("failed", err.Error())
				ime.Logger.Debug().Err(err).Msg("Failed to run command.")
				return
			}
			// update entry with success
			output, err := cmd.Output()
			if err != nil {
				ime.Logger.Debug().Err(err).Msg("Failed to get command output.")
			}
			update("success", string(output))
			ime.Logger.Info().Msg("Successfully finished command.")
			return
		case <-time.After(time.Duration(ime.Config.DefaultMaximumCommandRuntime) * time.Second):
			// update entry
			update("failed", "timeout")
			ime.Logger.Error().Msg("Command timed out.")
			if err := cmd.Process.Kill(); err != nil {
				ime.Logger.Error().Int("pid", cmd.Process.Pid).Msg("Failed to kill process with pid.")
			}
			return
		}
	}
}

// CancelRun will cancel a run and mark all commands as cancelled.
func (ime *InMemoryExecuter) CancelRun(ctx context.Context, id int) error {
	ime.runsLock.RLock()
	commands, ok := ime.runs[id]
	ime.runsLock.RUnlock()
	if !ok {
		ime.Logger.Error().Int("id", id).Msg("Run with ID not found")
		return errors.New("run with ID not found")
	}
	killError := false
	for _, c := range commands {
		if err := c.Process.Kill(); err != nil {
			ime.Logger.Err(err).Int("pid", c.Process.Pid).Msg("Failed to kill process with pid.")
			killError = true
		}
		// update entry
	}
	if killError {
		return errors.New("there was an error while cancelling running commands, " +
			"please inspect the log for more details")
	}
	ime.Logger.Debug().Msg("All commands successfully cancelled.")
	return nil
}
