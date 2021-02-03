package providers

import (
	"time"

	"github.com/krok-o/krok/pkg/models"
)

// UserCacheItem is a user which is authenticated and doesn't need to be checked if it exists or not
// for the duration of TTL.
type UserCacheItem struct {
	// constructed by time.Now().Add(TTL).
	TTL  time.Time
	User *models.User
}

// Expired returns if a user's TTL has expired.
func (a *UserCacheItem) Expired() bool {
	return time.Now().After(a.TTL)
}

// UserCache is a cache for models.User's.
type UserCache interface {
	Add(email string, id int)
	Remove(email string)
	Has(email string) (*UserCacheItem, bool)
	ClearTTL()
}
