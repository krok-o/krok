package livestore

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rs/zerolog"
)

// hostname can be dynamic, dependent on whether we are running on CI or locally.
var hostname = "localhost:5432"

// TestMain runs the tests for the package and allows us to bring up any external dependencies required.
func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	port, cleanup, err := createTestContainerIfNotCI()
	if err != nil {
		log.Fatal("error running test container: ", err)
	}

	if port != "" {
		hostname = "localhost:" + port
	}

	defer func() {
		if err := cleanup(); err != nil {
			log.Fatal(err)
		}
	}()

	return m.Run()
}

// createTestContainerIfNotCI uses an ephemeral postgres container to run a real test.
// the cleanup has to be called by the test runner.
func createTestContainerIfNotCI() (string, func() error, error) {
	logger := zerolog.New(os.Stderr)
	if _, ok := os.LookupEnv("CIRCLECI"); ok {
		logger.Debug().Msg("On circleci, skipping ephemeral container.")
		// skip circleci environment and do nothing on cleanup.
		// no-op teardown.
		return "", func() error { return nil }, nil
	}
	pool, err := dockertest.NewPool("")
	if err != nil {
		logger.Debug().Err(err).Msg("Failed to create new pool.")
		return "", func() error { return nil }, err
	}
	cwd, err := os.Getwd()
	if err != nil {
		logger.Debug().Err(err).Msg("Failed to get working director.")
		return "", func() error { return nil }, err
	}
	dbInit := filepath.Join(cwd, "..", "..", "dbinit")

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
		return "", func() error { return nil }, err
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
		return "", func() error { return nil }, err
	}

	hostname = "localhost:" + resource.GetPort("5432/tcp")
	logger.Debug().Str("hostname", hostname).Msg("Hostname set to ephemeral container port.")

	cleanup := func() error {
		return pool.Purge(resource)
	}

	return resource.GetPort("5432/tcp"), cleanup, nil
}
