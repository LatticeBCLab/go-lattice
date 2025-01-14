package lattice

import (
	"github.com/allegro/bigcache/v3"
	"time"
)

func NewMemoryCacheConfig(lifeDuration time.Duration, cleanInterval time.Duration) bigcache.Config {
	config := bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: 1024,

		// time after which entry can be evicted
		// 缓存的过期时间
		LifeWindow: lifeDuration,

		// Interval between removing expired entries (clean up).
		// If set to <= 0 then no action is performed.
		// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
		// 清理过期条目的时间间隔
		CleanWindow: cleanInterval,

		// rps * lifeWindow, used only in initial memory allocation
		// 在 LifeWindow 时间内可能的最大条目数，支持1024个账户的缓存
		MaxEntriesInWindow: 1024,

		// max entry size in bytes, used only in initial memory allocation
		// 每个条目的最大字节数，512byte = 0.5KB
		MaxEntrySize: 512,

		// prints information about additional memory allocation
		// 是否打印内存分配的详细信息
		Verbose: true,

		// cache will not allocate more memory than this limit, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		// 缓存系统的最大内存限制，以MB为单位，达到设置的上限时，新条目会覆盖旧条目，最大8192
		HardMaxCacheSize: 512,

		// callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		// 当最旧的条目被移除时触发的回调函数
		OnRemove: nil,

		// OnRemoveWithReason is a callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A constant representing the reason will be passed through.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		// Ignored if OnRemove is specified.
		// 与 OnRemove 类似，但此回调函数会带有移除条目的原因
		OnRemoveWithReason: nil,
	}
	return config
}
