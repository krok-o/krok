package livestore

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

// createTestContainerIfNotCI uses an ephemeral postgres container to run a real test.
// the cleanup has to be called by the test runner.
func createTestContainerIfNotCI() (func() error, error) {
	if _, ok := os.LookupEnv("CIRCLECI"); ok {
		// skip circleci environment and do nothing on cleanup.
		return func() error { return nil }, nil
	}
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13.1-alpine",
		Env: []string{
			"POSTGRES_USER=krok",
			"POSTGRES_PASSWORD=password123",
			"listen_addresses = '*'",
		},
		Mounts: []string{"dbinit:/docker-entrypoint-initdb.d"},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})

	if err != nil {
		log.Fatal(err)
	}

	if err = pool.Retry(func() error {
		var err error
		db, err := sql.Open("postgres", fmt.Sprintf("postgres://postgres:password123@localhost:%s/%s?sslmode=disable", resource.GetPort("5432/tcp"), "krok"))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	cleanup := func() error {
		return pool.Purge(resource)
	}
	return cleanup, nil
}
