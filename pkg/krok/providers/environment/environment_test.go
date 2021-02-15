package environment

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestDockerConverter_LoadValueFromFile(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	location, _ := ioutil.TempDir("", "TestDockerConverter_LoadValueFromFile")
	f := filepath.Join(location, "test_env")
	err := ioutil.WriteFile(f, []byte("getthis"), os.ModePerm)
	assert.NoError(t, err)

	d := DockerConverter{
		Dependencies: Dependencies{
			Logger: logger,
		},
		prefix: location,
	}

	value, err := d.LoadValueFromFile(f)
	assert.NoError(t, err)
	assert.Equal(t, "getthis", value)
}

func TestDockerConverter_LoadValueFromFileJustValue(t *testing.T) {
	d := NewDockerConverter(Dependencies{})

	value, err := d.LoadValueFromFile("getthis")
	assert.NoError(t, err)
	assert.Equal(t, "getthis", value)
}

func TestDockerConverter_LoadValueFromFileMissingFile(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	location, _ := ioutil.TempDir("", "TestDockerConverter_LoadValueFromFileMissingFile")

	d := DockerConverter{
		Dependencies: Dependencies{
			Logger: logger,
		},
		prefix: location,
	}

	_, err := d.LoadValueFromFile(filepath.Join(location, "missing"))
	assert.Error(t, err)
}
