package memory

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jiebutech/cache"
	"github.com/jiebutech/config"
	"github.com/spf13/cast"
	"sync"
	"sync/atomic"
	"time"
)

const DriverName = "memory"

func init() {
	cache.RegisterDriver(DriverName, NewMemoryCache)
}

func NewMemoryCache(conf config.CacheConfig) cache.Cache {
	c := &ins{
		rmKeyChan: make(chan string, 100),
		data:      [2]*sync.Map{{}, {}},
		idx:       0,
		lock:      sync.RWMutex{},
	}
	go c.listenCacheClear()
	return c
}

func (i *ins) listenCacheClear() {
	for {
		key := <-i.rmKeyChan
		_ = i.Delete(context.Background(), key)
	}
}

type ins struct {
	rmKeyChan chan string
	data      [2]*sync.Map
	idx       int32
	lock      sync.RWMutex
}

func (i *ins) dataContainer() *sync.Map {
	//i.lock.RLock()
	//defer i.lock.RUnlock()
	return i.data[i.idx]
}

func (i *ins) Get(ctx context.Context, key string) string {
	value, ok := i.dataContainer().Load(key)
	if !ok {
		return ""
	}
	strValue, ok := value.(string)
	if !ok {
		return ""
	}
	return strValue
}

func (i *ins) GetAndUnmarshal(ctx context.Context, key string, container interface{}) error {
	str := i.Get(ctx, key)
	if str == "" {
		return errors.New("cache value is empty")
	}
	return json.Unmarshal([]byte(str), container)
}

func (i *ins) Set(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error {
	data := i.dataContainer()
	str := cast.ToString(value)
	if str == "" {
		return errors.New("The value is empty ")
	}
	data.Store(key, str)
	if len(ttl) > 0 && ttl[0] > 0 {
		go func() {
			// 定时器删除缓存数据
			timer := time.NewTimer(ttl[0])
			<-timer.C
			i.rmKeyChan <- key
		}()
	}
	return nil
}

func (i *ins) HSet(ctx context.Context, key, field string, value interface{}) error {
	panic("The memory cache only support key-value struct")
}

func (i *ins) HDelete(ctx context.Context, key string, field ...string) error {
	panic("The memory cache only support key-value struct")
}

func (i *ins) Delete(ctx context.Context, key string) error {
	data := i.dataContainer()
	data.Delete(key)
	return nil
}

func (i *ins) Clear(ctx context.Context) {
	oldIdx := i.idx
	// 切换到另外一个container
	atomic.CompareAndSwapInt32(&i.idx, i.idx, i.idx+1)
	// 移除旧的container
	i.data[oldIdx] = &sync.Map{}

}
