package cachee

import (
	"sync"
	"time"
)

var mu sync.RWMutex = sync.RWMutex{}
var cache map[string]interface{} = make(map[string]interface{})

func Get(k string) (v interface{}, ok bool) {
	mu.RLock()
	defer mu.RUnlock()
	v, ok = cache[k]
	return
}

func Set(k string, v interface{}, expire time.Duration) {
	mu.Lock()
	defer mu.Unlock()

	cache[k] = v

	time.AfterFunc(expire, func() {
		mu.Lock()
		defer mu.Unlock()
		delete(cache, k)
	})
}

func Delete(k string) {
	mu.Lock()
	defer mu.Unlock()
	delete(cache, k)
	return
}

func GetIfNotSet(k string, v interface{}, expire time.Duration) (interface{}, bool) {
	if v, ok := cache[k]; ok {
		return v, ok
	}
	Set(k, v, expire)
	return v, false
}

func Keys() (keys []string) {
	mu.RLock()
	defer mu.RUnlock()
	for k := range cache {
		keys = append(keys, k)
	}
	return
}

func Values() (vals []interface{}) {
	mu.RLock()
	defer mu.RUnlock()
	for _, v := range cache {
		vals = append(vals, v)
	}
	return
}
