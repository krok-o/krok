package cache

import (
	"sync"
	"time"

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

const (
	// defaultTTL is the number of minutes to wait before purging an authenticated user.
	defaultTTL = 10 * time.Minute
)

// UserCache is a cache for authenticated users.
type UserCache struct {
	m map[string]*providers.UserCacheItem
	sync.RWMutex
}

// NewUserCache creates a new UserCache.
func NewUserCache() *UserCache {
	return &UserCache{
		m: make(map[string]*providers.UserCacheItem),
	}
}

// Add adds a user to the cache with a TTL and locking.
func (c *UserCache) Add(email string, id int) {
	c.Lock()
	defer c.Unlock()

	au := &providers.UserCacheItem{
		TTL: time.Now().Add(defaultTTL),
		User: &models.User{
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
func (c *UserCache) Has(email string) (*providers.UserCacheItem, bool) {
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
