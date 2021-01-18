package providers

import (
	"github.com/google/uuid"
)

// UUIDGenerator generates UUIDs.
type UUIDGenerator interface {
	Generate() (string, error)
}

type uuidGenerator struct{}

// NewUUIDGenerator creates a new uuidGenerator.
func NewUUIDGenerator() *uuidGenerator {
	return &uuidGenerator{}
}

// Generate will generate a unique ID for a resource.
func (u *uuidGenerator) Generate() (string, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return "", nil
	}

	return uid.String(), nil
}
