package application

import "context"

type Invalidator struct {
	cache *CacheService
}

func NewInvalidator(cache *CacheService) *Invalidator {
	return &Invalidator{cache: cache}
}

func (i *Invalidator) Invalidate(ctx context.Context) {
	if i == nil || i.cache == nil {
		return
	}
	_ = i.cache.InvalidateAllOverview(ctx)
}
