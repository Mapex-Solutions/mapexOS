package mocks

import (
	"context"
	"encoding/json"
	"time"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
)

// Compile-time check to ensure MockAppCache implements AppCache.
var _ common.AppCache = (*MockAppCache)(nil)

// MockAppCache is a mock implementation of common.AppCache (Cache + CacheGetOrSetEx).
// Used for service-private Redis cache testing (counter cache, general caching).
type MockAppCache struct {
	// Return values
	SetError        error
	SetExError      error
	GetError        error
	DelError        error
	GetOrSetExError error

	// Stored data for simulating cache behavior
	Store map[string]interface{}

	// Call tracking
	SetCalls        []CacheSetCallRecord
	SetExCalls      []CacheSetExCallRecord
	GetCalls        []string
	DelCalls        []string
	GetOrSetExCalls []common.GetOrSetParams
}

// NewMockAppCache creates a new MockAppCache.
func NewMockAppCache() *MockAppCache {
	return &MockAppCache{
		Store:           make(map[string]interface{}),
		SetCalls:        make([]CacheSetCallRecord, 0),
		SetExCalls:      make([]CacheSetExCallRecord, 0),
		GetCalls:        make([]string, 0),
		DelCalls:        make([]string, 0),
		GetOrSetExCalls: make([]common.GetOrSetParams, 0),
	}
}

// Set mocks the cache Set method.
func (m *MockAppCache) Set(_ context.Context, key string, value interface{}) error {
	m.SetCalls = append(m.SetCalls, CacheSetCallRecord{Key: key, Value: value})

	if m.SetError != nil {
		return m.SetError
	}

	m.Store[key] = value
	return nil
}

// SetEx mocks the cache SetEx method.
func (m *MockAppCache) SetEx(_ context.Context, key string, value interface{}, ttl time.Duration) error {
	m.SetExCalls = append(m.SetExCalls, CacheSetExCallRecord{Key: key, Value: value, TTL: ttl})

	if m.SetExError != nil {
		return m.SetExError
	}

	m.Store[key] = value
	return nil
}

// Get mocks the cache Get method.
func (m *MockAppCache) Get(_ context.Context, key string, dest interface{}) error {
	m.GetCalls = append(m.GetCalls, key)

	if m.GetError != nil {
		return m.GetError
	}

	val, exists := m.Store[key]
	if !exists {
		return m.GetError
	}

	// Simulate JSON marshal/unmarshal like real cache
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Del mocks the cache Del method.
func (m *MockAppCache) Del(_ context.Context, key string) error {
	m.DelCalls = append(m.DelCalls, key)

	if m.DelError != nil {
		return m.DelError
	}

	delete(m.Store, key)
	return nil
}

// GetOrSetEx mocks the cache GetOrSetEx method.
// On cache miss (no stored value), it executes the Callback and stores the result.
func (m *MockAppCache) GetOrSetEx(params common.GetOrSetParams) (any, error) {
	m.GetOrSetExCalls = append(m.GetOrSetExCalls, params)

	if m.GetOrSetExError != nil {
		return nil, m.GetOrSetExError
	}

	// Check if value exists in store
	val, exists := m.Store[params.CacheKey]
	if exists && params.Dest != nil {
		data, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, params.Dest); err != nil {
			return nil, err
		}
		return val, nil
	}

	// Cache miss - execute callback
	if params.Callback != nil {
		result, err := params.Callback()
		if err != nil {
			return nil, err
		}

		// Store in cache and populate Dest
		if result != nil {
			m.Store[params.CacheKey] = result
			if params.Dest != nil {
				data, err := json.Marshal(result)
				if err != nil {
					return nil, err
				}
				if err := json.Unmarshal(data, params.Dest); err != nil {
					return nil, err
				}
			}
		}

		return result, nil
	}

	return nil, nil
}
