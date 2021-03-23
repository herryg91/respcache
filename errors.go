package respcache

import "errors"

var ErrCacheEngine = errors.New("cache engine error")
var ErrUnmarshal = errors.New("failed to unmarshal")
var ErrMarshal = errors.New("failed to marshal")
var ErrOutPointer = errors.New("parameter `out` must be a pointer")
var ErrMismatchDataType = errors.New("parameter `out` and `fallback datas` has mismatch data type")
