package steps

import (
	"net/http"
)

// Bag keys this package writes. Other packages reading these keys
// import the constants from here.
const (
	// BagKeyTriggerID is the Mongo ObjectID hex of the trigger created
	// by CreateTrigger. Route groups of kind=trigger reference it on
	// Router.Trigger.TriggerId.
	BagKeyTriggerID = "triggers.triggerID"

	// BagKeyTriggerSinkServer holds the *http.Server the test sink
	// step started, so the Compensate path can shut it down. The
	// stored value type-asserts to *http.Server at the consumer site.
	BagKeyTriggerSinkServer = "triggers.sinkServer"

	// BagKeyTriggerSinkHits holds a *atomic.Int64 counter the sink
	// increments on each POST it receives. Asserts can read this
	// directly when the journey wants to bypass the events service
	// round-trip and verify the sink itself observed the POST.
	BagKeyTriggerSinkHits = "triggers.sinkHits"
)

// Compile-time assertion: the sink server stored on the bag is the
// std-lib *http.Server type so consumers can stop it via its public
// API without depending on a custom struct.
var _ *http.Server = (*http.Server)(nil)
