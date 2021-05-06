package providers

import (
	"io"
)

// Tar provides functionality to untar an uploaded tar.gz.
type Tar interface {
	// Untar will untar the contents of a reader.
	Untar(dst string, reader io.Reader) error
}
