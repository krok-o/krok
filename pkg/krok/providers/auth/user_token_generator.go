package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
)

var stdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

// UserTokenGenerator represents the user personal token generator.
type UserTokenGenerator struct{}

// NewUserTokenGenerator creates a new UserTokenGenerator.
func NewUserTokenGenerator() *UserTokenGenerator {
	return &UserTokenGenerator{}
}

// Generate generates a random, unique, personal access token for a user.
func (utg *UserTokenGenerator) Generate(length int) (string, error) {
	if length == 0 {
		return "", errors.New("length cannot be 0")
	}

	clen := len(stdChars)
	if clen < 2 || clen > 256 {
		return "", errors.New("wrong charset for length")
	}

	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.

	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			return "", fmt.Errorf("random bytes: %w", err)
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				// Skip this number to avoid modulo bias.
				continue
			}
			b[i] = stdChars[c%clen]
			i++
			if i == length {
				return string(b), nil
			}
		}
	}
}
