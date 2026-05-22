# 18 — Error Handler Per-Node

**Status:** Planned (post-launch)
**Priority:** P1
**Depends on:** Runtime service operational, Node output system (#7)

---

## Concept

Each node gets an optional `errorHandler` field in its Config map. The user configures retry strategy, timeouts, and notifications per node via the UI.

```json
{
  "errorHandler": {
    "enabled": true,
    "onError": "retry",
    "maxRetries": 3,
    "retryBackoff": "exponential",
    "initialDelay": 1000,
    "maxDelay": 30000,
    "backoffMultiplier": 2.0,
    "timeout": 60000,
    "notifyOnRetry": false,
    "notifyOnExhausted": true,
    "notificationChannel": "slack",
    "notificationTarget": "#workflow-alerts",
    "fallbackAction": "continue",
    "fallbackGotoLabel": ""
  }
}
```

---

## Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Toggle error handler for this node |
| `onError` | enum | `"stop"` | Action on error: `"stop"` \| `"continue"` \| `"retry"` \| `"goto"` |
| `maxRetries` | int | `3` | Max retry attempts (1-10) |
| `retryBackoff` | enum | `"exponential"` | Backoff strategy: `"fixed"` \| `"linear"` \| `"exponential"` |
| `initialDelay` | int (ms) | `1000` | Delay before first retry |
| `maxDelay` | int (ms) | `30000` | Cap on backoff growth |
| `backoffMultiplier` | float | `2.0` | For exponential: delay *= multiplier each retry |
| `timeout` | int (ms) | `60000` | Total timeout across all retries (0 = no limit) |
| `notifyOnRetry` | bool | `false` | Send notification on each retry attempt |
| `notifyOnExhausted` | bool | `true` | Send notification when retries exhausted |
| `notificationChannel` | enum | `""` | Channel: `"slack"` \| `"email"` \| `"webhook"` |
| `notificationTarget` | string | `""` | Channel/email/URL target |
| `fallbackAction` | enum | `"stop"` | After retries exhausted: `"stop"` \| `"continue"` \| `"goto"` |
| `fallbackGotoLabel` | string | `""` | If fallbackAction=goto, which goto receiver label |

---

## Retry Backoff Strategies

### Fixed
```
delay = initialDelay
attempt 1: wait 1000ms
attempt 2: wait 1000ms
attempt 3: wait 1000ms
```

### Linear
```
delay = initialDelay * attemptNumber
attempt 1: wait 1000ms
attempt 2: wait 2000ms
attempt 3: wait 3000ms
```

### Exponential
```
delay = min(initialDelay * (multiplier ^ (attempt-1)), maxDelay)
attempt 1: wait 1000ms
attempt 2: wait 2000ms
attempt 3: wait 4000ms
attempt 4: wait 8000ms (capped at maxDelay=30000ms)
```

---

## UI Design

New collapsible section in every node's config panel: **"Error Handling"**

### Layout
```
+------------------------------------------+
| Error Handling                      [v]  |
+------------------------------------------+
| Enable error handler        [ toggle ]   |
|                                          |
| On Error:  [dropdown: stop/continue/     |
|             retry/goto]                  |
|                                          |
| --- Retry Settings (if onError=retry) ---|
| Max retries:        [input: 3]           |
| Backoff strategy:   [dropdown]           |
| Initial delay (ms): [input: 1000]        |
| Max delay (ms):     [input: 30000]       |
| Backoff multiplier: [input: 2.0]         |
| Total timeout (ms): [input: 60000]       |
|                                          |
| --- After Retries Exhausted -------------|
| Fallback action:    [dropdown: stop/     |
|                      continue/goto]      |
| Goto label:         [input] (if goto)    |
|                                          |
| --- Notifications ---------------------- |
| Notify on retry:       [ toggle ]        |
| Notify on exhausted:   [ toggle ]        |
| Channel: [dropdown: slack/email/webhook] |
| Target:  [input: #channel / email / URL] |
+------------------------------------------+
```

### Behavior
- Section collapsed by default
- When `enabled=false`, all fields disabled/greyed out
- Retry fields only visible when `onError=retry`
- Goto label field only visible when `onError=goto` or `fallbackAction=goto`
- Notification fields only visible when at least one notify toggle is on

---

## Backend Changes (Future Implementation)

### 1. RuntimeService — Retry Loop

Wrap `executor.Execute()` with retry logic in the runtime service:

```go
func (s *RuntimeService) executeWithRetry(ctx context.Context, node ExecutionNode, state *WorkflowState) error {
    cfg := parseErrorHandler(node.Config)
    if !cfg.Enabled || cfg.OnError != "retry" {
        return s.executeNode(ctx, node, state)
    }

    var lastErr error
    deadline := time.Now().Add(time.Duration(cfg.Timeout) * time.Millisecond)

    for attempt := 1; attempt <= cfg.MaxRetries; attempt++ {
        lastErr = s.executeNode(ctx, node, state)
        if lastErr == nil {
            return nil // success
        }

        // Check total timeout
        if cfg.Timeout > 0 && time.Now().After(deadline) {
            break
        }

        // Notify on retry
        if cfg.NotifyOnRetry {
            s.sendNotification(cfg, node, attempt, lastErr)
        }

        // Wait with backoff
        delay := computeDelay(cfg, attempt)
        time.Sleep(delay)
    }

    // Retries exhausted
    if cfg.NotifyOnExhausted {
        s.sendNotification(cfg, node, cfg.MaxRetries, lastErr)
    }

    // Execute fallback
    switch cfg.FallbackAction {
    case "continue":
        return nil
    case "goto":
        return &GotoError{Label: cfg.FallbackGotoLabel}
    default: // "stop"
        return lastErr
    }
}
```

### 2. NodeState — Track Retry Count

Add to `NodeState`:
```go
type NodeState struct {
    // ... existing fields
    RetryCount int       `json:"retryCount,omitempty"`
    LastError  string    `json:"lastError,omitempty"`
}
```

### 3. Notification System

Publish notification events to NATS:
```
Subject: workflow.notification.<orgId>
Payload: {
  "type": "node_retry" | "node_retries_exhausted",
  "instanceId": "...",
  "nodeId": "...",
  "nodeName": "...",
  "attempt": 3,
  "maxRetries": 5,
  "error": "connection timeout",
  "channel": "slack",
  "target": "#workflow-alerts"
}
```

A separate notification consumer processes these events and dispatches to Slack/email/webhook.

### 4. Validation

Add to `node_validator.go` (when implementing):
```go
func validateErrorHandler(cfg map[string]interface{}) []string {
    eh := mapget.Map(cfg, "errorHandler")
    if eh == nil || !mapget.Bool(eh, "enabled") {
        return nil
    }
    // Validate onError enum, maxRetries range, delay values, etc.
}
```

---

## Migration

No migration needed — `errorHandler` is a new optional field in Config map. Existing nodes without it behave as `onError=stop` (current behavior).
