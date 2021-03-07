package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// UserTokenGenerator represents the user personal token generator.
type UserTokenGenerator struct{}

// NewUserTokenGenerator creates a new UserTokenGenerator.
func NewUserTokenGenerator() *UserTokenGenerator {
	return &UserTokenGenerator{}
}

// Generate generates a random, unique, personal access token for a user.
func (utg *UserTokenGenerator) Generate() (string, error) {
	b := make([]byte, 60)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("rand read: %w", err)
	}
	return hex.EncodeToString(b), nil
}
