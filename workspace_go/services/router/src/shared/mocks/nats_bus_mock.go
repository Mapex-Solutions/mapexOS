package mocks

import (
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// Compile-time interface check
var _ natsModel.CorePublisher = (*MockNatsBus)(nil)

// MockNatsBus is a mock implementation of the natsModel.CorePublisher interface.
type MockNatsBus struct {
	// Control fields for testing behavior
	PublishError error
	FlushError   error

	// Call tracking
	PublishCalls []natsModel.PublishCoreConfig
	FlushCalls   int
}

// NewMockNatsBus creates a new MockNatsBus.
func NewMockNatsBus() *MockNatsBus {
	return &MockNatsBus{
		PublishCalls: make([]natsModel.PublishCoreConfig, 0),
	}
}

// PublishCore implements the natsModel.CorePublisher interface.
func (m *MockNatsBus) PublishCore(config natsModel.PublishCoreConfig) error {
	m.PublishCalls = append(m.PublishCalls, config)

	if m.PublishError != nil {
		return m.PublishError
	}

	return nil
}

// FlushConnection implements the natsModel.CorePublisher interface.
func (m *MockNatsBus) FlushConnection() error {
	m.FlushCalls++

	if m.FlushError != nil {
		return m.FlushError
	}

	return nil
}

// GetPublishCallCount returns the number of publish calls made.
func (m *MockNatsBus) GetPublishCallCount() int {
	return len(m.PublishCalls)
}

// GetLastPublishCall returns the last publish call made.
func (m *MockNatsBus) GetLastPublishCall() *natsModel.PublishCoreConfig {
	if len(m.PublishCalls) == 0 {
		return nil
	}
	return &m.PublishCalls[len(m.PublishCalls)-1]
}

// GetPublishCallsBySubject returns all publish calls for a specific subject.
func (m *MockNatsBus) GetPublishCallsBySubject(subject string) []natsModel.PublishCoreConfig {
	calls := make([]natsModel.PublishCoreConfig, 0)
	for _, call := range m.PublishCalls {
		if call.Subject == subject {
			calls = append(calls, call)
		}
	}
	return calls
}

// Reset clears all recorded calls.
func (m *MockNatsBus) Reset() {
	m.PublishCalls = make([]natsModel.PublishCoreConfig, 0)
	m.PublishError = nil
	m.FlushError = nil
	m.FlushCalls = 0
}
