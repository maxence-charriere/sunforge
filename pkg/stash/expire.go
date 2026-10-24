package stash

import (
	"context"
	"sync"
	"time"
)

type Expire struct {
	// The duration while an item is cached.
	ItemTTL time.Duration

	once  sync.Once
	mutex sync.RWMutex
	items map[string]*memItem
	queue []*memItem
}

func (c *Expire) Get(ctx context.Context, key string) (any, bool) {
	c.once.Do(c.init)
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	i, isCached := c.items[key]
	if !isCached || i.expiresAt.Before(time.Now()) {
		return nil, false
	}
	return i.value, true
}

func (c *Expire) Set(ctx context.Context, key string, i any) {
	c.once.Do(c.init)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if i, isCached := c.items[key]; isCached {
		c.del(i)
	}

	c.expire()
	c.add(&memItem{
		key:       key,
		expiresAt: time.Now().Add(c.ItemTTL),
		value:     i,
	})
}

func (c *Expire) Del(ctx context.Context, key string) {
	c.once.Do(c.init)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if i, isCached := c.items[key]; isCached {
		c.del(i)
		c.expire()
	}
}

func (c *Expire) init() {
	c.items = make(map[string]*memItem)
}

func (c *Expire) add(i *memItem) {
	c.items[i.key] = i
	c.queue = append(c.queue, i)
}

func (c *Expire) del(i *memItem) {
	if !i.deleted {
		delete(c.items, i.key)
		i.value = nil
		i.deleted = true
	}
}

func (c *Expire) expire() {
	now := time.Now()

	i := 0
	for i < len(c.queue) && c.queue[i].isExpired(now) {
		c.del(c.queue[i])
		i++
	}

	copy(c.queue, c.queue[i:])
	c.queue = c.queue[:len(c.queue)-i]
}

type memItem struct {
	key       string
	expiresAt time.Time
	value     any
	deleted   bool
}

func (i *memItem) isExpired(now time.Time) bool {
	return i.deleted || i.expiresAt.Before(now)
}
