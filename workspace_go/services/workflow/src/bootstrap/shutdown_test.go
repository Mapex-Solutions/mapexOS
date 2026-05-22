package bootstrap

import (
	"testing"
)

func TestConsumerRegistry_RegisterNil(t *testing.T) {
	r := &ConsumerRegistry{}
	r.Register(nil) // should not panic or add

	if len(r.consumers) != 0 {
		t.Fatalf("expected 0 consumers after registering nil, got %d", len(r.consumers))
	}
}

func TestConsumerRegistry_StopAll_Empty(t *testing.T) {
	r := &ConsumerRegistry{}
	r.StopAll() // should not panic on empty registry
}
