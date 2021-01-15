package providers

import "github.com/krok-o/krok/pkg/models"

// URLGenerator generates a unique URL for a repository.
type URLGenerator interface {
	Generate(repo *models.Repository) (string, error)
}
