package asserts

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/core/saga"
	triggerSteps "github.com/Mapex-Solutions/MapexOS/e2eTests/services/triggers/triggers/steps"
)

// AssertSmtpReceivedEventually polls the SMTP sink hit counter until
// it observes at least `expected` deliveries or the timeout elapses.
// Smoke oracle for the Email trigger journey.
//
// Reads (bag):
//   - triggerSteps.BagKeySmtpHits  *atomic.Int64  set by StartSmtpSink
func AssertSmtpReceivedEventually(expected int64) saga.Assert {
	return AssertSmtpReceivedEventuallyWithTimeout(expected, 90*time.Second, 500*time.Millisecond)
}

// AssertSmtpReceivedEventuallyWithTimeout overrides the polling budget.
func AssertSmtpReceivedEventuallyWithTimeout(expected int64, timeout, tick time.Duration) saga.Assert {
	return saga.Assert{
		Name: fmt.Sprintf("triggers/triggers.AssertSmtpReceivedEventually[%d]", expected),
		Check: func(c *saga.Context) error {
			v, ok := c.Get(triggerSteps.BagKeySmtpHits)
			if !ok {
				return fmt.Errorf("bag key %q missing — did StartSmtpSink run?", triggerSteps.BagKeySmtpHits)
			}
			hits, ok := v.(*atomic.Int64)
			if !ok {
				return fmt.Errorf("bag key %q is not *atomic.Int64 (%T)", triggerSteps.BagKeySmtpHits, v)
			}
			deadline := time.Now().Add(timeout)
			for {
				got := hits.Load()
				if got >= expected {
					return nil
				}
				if time.Now().After(deadline) {
					return fmt.Errorf("smtp hits: want >=%d, got %d after %v", expected, got, timeout)
				}
				select {
				case <-c.Stdctx.Done():
					return fmt.Errorf("smtp-hit poll cancelled: %w", c.Stdctx.Err())
				case <-time.After(tick):
				}
			}
		},
	}
}

// AssertSmtpLastMessageContent validates the most recent captured
// SMTP message against the expected envelope + subject fragment.
// Empty expected* args are skipped — pass only the fields you care
// about. expectedSubjectFragment is a substring match.
//
// Reads (bag):
//   - triggerSteps.BagKeySmtpLastMessage  **SmtpCapturedMessage  set by StartSmtpSink
func AssertSmtpLastMessageContent(expectedFrom, expectedTo, expectedSubjectFragment string) saga.Assert {
	return saga.Assert{
		Name: "triggers/triggers.AssertSmtpLastMessageContent",
		Check: func(c *saga.Context) error {
			v, ok := c.Get(triggerSteps.BagKeySmtpLastMessage)
			if !ok {
				return fmt.Errorf("bag key %q missing — did StartSmtpSink run?", triggerSteps.BagKeySmtpLastMessage)
			}
			slot, ok := v.(**triggerSteps.SmtpCapturedMessage)
			if !ok {
				return fmt.Errorf("bag key %q is not **SmtpCapturedMessage (%T)", triggerSteps.BagKeySmtpLastMessage, v)
			}
			msg := *slot
			if msg == nil {
				return fmt.Errorf("smtp last message: nil — no delivery captured yet")
			}
			if expectedFrom != "" && msg.From != expectedFrom {
				return fmt.Errorf("smtp from: want %q, got %q", expectedFrom, msg.From)
			}
			if expectedTo != "" && msg.To != expectedTo {
				return fmt.Errorf("smtp to: want %q, got %q", expectedTo, msg.To)
			}
			if expectedSubjectFragment != "" && !strings.Contains(msg.Subject, expectedSubjectFragment) {
				return fmt.Errorf("smtp subject: want fragment %q in %q", expectedSubjectFragment, msg.Subject)
			}
			return nil
		},
	}
}
