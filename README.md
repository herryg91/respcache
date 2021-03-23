# respcache

## What Is?
Library to create simple response cache. Currently respcache just support redis as a cache (github.com/gomodule/redigo/redis).

Fyi: Redis is an open source (BSD licensed), in-memory data structure store, used as a **database**, **cache**, and message broker (https://redis.io/). There 's a big difference on using redis as database and redis as a cache. In this library, we treat redis as a cache, therefore usually we call this library in handler/driver layer

## How to Use?
**Install**
```go
go get github.com/herryg91/respcache
```

**respcache: redis**

disclaimer: thundering herd still can happen when handle high traffic + using multiple nodes
```go
import (
    "github.com/gomodule/redigo/redis"
    "github.com/herryg91/respcache"
)

rdsPool := &redis.Pool{...} // init redis

rc := respcache.NewRedisCache(rdsPool)

var output OutputStruct
ttl := 60  // 1 minute
iscached, err := rc.Run("key_redis", ttl, &output, func() (interface{}, error) {
    // implement fallback function if redis key haven't set
    // return data, error
})
```

## Future
- redis respcache with thread safe
- go-cache implementation
- memcached implementation

