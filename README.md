# Cache

### 缓存包

### Quick Start

1. import
> import "github.com/jiebutech/cache/driver/redis"

2. New client
```go
import (
    "github.com/jiebutech/cache/driver/redis"
	"github.com/jiebutech/cache"
)

func main()  {
	cfg := &config.CacheConfig{}
	c := cache.New(redis.DriverName, cfg)
	world := c.Get("hello")
	println(world)
}


```
