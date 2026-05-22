package mocks

import (
	"context"
	"time"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
)

// Compile-time interface check
var _ common.TieredCache = (*MockTieredCache)(nil)

// SetCallRecord tracks calls to Set method.
type SetCallRecord struct {
	Key   string
	Value []byte
	TTL   time.Duration
}

// MockTieredCache is a mock implementation of the TieredCache interface.
type MockTieredCache struct {
	// Data storage for mocking
	Data map[string][]byte

	// Control fields for testing behavior
	GetError        error
	SetError        error
	DeleteError     error
	InvalidateError error
	GetTierHit      int // 0=L0, 1=L1, 2=L2, -1=miss

	// Call tracking
	GetCalls        []string
	SetCalls        []SetCallRecord
	DeleteCalls     []string
	InvalidateCalls []string
}

// NewMockTieredCache creates a new MockTieredCache with initialized data map.
func NewMockTieredCache() *MockTieredCache {
	return &MockTieredCache{
		Data:            make(map[string][]byte),
		GetCalls:        make([]string, 0),
		SetCalls:        make([]SetCallRecord, 0),
		DeleteCalls:     make([]string, 0),
		InvalidateCalls: make([]string, 0),
		GetTierHit:      0,
	}
}

// Get retrieves a value from the mock cache.
func (m *MockTieredCache) Get(ctx context.Context, key string) ([]byte, int, error) {
	m.GetCalls = append(m.GetCalls, key)

	if m.GetError != nil {
		return nil, -1, m.GetError
	}

	data, ok := m.Data[key]
	if !ok {
		return nil, -1, nil
	}

	return data, m.GetTierHit, nil
}

// Set stores a value in the mock cache.
func (m *MockTieredCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	m.SetCalls = append(m.SetCalls, SetCallRecord{Key: key, Value: value, TTL: ttl})

	if m.SetError != nil {
		return m.SetError
	}

	m.Data[key] = value
	return nil
}

// Delete removes a value from the mock cache.
func (m *MockTieredCache) Delete(ctx context.Context, key string) error {
	m.DeleteCalls = append(m.DeleteCalls, key)

	if m.DeleteError != nil {
		return m.DeleteError
	}

	delete(m.Data, key)
	return nil
}

// Invalidate removes a value from local cache only.
func (m *MockTieredCache) Invalidate(ctx context.Context, key string) error {
	m.InvalidateCalls = append(m.InvalidateCalls, key)

	if m.InvalidateError != nil {
		return m.InvalidateError
	}

	delete(m.Data, key)
	return nil
}

// Stats returns mock cache statistics.
func (m *MockTieredCache) Stats() common.LocalCacheStats {
	return common.LocalCacheStats{}
}

// GetFromL0 retrieves directly from RAM cache (mock implementation).
func (m *MockTieredCache) GetFromL0(key string) ([]byte, bool) {
	data, ok := m.Data[key]
	return data, ok
}

// GetFromL1 retrieves directly from disk cache (mock implementation).
func (m *MockTieredCache) GetFromL1(ctx context.Context, key string) ([]byte, error) {
	data, ok := m.Data[key]
	if !ok {
		return nil, nil
	}
	return data, nil
}

// Warmup preloads keys into L0/L1 from L2 (mock implementation - no-op).
func (m *MockTieredCache) Warmup(ctx context.Context, keys []string) error {
	return nil
}
