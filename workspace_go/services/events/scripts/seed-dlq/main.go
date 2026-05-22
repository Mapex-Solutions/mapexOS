package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	dlqContract "github.com/Mapex-Solutions/MapexOS/contracts/common/dlq"
)

/**
 * Seed DLQ — Publishes test DLQ messages to the canonical DLQ stream
 * (${ENV}-MAPEXOS-MAPEXOSGOKIT-DLQ) for development testing.
 *
 * Usage:
 *   go run ./scripts/seed-dlq
 *   go run ./scripts/seed-dlq --count 50
 *   NATS_URL=nats://192.168.15.5:4222 go run ./scripts/seed-dlq
 */

// DLQMessage mirrors the infrastructure DLQMessage struct
type DLQMessage struct {
	ID              string            `json:"id"`
	EventTrackerId  string            `json:"eventTrackerId"`
	OrgId           string            `json:"orgId"`
	PathKey         string            `json:"pathKey"`
	ServiceName     string            `json:"serviceName"`
	ServiceType     string            `json:"serviceType"`
	EventType       string            `json:"eventType"`
	OriginalSubject string            `json:"originalSubject"`
	OriginalStream  string            `json:"originalStream"`
	OriginalData    json.RawMessage   `json:"originalData"`
	OriginalHeaders map[string]string `json:"originalHeaders,omitempty"`
	LastError       string            `json:"lastError"`
	ErrorCount      int               `json:"errorCount"`
	FirstDelivery   time.Time         `json:"firstDelivery"`
	LastDelivery    time.Time         `json:"lastDelivery"`
	TotalDeliveries int               `json:"totalDeliveries"`
	ConsumerName    string            `json:"consumerName"`
	SentToDLQAt     time.Time         `json:"sentToDLQAt"`
}

// Service scenarios — realistic failure patterns
var scenarios = []struct {
	ServiceName     string
	ServiceType     string
	EventType       string
	OriginalSubject string
	OriginalStream  string
	ConsumerName    string
	Errors          []string
	Payloads        []map[string]interface{}
}{
	{
		ServiceName:     "workflow-service",
		ServiceType:     "workflow",
		EventType:       "workflow.execution",
		OriginalSubject: "mapexos.workflow.execute",
		OriginalStream:  "MAPEXOS-WORKFLOW",
		ConsumerName:    "workflow-service-runtime-execute",
		Errors: []string{
			"context deadline exceeded: node 'http_request_1' timed out after 30s",
			"plugin execution failed: telegram plugin returned status 429 Too Many Requests",
			"subworkflow deadlock: child execution abc-123 never completed callback",
			"CAS conflict: revision mismatch after 5 retries on KV checkpoint",
			"panic recovered: runtime error: index out of range [3] with length 2",
		},
		Payloads: []map[string]interface{}{
			{"workflowUUID": "wf-001", "instanceId": "inst-abc", "definitionId": "def-123", "nodeId": "http_request_1", "attempt": 1},
			{"workflowUUID": "wf-002", "instanceId": "inst-def", "definitionId": "def-456", "nodeId": "telegram_send", "attempt": 3},
			{"workflowUUID": "wf-003", "instanceId": "inst-ghi", "definitionId": "def-789", "nodeId": "subworkflow_1", "attempt": 1},
		},
	},
	{
		ServiceName:     "triggers-service",
		ServiceType:     "triggers",
		EventType:       "trigger.execute",
		OriginalSubject: "mapexos.triggers.execute",
		OriginalStream:  "MAPEXOS-TRIGGERS",
		ConsumerName:    "triggers-service-trigger-execute",
		Errors: []string{
			"HTTP trigger failed: POST https://api.example.com/webhook returned 503 Service Unavailable",
			"MQTT publish failed: connection refused to broker mqtt://10.0.1.50:1883",
			"email trigger failed: SMTP authentication error: invalid credentials for smtp.company.com",
			"Slack trigger failed: channel_not_found: #alerts-prod does not exist",
		},
		Payloads: []map[string]interface{}{
			{"triggerId": "trg-http-001", "triggerType": "http", "category": "technical", "source": "router", "url": "https://api.example.com/webhook"},
			{"triggerId": "trg-mqtt-001", "triggerType": "mqtt", "category": "technical", "source": "router", "topic": "devices/alerts"},
			{"triggerId": "trg-email-001", "triggerType": "email", "category": "communication", "source": "router", "to": "ops@company.com"},
			{"triggerId": "trg-slack-001", "triggerType": "slack", "category": "communication", "source": "router", "channel": "#alerts-prod"},
		},
	},
	{
		ServiceName:     "router-service",
		ServiceType:     "router",
		EventType:       "route.execute",
		OriginalSubject: "mapexos.router.execute",
		OriginalStream:  "MAPEXOS-ROUTER",
		ConsumerName:    "router-service-route-execute",
		Errors: []string{
			"route evaluation failed: nil pointer dereference in condition evaluator for route group 'production-alerts'",
			"publish failed: max inflight messages reached on stream MAPEXOS-TRIGGERS",
			"asset lookup failed: asset 'sensor-temp-042' not found in org org-001",
		},
		Payloads: []map[string]interface{}{
			{"routeGroupId": "rg-001", "assetId": "asset-temp-042", "templateId": "tmpl-sensor", "matchedRoutes": 3},
			{"routeGroupId": "rg-002", "assetId": "asset-pressure-007", "templateId": "tmpl-pressure", "matchedRoutes": 0},
		},
	},
	{
		ServiceName:     "events-service",
		ServiceType:     "events",
		EventType:       "events.store",
		OriginalSubject: "mapexos.events.store",
		OriginalStream:  "MAPEXOS-EVENTS-STORE",
		ConsumerName:    "events-service-events-store",
		Errors: []string{
			"ClickHouse batch insert failed: code: 252, message: Too many parts (600). Merges are processing significantly slower than inserts",
			"ClickHouse connection timeout: dial tcp 10.0.1.100:9000: i/o timeout after 5s",
			"invalid event data: missing required field 'templateId' in processed event",
		},
		Payloads: []map[string]interface{}{
			{"threadId": "thread-001", "assetId": "asset-001", "templateId": "tmpl-001", "eventType": "telemetry", "source": "http_gateway"},
			{"threadId": "thread-002", "assetId": "asset-002", "eventType": "telemetry", "source": "mqtt_gateway"},
		},
	},
	{
		ServiceName:     "assets-service",
		ServiceType:     "assets",
		EventType:       "lists.name_updated",
		OriginalSubject: "mapexos.lists.name_updated",
		OriginalStream:  "MAPEXOS-LISTS",
		ConsumerName:    "assets-service-list-name-updated",
		Errors: []string{
			"MongoDB update failed: connection pool exhausted, no available connections after 10s wait",
			"bulk update timeout: updating 1500 asset templates with new list name exceeded 30s deadline",
		},
		Payloads: []map[string]interface{}{
			{"listId": "list-001", "oldName": "Device Types", "newName": "Equipment Categories", "affectedTemplates": 42},
			{"listId": "list-002", "oldName": "Locations", "newName": "Site Locations", "affectedTemplates": 15},
		},
	},
	{
		ServiceName:     "mapex-iam-service",
		ServiceType:     "mapex-iam",
		EventType:       "cache.invalidation",
		OriginalSubject: "mapexos.cache.invalidation.roles",
		OriginalStream:  "MAPEXOS_CACHE_INVALIDATION",
		ConsumerName:    "mapex-iam-service-cache-invalidation",
		Errors: []string{
			"Redis FLUSHDB failed: READONLY You can't write against a read only replica",
			"cache invalidation timeout: clearing tiered cache for org org-001 exceeded 15s",
		},
		Payloads: []map[string]interface{}{
			{"cacheType": "roles", "orgId": "org-001", "reason": "role_updated", "roleId": "role-admin"},
			{"cacheType": "permissions", "orgId": "org-002", "reason": "group_membership_changed", "groupId": "grp-ops"},
		},
	},
	{
		ServiceName:     "js-executor-service",
		ServiceType:     "js-executor",
		EventType:       "jsexec.execute",
		OriginalSubject: "mapexos.jsexec.execute",
		OriginalStream:  "MAPEXOS-JSEXEC",
		ConsumerName:    "js-executor-service-script-execute",
		Errors: []string{
			"V8 execution timeout: script exceeded 5000ms CPU time limit",
			"V8 memory limit exceeded: script allocated 128MB, limit is 64MB",
			"script runtime error: TypeError: Cannot read properties of undefined (reading 'temperature')",
			"script compilation error: SyntaxError: Unexpected token '}' at line 42",
		},
		Payloads: []map[string]interface{}{
			{"scriptId": "script-001", "assetId": "asset-sensor-001", "templateId": "tmpl-sensor", "executionTime": 5001},
			{"scriptId": "script-002", "assetId": "asset-gateway-003", "templateId": "tmpl-gateway", "memoryUsed": 134217728},
		},
	},
	{
		ServiceName:     "js-workflow-executor",
		ServiceType:     "js-workflow-executor",
		EventType:       "workflow.js.code",
		OriginalSubject: "workflow.js.code",
		OriginalStream:  "WORKFLOW-JS-CODE",
		ConsumerName:    "js-workflow-executor-code",
		Errors: []string{
			"Piscina worker timeout: script execution exceeded 10s deadline",
			"V8 isolate OOM: workflow code node allocated 256MB, limit is 128MB",
			"script runtime error: ReferenceError: context is not defined at line 15",
			"worker thread crashed: SIGABRT during V8 garbage collection",
		},
		Payloads: []map[string]interface{}{
			{"workflowId": "wf-001", "nodeId": "code_1", "instanceId": "inst-abc", "orgId": "org-001"},
			{"workflowId": "wf-002", "nodeId": "code_transform", "instanceId": "inst-def", "orgId": "org-001"},
		},
	},
}

// Real org IDs from MongoDB seed (deployment/docker-compose/.../organizations.json)
var orgIDs = []struct {
	OrgID   string
	PathKey string
}{
	{"0000000000000000000aa001", "000001"},
}

func main() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	natsUser := os.Getenv("NATS_USERNAME")
	if natsUser == "" {
		natsUser = "service"
	}
	natsPass := os.Getenv("NATS_PASSWORD")
	if natsPass == "" {
		natsPass = "service_secret"
	}

	count := 20
	if len(os.Args) > 2 && os.Args[1] == "--count" {
		fmt.Sscanf(os.Args[2], "%d", &count)
	}

	log.Printf("[SEED-DLQ] Connecting to NATS at %s (user: %s)", natsURL, natsUser)

	nc, err := nats.Connect(natsURL, nats.UserInfo(natsUser, natsPass))
	if err != nil {
		log.Fatalf("[SEED-DLQ] Failed to connect: %v", err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("[SEED-DLQ] Failed to get JetStream context: %v", err)
	}

	// Ensure DLQ stream exists (canonical name resolved at startup from GO_ENV).
	_, err = js.StreamInfo(dlqContract.Stream)
	if err != nil {
		log.Printf("[SEED-DLQ] Stream %s not found, creating...", dlqContract.Stream)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     dlqContract.Stream,
			Subjects: []string{dlqContract.Subject},
			Storage:  nats.FileStorage,
		})
		if err != nil {
			log.Fatalf("[SEED-DLQ] Failed to create stream: %v", err)
		}
		log.Printf("[SEED-DLQ] Stream %s created", dlqContract.Stream)
	}

	log.Printf("[SEED-DLQ] Publishing %d test DLQ messages...", count)

	for i := 0; i < count; i++ {
		scenario := scenarios[rand.Intn(len(scenarios))]
		org := orgIDs[rand.Intn(len(orgIDs))]
		errMsg := scenario.Errors[rand.Intn(len(scenario.Errors))]
		payload := scenario.Payloads[rand.Intn(len(scenario.Payloads))]

		// Add org context to payload
		payload["orgId"] = org.OrgID

		payloadBytes, _ := json.Marshal(payload)

		errorCount := rand.Intn(10) + 1
		firstDelivery := time.Now().Add(-time.Duration(rand.Intn(3600)) * time.Second)

		msg := DLQMessage{
			ID:              uuid.New().String(),
			EventTrackerId:  uuid.New().String(),
			OrgId:           org.OrgID,
			PathKey:         org.PathKey,
			ServiceName:     scenario.ServiceName,
			ServiceType:     scenario.ServiceType,
			EventType:       scenario.EventType,
			OriginalSubject: scenario.OriginalSubject,
			OriginalStream:  scenario.OriginalStream,
			OriginalData:    payloadBytes,
			OriginalHeaders: map[string]string{
				"Nats-Msg-Id":    uuid.New().String(),
				"X-Org-Id":       org.OrgID,
				"X-Path-Key":     org.PathKey,
				"Content-Type":   "application/json",
				"X-Retry-Count":  fmt.Sprintf("%d", errorCount),
			},
			LastError:       errMsg,
			ErrorCount:      errorCount,
			FirstDelivery:   firstDelivery,
			LastDelivery:    time.Now(),
			TotalDeliveries: errorCount + rand.Intn(3),
			ConsumerName:    scenario.ConsumerName,
			SentToDLQAt:     time.Now(),
		}

		data, err := json.Marshal(msg)
		if err != nil {
			log.Printf("[SEED-DLQ] Failed to marshal message %d: %v", i+1, err)
			continue
		}

		_, err = js.Publish("dlq.mapexos", data)
		if err != nil {
			log.Printf("[SEED-DLQ] Failed to publish message %d: %v", i+1, err)
			continue
		}

		log.Printf("[SEED-DLQ] [%d/%d] %s → %s (%s)", i+1, count, scenario.ServiceName, scenario.EventType, errMsg[:min(60, len(errMsg))])
	}

	log.Printf("[SEED-DLQ] Done! Published %d DLQ messages to MAPEXOS-DLQ", count)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
