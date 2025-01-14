package lattice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/LatticeBCLab/go-lattice/common/types"
	"github.com/LatticeBCLab/go-lattice/lattice/client"
	"github.com/allegro/bigcache/v3"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

// NewMemoryBlockCache 初始化一个内存缓存
//
// Parameters:
//   - enable bool: 是否启用缓存
//   - httpApi client.HttpApi
//   - daemonHashExpirationDuration time.Duration: 守护区块哈希的过期时长
//   - lifeDuration time.Duration: 缓存的存活时长
//   - cleanInterval time.Duration: 过期缓存的清理间隔时长
//
// Returns:
//   - BlockCache
func NewMemoryBlockCache(daemonHashExpirationDuration time.Duration, lifeDuration time.Duration, cleanInterval time.Duration) BlockCache {
	memoryCacheApi, err := bigcache.New(context.Background(), NewMemoryCacheConfig(lifeDuration, cleanInterval))
	if err != nil {
		panic(err)
	}
	return &memoryBlockCache{
		enable:                       true,
		memoryCacheApi:               memoryCacheApi,
		daemonHashExpirationDuration: daemonHashExpirationDuration,
	}
}

func newDisabledMemoryBlockCache(httpApi client.HttpApi) BlockCache {
	return &memoryBlockCache{
		enable:  false,
		httpApi: httpApi,
	}
}

type BlockCache interface {
	SetHttpApi(httpApi client.HttpApi)

	// SetBlock 设置区块缓存
	//
	// Parameters:
	//   - key string: 缓存的Key
	//   - block *types.LatestBlock: 缓存的区块
	//
	// Returns:
	//   - error
	SetBlock(chainId, address string, block *types.LatestBlock) error

	// GetBlock 获取区块缓存
	//
	// Parameters:
	//   - key string: 缓存的Key
	//
	// Returns:
	//   - *types.LatestBlock: 缓存的区块信息
	//   - error
	GetBlock(chainId, address string) (*types.LatestBlock, error)
}

// type redisBlockCache struct{}

type memoryBlockCache struct {
	enable                       bool               // 是否启用缓存
	httpApi                      client.HttpApi     // 节点的http客户端
	memoryCacheApi               *bigcache.BigCache // big cache
	daemonHashExpireAtMap        sync.Map           // 守护哈希的过期时间，每个链维护一个
	daemonHashExpirationDuration time.Duration      // 守护区块哈希的过期时长
}

func (c *memoryBlockCache) SetHttpApi(httpApi client.HttpApi) {
	c.httpApi = httpApi
}

func (c *memoryBlockCache) SetBlock(chainId, address string, block *types.LatestBlock) error {
	if !c.enable {
		return nil
	}
	bytes, err := json.Marshal(block)
	if err != nil {
		log.Error().Err(err).Msgf("json序列化block失败，chainId: %s, accountAddress: %s", chainId, address)
		return err
	}
	if err := c.memoryCacheApi.Set(fmt.Sprintf("%s_%s", chainId, address), bytes); err != nil {
		log.Error().Err(err).Msgf("设置区块缓存信息失败，chainId: %s, accountAddress: %s", chainId, address)
		return err
	}

	_, ok := c.daemonHashExpireAtMap.Load(chainId)
	if !ok {
		c.daemonHashExpireAtMap.Store(chainId, time.Now().Add(c.daemonHashExpirationDuration))
	}

	return nil
}

func (c *memoryBlockCache) GetBlock(chainId, address string) (*types.LatestBlock, error) {
	if !c.enable {
		return c.httpApi.GetLatestBlock(context.Background(), chainId, address)
	}
	// load cached block from memory
	cacheBlockBytes, err := c.memoryCacheApi.Get(fmt.Sprintf("%s_%s", chainId, address))
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return c.httpApi.GetLatestBlock(context.Background(), chainId, address)
		}
		log.Error().Err(err).Msgf("获取区块缓存信息失败，chainId: %s, accountAddress: %s", chainId, address)
		return nil, err
	}
	cacheBlock := new(types.LatestBlock)
	if err := json.Unmarshal(cacheBlockBytes, cacheBlock); err != nil {
		log.Error().Err(err).Msgf("json序列化block失败，chainId: %s, accountAddress: %s", chainId, address)
		return nil, err
	}
	// judge daemon hash expiration time
	daemonHashExpireAt, ok := c.daemonHashExpireAtMap.Load(chainId)
	if !ok {
		daemonHashExpireAt = time.Now().Add(c.daemonHashExpirationDuration)
		c.daemonHashExpireAtMap.LoadOrStore(chainId, daemonHashExpireAt)
	}
	if time.Now().After(daemonHashExpireAt.(time.Time)) {
		log.Debug().Msgf("守护区块哈希已过期，开始更新守护区块哈希，chainId: %s, accountAddress: %s", chainId, address)
		block, err := c.httpApi.GetLatestBlock(context.Background(), chainId, address)
		if err != nil {
			log.Error().Err(err).Msgf("请求节点获取最新区块信息失败，chainId: %s, accountAddress: %s", chainId, address)
			return nil, err
		}
		c.daemonHashExpireAtMap.Store(chainId, time.Now().Add(c.daemonHashExpirationDuration))
		cacheBlock.DaemonBlockHash = block.DaemonBlockHash
	}

	return cacheBlock, nil
}
