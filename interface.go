package respcache

type CacheFallback func() (interface{}, error)
type RespCache interface {
	Run(key string, ttl int, out interface{}, fallbackFn CacheFallback) (iscached bool, err error)
}
