package executor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
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
	Logger        zerolog.Logger
	CommandRuns   providers.CommandRunStorer
	CommandStorer providers.CommandStorer
	Clock         providers.Clock
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
func (ime *InMemoryExecuter) CreateRun(ctx context.Context, event *models.Event, commands []*models.Command) error {
	log := ime.Logger.
		With().
		Int("event_id", event.ID).
		Int("repository_id", event.RepositoryID).
		Int("platform_id", event.VCS).
		Int("commands", len(commands)).
		Logger()

	platform, found := models.SupportedPlatforms[event.VCS]
	if !found {
		return fmt.Errorf("failed to find %d in supported platforms", event.VCS)
	}

	log = log.With().Str("platform", platform.Name).Logger()

	log.Info().Msg("Starting run")
	cmds := make([]*exec.Cmd, 0)
	// Start these here with the runner go routine
	for _, c := range commands {
		if !c.Enabled {
			log.Debug().Str("name", c.Name).Msg("Skipping as command is disabled.")
			continue
		}

		if ok, err := ime.CommandStorer.IsPlatformSupported(ctx, c.ID, event.VCS); err != nil {
			log.Debug().Err(err).Msg("Failed to get is platform is supported by command.")
			return err
		} else if !ok {
			log.Debug().Str("name", c.Name).Msg("Command does not support platform.")
			continue
		}

		settings, err := ime.CommandStorer.ListSettings(ctx, c.ID)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to get settings for command.")
			return err
		}

		// We aren't going to save these because it could be things like tokens which are
		// confidential. The platform must always be the first arg.
		args := []string{
			fmt.Sprintf("platform:%s", platform.Name),
		}
		for _, s := range settings {
			args = append(args, fmt.Sprintf("%s:%s", s.Key, s.Value))
		}

		commandRun := &models.CommandRun{
			EventID:     event.ID,
			CommandName: c.Name,
			Status:      "created",
			Outcome:     "",
			CreateAt:    ime.Clock.Now(),
		}
		commandRun, err = ime.CommandRuns.CreateRun(ctx, commandRun)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to create run for command")
			return err
		}
		location := filepath.Join(c.Location, c.Name)
		// run the plugin, which should be an executable.
		cmd := exec.Command(location, strings.Join(args, ","))
		cmds = append(cmds, cmd)
		log.Debug().Str("location", location).Msg("Preparing to run command at location...")
		// this needs its own context, since the context from above is already cancelled.
		go ime.runCommand(cmd, commandRun.ID, []byte(event.Payload))
	}
	ime.runsLock.Lock()
	ime.runs[event.ID] = cmds
	ime.runsLock.Unlock()
	return nil
}

// runCommand takes a single command and executes it, waiting for it to finish,
// or time out. Either way, it will update the corresponding command row.
func (ime *InMemoryExecuter) runCommand(cmd *exec.Cmd, commandRunID int, payload []byte) {
	done := make(chan error, 1)
	update := func(status string, outcome string) {
		ime.Logger.Debug().Int("command_run_id", commandRunID).Str("status", status).Str("outcome", outcome).Msg("Updating command run entry.")
		if err := ime.CommandRuns.UpdateRunStatus(context.Background(), commandRunID, status, outcome); err != nil {
			ime.Logger.Debug().Err(err).Msg("Updating status of command failed.")
		}
	}
	buffer := bytes.Buffer{}
	buffer.Write(payload)
	cmd.Stdin = &buffer
	stdErr := bytes.Buffer{}
	cmd.Stderr = &stdErr
	stdOut := bytes.Buffer{}
	cmd.Stdout = &stdOut

	if err := cmd.Start(); err != nil {
		update("failed", err.Error())
		ime.Logger.Debug().Err(err).Msg("Failed to start command.")
		return
	}

	go func() {
		done <- cmd.Wait()
	}()

	for {
		select {
		case err := <-done:
			if err != nil {
				update("failed", stdErr.String())
				ime.Logger.Debug().Err(err).Msg("Failed to run command.")
				return
			}
			update("success", stdOut.String())
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

// CancelRun will cancel a run and mark all commands as cancelled then remove the entry from the run map.
// If the kill was unsuccessful, the user can try running it again.
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
	}
	if killError {
		return errors.New("there was an error while cancelling running commands, " +
			"please inspect the log for more details")
	}
	ime.runsLock.Lock()
	delete(ime.runs, id)
	ime.runsLock.Unlock()
	ime.Logger.Debug().Msg("All commands successfully cancelled.")
	return nil
}
