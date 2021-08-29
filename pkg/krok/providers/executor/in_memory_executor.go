package executor

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/rs/zerolog"

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

// Config defines configuration for this provider
type Config struct {
	DefaultMaximumCommandRuntime int
}

// Dependencies defines dependencies for this provider
type Dependencies struct {
	Logger        zerolog.Logger
	CommandRuns   providers.CommandRunStorer
	CommandStorer providers.CommandStorer
	Clock         providers.Clock
}

// InMemoryExecutor defines an Executor which runs commands
// alongside Krok. It saves runs in a map and constantly updates it.
// Cancelling will go over all processes belonging to that run
// and kill them.
type InMemoryExecutor struct {
	Config
	Dependencies

	// For each event, a list of commandName=>containerIDs.
	// ContainerIDs are filled in as the containers are pulled and started.
	runs     map[int]*sync.Map
	runsLock sync.RWMutex
}

// NewInMemoryExecutor creates a new InMemoryExecutor which will hold all runs in its memory.
// In case of a crash, human intervention will be required.
// TODO: Later, save runs in db with the process id to cancel so Krok can pick up runs again.
func NewInMemoryExecutor(cfg Config, deps Dependencies) *InMemoryExecutor {
	m := make(map[int]*sync.Map)
	return &InMemoryExecutor{
		Config:       cfg,
		Dependencies: deps,
		runs:         m,
	}
}

// CreateRun creates a run for an event.
func (ime *InMemoryExecutor) CreateRun(ctx context.Context, event *models.Event, commands []*models.Command) error {
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
	containers := &sync.Map{}
	payload := base64.StdEncoding.EncodeToString([]byte(event.Payload))
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
			fmt.Sprintf("--platform=%s", platform.Name),
			fmt.Sprintf("--event-type=%s", event.EventType),
			fmt.Sprintf("--payload=%s", payload),
		}
		for _, s := range settings {
			args = append(args, fmt.Sprintf("--%s=%s", s.Key, s.Value))
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
		log.Debug().Str("image", c.Image).Msg("Preparing to run command...")
		containers.Store(c.Name, "")
		go ime.pullAndCreateContainer(c.Name, c.Image, args, event.ID, commandRun.ID)
	}
	ime.runsLock.Lock()
	ime.runs[event.ID] = containers
	ime.runsLock.Unlock()
	return nil
}

func (ime *InMemoryExecutor) pullAndCreateContainer(commandName, image string, args []string, eventID int, commandRunID int) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		ime.updateStatus("failed", err.Error(), commandRunID)
		ime.Logger.Debug().Err(err).Msg("Failed to create docker client.")
		return
	}
	output, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		ime.updateStatus("failed", fmt.Sprintf("failed to pull image: %s", err), commandRunID)
		ime.Logger.Debug().Err(err).Msg("Failed to pull image.")
		return
	}
	if _, err := io.Copy(os.Stdout, output); err != nil {
		ime.updateStatus("failed", fmt.Sprintf("failed to pull image: %s", err), commandRunID)
		ime.Logger.Debug().Err(err).Msg("Failed to pull image.")
		return
	}

	ime.Logger.Info().Msg("Creating container...")
	cont, err := cli.ContainerCreate(context.Background(), &container.Config{
		AttachStdout: true,
		AttachStderr: true,
		Image:        image,
		Cmd:          args,
	}, nil, nil, nil, "")
	if err != nil {
		ime.updateStatus("failed", err.Error(), commandRunID)
		ime.Logger.Debug().Err(err).Strs("warnings", cont.Warnings).Msg("Failed to create container.")
		return
	}
	ime.runsLock.Lock()
	ime.runs[eventID].Store(commandName, cont.ID)
	ime.runsLock.Unlock()
	go ime.startAndWaitForContainer(commandName, cont.ID, eventID, commandRunID)
}

// TODO: this should return an error and we should log that.
func (ime *InMemoryExecutor) updateStatus(status, outcome string, commandRunID int) {
	outcome = strconv.Quote(outcome)
	ime.Logger.Debug().Int("command_run_id", commandRunID).Str("status", status).Str("outcome", outcome).Msg("Updating command run entry.")
	if err := ime.CommandRuns.UpdateRunStatus(context.Background(), commandRunID, status, outcome); err != nil {
		ime.Logger.Debug().Err(err).Msg("Updating status of command failed.")
	}
}

// runCommand takes a single command and executes it, waiting for it to finish,
// or time out. Either way, it will update the corresponding command row.
func (ime *InMemoryExecutor) startAndWaitForContainer(commandName, containerID string, eventID, commandRunID int) {
	done := make(chan error, 1)
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		ime.updateStatus("failed", err.Error(), commandRunID)
		ime.Logger.Debug().Err(err).Msg("Failed to create docker client.")
		return
	}
	defer func() {
		// we remove the container in a `defer` instead of autoRemove, to be able to read out the logs.
		// If we use AutoRemove, the container is gone by the time we want to read the output.
		// Could try streaming the logs. But this is enough for now.
		if err := cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{
			Force: true,
		}); err != nil {
			ime.Logger.Debug().Err(err).Str("container_id", containerID).Msg("Failed to remove container.")
		}

		// we also delete this command run from memory since it has been saved in the db.
		ime.runsLock.Lock()
		ime.runs[eventID].Delete(commandName)
		// if there are no more runs for this event, remove the event entry too.
		empty := true
		ime.runs[eventID].Range(func(key, value interface{}) bool {
			empty = false
			return false
		})
		if empty {
			delete(ime.runs, eventID)
		}
		ime.runsLock.Unlock()
	}()

	ime.Logger.Info().Msg("Starting container...")
	if err := cli.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{}); err != nil {
		ime.updateStatus("failed", err.Error(), commandRunID)
		return
	}

	go func() {
		exit, err := cli.ContainerWait(context.Background(), containerID, container.WaitConditionNotRunning)
		select {
		case e := <-err:
			done <- e
		case e := <-exit:
			if e.StatusCode != 0 {
				if e.Error != nil {
					done <- errors.New(e.Error.Message)
				} else {
					done <- fmt.Errorf("status code: %d", e.StatusCode)
				}
			} else {
				done <- nil
			}
		}
	}()

	for {
		select {
		case err := <-done:
			log, logErr := cli.ContainerLogs(context.Background(), containerID, types.ContainerLogsOptions{
				ShowStderr: true,
				ShowStdout: true,
			})
			if logErr != nil {
				ime.updateStatus("failed", logErr.Error(), commandRunID)
				return
			}
			buffer := &bytes.Buffer{}
			logs := "no logs available"
			if _, err := stdcopy.StdCopy(buffer, buffer, log); err != nil {
				ime.Logger.Debug().Err(err).Msg("Failed to de-multiplex the docker log.")
			} else {
				logs = buffer.String()
			}

			if err != nil {
				ime.updateStatus("failed", logs, commandRunID)
				ime.Logger.Debug().Err(err).Msg("Failed to run command.")
				return
			}
			ime.updateStatus("success", logs, commandRunID)
			ime.Logger.Info().Msg("Successfully finished command.")
			return
		case <-time.After(time.Duration(ime.Config.DefaultMaximumCommandRuntime) * time.Second):
			// update entry
			ime.updateStatus("failed", "timeout", commandRunID)
			ime.Logger.Error().Msg("Command timed out.")
			if err := cli.ContainerKill(context.Background(), containerID, "SIGKILL"); err != nil {
				ime.Logger.Error().Str("container_id", containerID).Msg("Failed to kill process with pid.")
			}
			return
		}
	}
}

// CancelRun will cancel a run and mark all commands as cancelled then remove the entry from the run map.
// If the kill was unsuccessful, the user can try running it again.
func (ime *InMemoryExecutor) CancelRun(ctx context.Context, id int) error {
	ime.runsLock.RLock()
	commands, ok := ime.runs[id]
	ime.runsLock.RUnlock()
	if !ok {
		ime.Logger.Error().Int("id", id).Msg("Run with ID not found")
		return errors.New("run with ID not found")
	}
	killError := false
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		ime.Logger.Debug().Err(err).Msg("Failed to create docker client.")
		return err
	}

	commands.Range(func(key, value interface{}) bool {
		if value.(string) == "" {
			ime.Logger.Debug().Str("command_name", key.(string)).Msg("Command has no container running.")
			return true
		}
		if err := cli.ContainerKill(context.Background(), value.(string), "SIGKILL"); err != nil {
			ime.Logger.Err(err).Msg("Failed to kill container.")
			killError = true
		}
		return true
	})
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
