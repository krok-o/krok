package providers

import (
	"github.com/google/uuid"
)

// UUIDGenerator generates UUIDs.
type UUIDGenerator interface {
	Generate() (string, error)
}

// Generator defines a wrapper for uuid generator for testing purposes.
type Generator struct{}

// NewUUIDGenerator creates a new Generator.
func NewUUIDGenerator() *Generator {
	return &Generator{}
}

// Generate will generate a unique ID for a resource.
func (u *Generator) Generate() (string, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return "", nil
	}

	return uid.String(), nil
}
