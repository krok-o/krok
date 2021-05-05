package tar

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestUntar(t *testing.T) {
	log := zerolog.New(os.Stderr)
	tmp, err := ioutil.TempDir("", "untar_01")
	assert.NoError(t, err)
	tar := NewTarer(Dependencies{
		Logger: log,
	})
	archive, err := os.Open(filepath.Join("testdata", "test.tar.gz"))
	assert.NoError(t, err)
	err = tar.Untar(tmp, archive)
	assert.NoError(t, err)
	content, err := ioutil.ReadFile(filepath.Join(tmp, "test"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("this is a test file\n"), content)
}
