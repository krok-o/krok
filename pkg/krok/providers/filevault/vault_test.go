package filevault

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestFileStorer_Flow(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	location, _ := ioutil.TempDir("", "TestFileStorer_Flow")
	fileStore, err := NewFileStorer(Config{
		Location: location,
		Key:      "password123",
	}, Dependencies{Logger: logger})
	assert.NoError(t, err)
	err = fileStore.Init()
	assert.NoError(t, err)
	err = fileStore.Write([]byte("data"))
	assert.NoError(t, err)
	b, err := fileStore.Read()
	assert.NoError(t, err)
	assert.Equal(t, []byte("data"), b)

	// init again... read write should not damage data.
	err = fileStore.Init()
	assert.NoError(t, err)
	b, err = fileStore.Read()
	assert.NoError(t, err)
	b = append(b, []byte("newdata")...)
	err = fileStore.Write(b)
	assert.NoError(t, err)
	b, err = fileStore.Read()
	assert.NoError(t, err)
	assert.Equal(t, []byte("datanewdata"), b)

	// We always write the entire data back and forth.
	err = fileStore.Write([]byte("data2"))
	assert.NoError(t, err)
	b, err = fileStore.Read()
	assert.NoError(t, err)
	assert.Equal(t, []byte("data2"), b)
}
