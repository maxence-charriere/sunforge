package stash

import (
	"container/list"
	"context"
	"sync"
	"time"
)

// ExpiringCache stores values until their time-to-live expires.
// It is safe for concurrent use and must not be copied after first use.
type ExpiringCache struct {
	// ItemTTL is the duration each value remains cached after Set.
	// Values expire only after this duration has elapsed, not at the exact deadline.
	// Defaults to one minute. Negative durations disable caching.
	// Set ItemTTL before first use and do not modify it afterward.
	ItemTTL time.Duration

	once  sync.Once
	mutex sync.RWMutex
	items map[string]*list.Element
	queue list.List
	now   func() time.Time // Optional clock for deterministic tests; set before use.
}

func (c *ExpiringCache) Get(ctx context.Context, key string) (any, bool) {
	c.once.Do(c.init)
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	e, ok := c.items[key]
	if !ok {
		return nil, false
	}

	i := e.Value.(*item)
	if c.now().After(i.expiresAt) {
		return nil, false
	}
	return i.value, true
}

func (c *ExpiringCache) Set(ctx context.Context, key string, value any) {
	c.once.Do(c.init)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := c.now()
	c.removeExpiredItems(now)
	if c.ItemTTL <= 0 {
		return
	}

	if e, ok := c.items[key]; ok {
		i := e.Value.(*item)
		i.value = value
		i.expiresAt = now.Add(c.ItemTTL)
		c.queue.MoveToBack(e)
		return
	}

	c.items[key] = c.queue.PushBack(&item{
		key:       key,
		expiresAt: now.Add(c.ItemTTL),
		value:     value,
	})
}

func (c *ExpiringCache) Del(ctx context.Context, key string) {
	c.once.Do(c.init)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if e, ok := c.items[key]; ok {
		c.remove(e)
	}
	c.removeExpiredItems(c.now())
}

func (c *ExpiringCache) init() {
	if c.ItemTTL == 0 {
		c.ItemTTL = time.Minute
	}
	if c.now == nil {
		c.now = time.Now
	}
	c.items = make(map[string]*list.Element)
}

func (c *ExpiringCache) remove(e *list.Element) {
	delete(c.items, e.Value.(*item).key)
	c.queue.Remove(e)
}

func (c *ExpiringCache) removeExpiredItems(now time.Time) {
	for e := c.queue.Front(); e != nil; e = c.queue.Front() {
		if !now.After(e.Value.(*item).expiresAt) {
			break
		}
		c.remove(e)
	}
}

type item struct {
	key       string
	expiresAt time.Time
	value     any
}
