// cache/cache.go
package cache

import "errors"

var ErrCacheMiss = errors.New("cache miss") // Добавляем кастомную ошибку для "ключ не найден"
