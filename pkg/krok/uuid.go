package krok

import (
	"github.com/google/uuid"
)

// GenerateUUID will generate a unique ID for a resource.
func GenerateUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
