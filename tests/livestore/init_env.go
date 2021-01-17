package livestore

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rs/zerolog"
)

// hostname can be dynamic, dependent on whether we are running on CI or locally.
var hostname = "localhost:5432"

// createTestContainerIfNotCI uses an ephemeral postgres container to run a real test.
// the cleanup has to be called by the test runner.
func createTestContainerIfNotCI() error {
	logger := zerolog.New(os.Stderr)
	if _, ok := os.LookupEnv("CIRCLECI"); ok {
		logger.Debug().Msg("On circleci, skipping ephemeral container.")
		// skip circleci environment and do nothing on cleanup.
		return nil
	}
	pool, err := dockertest.NewPool("")
	if err != nil {
		logger.Debug().Err(err).Msg("Failed to create new pool.")
		return err
	}
	cwd, err := os.Getwd()
	if err != nil {
		logger.Debug().Err(err).Msg("Failed to get working director.")
		return err
	}
	dbInit := filepath.Join(cwd, "dbinit")
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13.1-alpine",
		Env: []string{
			"POSTGRES_USER=krok",
			"POSTGRES_PASSWORD=password123",
		},
		Mounts: []string{dbInit + ":/docker-entrypoint-initdb.d"},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})

	if err != nil {
		logger.Debug().Err(err).Msg("Failed to run with options.")
		return err
	}

	if err = pool.Retry(func() error {
		var err error
		db, err := sql.Open("postgres", fmt.Sprintf("postgres://krok:password123@localhost:%s/%s?sslmode=disable", resource.GetPort("5432/tcp"), "krok"))
		if err != nil {
			logger.Debug().Err(err).Msg("Failed to open new connection.")
			return err
		}
		return db.Ping()
	}); err != nil {
		logger.Debug().Err(err).Msg("Retry eventually failed.")
		return err
	}

	hostname = "localhost:" + resource.GetPort("5432/tcp")
	logger.Debug().Str("hostname", hostname).Msg("Hostname set to ephemeral container port.")

	return nil
}

// set up ephemeral docker test container if not on ci
func init() {
	if err := createTestContainerIfNotCI(); err != nil {
		log.Fatal(err)
	}
}
