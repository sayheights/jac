package jac

import (
	"context"
	"sync"
	"time"
)

const (
	defaultEvictionInterval = time.Hour * 1
)

type CacheItem struct {
	Response   *Response
	Expiration time.Time
}

func newCacheItem(response *Response, dur time.Duration) *CacheItem {
	return &CacheItem{
		Response:   response,
		Expiration: time.Now().Add(dur),
	}
}

func (c *CacheItem) isExpired() bool {
	return time.Now().After(c.Expiration)
}

type Cache interface {
	Get(key string) *Response
	Set(key string, item *CacheItem)
}

type InMemoryCache struct {
	mu             sync.RWMutex
	items          map[string]*CacheItem
	interval       time.Duration
	cancelEviction context.CancelFunc
}

type InMemoryCacheOption func(cache *InMemoryCache)

func WithEvictionInterval(interval time.Duration) InMemoryCacheOption {
	return func(cache *InMemoryCache) {
		cache.interval = interval
	}
}

func NewInMemoryCache(opts ...InMemoryCacheOption) *InMemoryCache {
	cache := &InMemoryCache{interval: defaultEvictionInterval, items: map[string]*CacheItem{}}
	ctx, cancel := context.WithCancel(context.Background())
	cache.cancelEviction = cancel

	for _, opt := range opts {
		opt(cache)
	}

	go cacheEvictor(ctx, cache)

	return cache
}

func cacheEvictor(ctx context.Context, cache *InMemoryCache) {
	ticker := time.NewTicker(cache.interval)
	for {
		select {
		case <-ticker.C:
			cache.deleteExpired()
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (i *InMemoryCache) Get(key string) *Response {
	i.mu.RLock()
	defer i.mu.RUnlock()

	item, ok := i.items[key]
	if !ok {
		return nil
	}

	if item.isExpired() {
		delete(i.items, key)
		return nil
	}

	return item.Response
}

func (i *InMemoryCache) Set(key string, item *CacheItem) {
	i.mu.Lock()
	i.items[key] = item
	i.mu.Unlock()
}

func (i *InMemoryCache) StopEvictor() {
	i.cancelEviction()
}

func (i *InMemoryCache) deleteExpired() {
	i.mu.Lock()
	for key, item := range i.items {
		if item.isExpired() {
			delete(i.items, key)
		}
	}
	i.mu.Unlock()
}
