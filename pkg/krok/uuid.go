package krok

import (
	"github.com/google/uuid"
)

// GenerateResourceID will generate a unique ID for a resource.
func GenerateResourceID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
