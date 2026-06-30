package cache

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"
)

type Cache struct {
	cache *cache.Cache
	ttl   time.Duration
}

func New() *Cache {
	ttlMinutes := 30
	if ttl := os.Getenv("CACHE_TTL_MINUTES"); ttl != "" {
		if v, err := strconv.Atoi(ttl); err == nil && v > 0 {
			ttlMinutes = v
		}
	}

	ttl := time.Duration(ttlMinutes) * time.Minute

	c := cache.New(ttl, 10*time.Minute)

	if ttlMinutes > 0 {
		log.Printf("[Cache] Enabled: TTL=%d minutes", ttlMinutes)
	} else {
		log.Println("[Cache] Disabled")
	}

	return &Cache{
		cache: c,
		ttl:   ttl,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	if c.ttl == 0 {
		return nil, false
	}
	return c.cache.Get(key)
}

func (c *Cache) Set(key string, value interface{}) {
	if c.ttl == 0 {
		return
	}
	c.cache.Set(key, value, c.ttl)
}
func (c *Cache) Delete(key string) {
	c.cache.Delete(key)
}
func (c *Cache) Flush() {
	c.cache.Flush()
	log.Println("[Cache] Flushed")
}
func (c *Cache) Keys() []string {
	items := c.cache.Items()
	keys := make([]string, 0, len(items))
	for k := range items {
		keys = append(keys, k)
	}
	return keys
}
func (c *Cache) GetTTL() time.Duration {
	return c.ttl
}
func (c *Cache) IsEnabled() bool {
	return c.ttl > 0
}
