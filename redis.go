package respcache

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

type resp_cache_redis struct {
	rdsPool *redis.Pool
}

func NewRedisCache(rdsPool *redis.Pool) RespCache {
	return &resp_cache_redis{
		rdsPool: rdsPool,
	}
}

func (rc *resp_cache_redis) Run(key string, ttl int, out interface{}, fallbackFn CacheFallback) (iscached bool, err error) {
	iscached = false
	err = nil

	/* get from cache */
	cached_datas := ""
	iscached, cached_datas = rc.get(key)
	if iscached && cached_datas != "" {
		err_parse := json.Unmarshal([]byte(cached_datas), &out)
		if err_parse == nil {
			return
		}

		logrus.Warn(fmt.Sprintf("%s: %s | key: %s, cached_datas: %s", ErrUnmarshal.Error(), err_parse.Error(), key, cached_datas))
	}

	/* fallback if not data or fail get from cache */
	iscached = false
	var fallbackDatas interface{}
	fallbackDatas, err = fallbackFn()
	if err != nil {
		return
	} else if fallbackDatas == nil {
		return
	}
	outVal := reflect.ValueOf(out)
	if outVal.Kind() != reflect.Ptr {
		err = fmt.Errorf("%w: %s", ErrOutPointer, "got: "+reflect.TypeOf(out).String())
		return
	} else if outVal.Elem().Type() != reflect.TypeOf(fallbackDatas) {
		err = fmt.Errorf("%w. out: %s, fallback: %s", ErrMismatchDataType, outVal.Elem().Type().String(), reflect.TypeOf(fallbackDatas).String())
		return
	}
	outVal.Elem().Set(reflect.ValueOf(fallbackDatas))

	/* set to cache */
	rc.set(key, ttl, fallbackDatas)
	return
}

func (rc *resp_cache_redis) get(key string) (iscached bool, resp string) {
	rdsConn := rc.rdsPool.Get()
	defer rdsConn.Close()

	iscached = true
	var err error
	resp, err = redis.String(rdsConn.Do("GET", key))
	if err != nil {
		iscached = false
		resp = ""
		if !errors.Is(err, redis.ErrNil) {
			logrus.Warn(fmt.Sprintf("%s (redis): %s | key: %s", ErrCacheEngine.Error(), err.Error(), key))
		}
	}
	return
}

func (rc *resp_cache_redis) set(key string, ttl int, datas interface{}) {
	rdsConn := rc.rdsPool.Get()
	defer rdsConn.Close()
	jsonResp, err := json.Marshal(datas)
	if err != nil {
		logrus.Warn(fmt.Sprintf("%s: %s | datas: %v", ErrMarshal.Error(), err.Error(), datas))
		return
	}

	args := []interface{}{}
	args = append(args, key, string(jsonResp))
	if ttl > 0 {
		args = append(args, "EX", ttl)
	}
	_, err = rdsConn.Do("SET", args...)
	if err != nil {
		logrus.Warn(fmt.Sprintf("Failed to set cached (key: %s): %v", key, err))
	}
}
