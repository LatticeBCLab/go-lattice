package lattice

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/rs/zerolog/log"
)

func NewCache(lifeDuration, cleanInterval time.Duration) Cache {
	api, err := bigcache.New(context.Background(), NewMemoryCacheConfig(lifeDuration, cleanInterval))
	if err != nil {
		log.Fatal().Err(err).Msg("初始化内存缓存失败")
	}
	return &cache{api: api}
}

type cache struct {
	api *bigcache.BigCache
}

type Cache interface {
	// Set value to big cache
	Set(key string, value []byte) error
	// Get value from big cache
	Get(key string) ([]byte, error)
}

func (c *cache) Set(key string, value []byte) error {
	if err := c.api.Set(key, value); err != nil {
		log.Error().Err(err).Msgf("新增缓存信息失败：key: %s，value: %s", key, value)
		return err
	}
	return nil
}

func (c *cache) Get(key string) ([]byte, error) {
	value, err := c.api.Get(key)
	if err != nil {
		log.Error().Err(err).Msgf("查询缓存信息失败：key: %s", key)
		return nil, err
	}
	return value, nil
}
