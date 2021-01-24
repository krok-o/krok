package cache

import (
	"sync"
	"time"

	"github.com/krok-o/krok/pkg/models"
)

const (
	// defaultTTL is the number of minutes to wait before purging an authenticated user.
	defaultTTL = 10 * time.Minute
)

// authUser is a user which is authenticated and doesn't need to be checked if it exists or not
// for the duration of TTL.
type authUser struct {
	// constructed by time.Now().Add(TTL).
	ttl  time.Time
	user *models.User
}

// Expired returns if a user's TTL has expired.
func (a *authUser) Expired() bool {
	return time.Now().After(a.ttl)
}

// UserCache is a cache for authenticated users.
type UserCache struct {
	m map[string]*authUser
	sync.RWMutex
}

// NewUserCache creates a new UserCache.
func NewUserCache() *UserCache {
	return &UserCache{
		m: make(map[string]*authUser),
	}
}

// Add adds a user to the cache with a TTL and locking.
func (c *UserCache) Add(email string, id int) {
	c.Lock()
	defer c.Unlock()

	au := &authUser{
		ttl: time.Now().Add(defaultTTL),
		user: &models.User{
			Email: email,
			ID:    id,
		},
	}
	c.m[email] = au
}

// Remove removes a user from the cache.
func (c *UserCache) Remove(email string) {
	c.Lock()
	defer c.Unlock()

	delete(c.m, email)
}

// Has returns whether we already saved the current user or not.
func (c *UserCache) Has(email string) (*authUser, bool) {
	c.RLock()
	defer c.RUnlock()

	v, ok := c.m[email]
	return v, ok
}

// ClearTTL removes old users.
func (c *UserCache) ClearTTL() {
	c.Lock()
	defer c.Unlock()

	// I don't expect more than say, a 1000 users online at a given time.
	for k, u := range c.m {
		// times up, delete the user. which means the user's information will have to be re-fetched from the db.
		if u.Expired() {
			delete(c.m, k)
		}
	}
}
