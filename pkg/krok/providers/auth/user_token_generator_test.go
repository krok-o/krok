package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUserTokenGenerator(t *testing.T) {
	t.Run("successful 60 length", func(t *testing.T) {
		token, err := NewUserTokenGenerator().Generate(60)
		require.NoError(t, err)
		require.Len(t, token, 60)
	})

	t.Run("successful 10 length", func(t *testing.T) {
		token, err := NewUserTokenGenerator().Generate(10)
		require.NoError(t, err)
		require.Len(t, token, 10)
	})

	t.Run("0 returns err", func(t *testing.T) {
		_, err := NewUserTokenGenerator().Generate(0)
		require.EqualError(t, err, "length cannot be 0")
	})
}
