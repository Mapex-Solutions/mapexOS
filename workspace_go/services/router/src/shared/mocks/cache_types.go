package mocks

import "time"

// CacheSetCallRecord tracks calls to Set.
type CacheSetCallRecord struct {
	Key   string
	Value interface{}
}

// CacheSetExCallRecord tracks calls to SetEx.
type CacheSetExCallRecord struct {
	Key   string
	Value interface{}
	TTL   time.Duration
}
