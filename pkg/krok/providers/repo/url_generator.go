package repo

import (
	"fmt"
	"net/url"
	"path"
	"strconv"

	"github.com/krok-o/krok/pkg/models"
)

type URLGenerator struct {
	Hostname string
}

func NewURLGenerator(hostname string) *URLGenerator {
	return &URLGenerator{Hostname: hostname}
}

func (ug *URLGenerator) Generate(repo *models.Repository) (string, error) {
	u, err := url.Parse(ug.Hostname)
	if err != nil {
		return "", fmt.Errorf("url parse: %w", err)
	}

	u.Path = path.Join(u.Path, strconv.Itoa(repo.ID), strconv.Itoa(repo.VCS), "callback")
	return u.String(), nil
}
