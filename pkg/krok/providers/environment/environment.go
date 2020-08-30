package environment

import (
	"io/ioutil"
	"strings"

	"github.com/rs/zerolog"
)

const (
	dockerSecretPrefix = "/run/secrets"
)

// Config has the configuration options for the environment
type Config struct {
}

// Dependencies defines the dependencies of this commenter
type Dependencies struct {
	Logger zerolog.Logger
}

// DockerConverter is a docker environment secret converter.
type DockerConverter struct {
	Config
	Dependencies
}

// NewDockerConverter creates a new DockerConverter.
func NewDockerConverter(cfg Config, deps Dependencies) *DockerConverter {
	return &DockerConverter{Config: cfg, Dependencies: deps}
}

// LoadValueFromFile provides the ability to load a secret from a docker
// mounted secret file if the value contains `/run/secret`.
func (d *DockerConverter) LoadValueFromFile(f string) (string, error) {
	// if we don't have that prefix, simply return the content.
	if !strings.HasPrefix(f, dockerSecretPrefix) {
		return f, nil
	}
	// Load the content from file
	d.Logger.Debug().Str("value", f).Msg("Loading value from secret file.")
	data, err := ioutil.ReadFile(f)
	if err != nil {
		d.Logger.Error().Err(err).Msg("Failed to read docker secret file.")
		return "", err
	}
	return string(data), nil
}
