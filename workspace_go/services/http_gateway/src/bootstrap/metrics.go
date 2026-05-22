package bootstrap

import (
	"github.com/prometheus/client_golang/prometheus"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/metrics"
	"go.uber.org/dig"
)

// HttpGatewayMetrics contains all metrics specific to the http_gateway service.
// This struct is provided to the DIG container so any module can inject it.
type HttpGatewayMetrics struct {
	Registry *metrics.Registry

	// --- Event Auth (CustomAuthMiddleware) ---

	// Webhook auth attempts by strategy and result
	EventAuthTotal *prometheus.CounterVec
	// Webhook auth validation latency per strategy
	EventAuthDuration *prometheus.HistogramVec
	// Auth failures that triggered security event to events.raw
	EventAuthFailures *prometheus.CounterVec

	// --- Event Ingestion ---

	// Events processed through the full handler flow
	EventsProcessed *prometheus.CounterVec
	// NATS publish attempts by subject and result
	EventsPublished *prometheus.CounterVec
	// Handler processing time excluding auth
	EventProcessingDuration prometheus.Histogram
	// Incoming webhook body size
	EventPayloadSize prometheus.Histogram

	// --- Data Source CRUD ---

	// Data source CRUD operations by type and result
	DsOperations *prometheus.CounterVec
	// Data source CRUD operation latency
	DsOperationDuration *prometheus.HistogramVec
	// Items returned per list query
	DsListResultsCount prometheus.Histogram

	// --- Data Source Cache ---

	// Data source cache hit/miss counter
	DsCacheTotal *prometheus.CounterVec

	// --- Heartbeat (POST /api/v1/heartbeat — TKT-2026-0034 explicit mode) ---

	// POST /api/v1/heartbeat outcomes labeled by status (success|error)
	HeartbeatsTotal *prometheus.CounterVec
	// End-to-end EventService.ProcessHeartbeat duration in seconds
	HeartbeatDuration prometheus.Histogram
}

// InitMetrics creates the metrics registry, declares all service-specific
// metrics, and provides the metrics struct to the DIG container.
func InitMetrics(c *dig.Container) {
	c.Provide(func() *HttpGatewayMetrics {
		reg := metrics.NewRegistry("httpgw")

		// Go runtime and process collectors — always enabled
		reg.EnableGoCollector()
		reg.EnableProcessCollector()

		m := &HttpGatewayMetrics{
			Registry: reg,

			// --- Event Auth ---
			EventAuthTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "event",
				Name:      "auth_total",
				Help:      "Webhook auth attempts by strategy and result",
			}, []string{"auth_type", "result"}),

			EventAuthDuration: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "event",
				Name:      "auth_duration_seconds",
				Help:      "Webhook auth validation latency per strategy",
				Buckets:   prometheus.DefBuckets,
			}, []string{"auth_type"}),

			EventAuthFailures: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "event",
				Name:      "auth_failures_total",
				Help:      "Auth failures that triggered security event to events.raw",
			}, []string{"auth_type"}),

			// --- Event Ingestion ---
			EventsProcessed: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "event",
				Name:      "processed_total",
				Help:      "Events processed through the full handler flow",
			}, []string{"status"}),

			EventsPublished: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "event",
				Name:      "published_total",
				Help:      "NATS publish attempts by subject and result",
			}, []string{"subject", "status"}),

			EventProcessingDuration: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "event",
				Name:      "processing_duration_seconds",
				Help:      "Handler processing time excluding auth",
				Buckets:   prometheus.DefBuckets,
			}),

			EventPayloadSize: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "event",
				Name:      "payload_size_bytes",
				Help:      "Incoming webhook body size",
				Buckets:   []float64{100, 500, 1000, 5000, 10000, 50000, 100000},
			}),

			// --- Data Source CRUD ---
			DsOperations: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "ds",
				Name:      "operations_total",
				Help:      "Data source CRUD operations by type and result",
			}, []string{"operation", "status"}),

			DsOperationDuration: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "ds",
				Name:      "operation_duration_seconds",
				Help:      "Data source CRUD operation latency",
				Buckets:   prometheus.DefBuckets,
			}, []string{"operation"}),

			DsListResultsCount: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "ds",
				Name:      "list_results_count",
				Help:      "Items returned per list query",
				Buckets:   []float64{0, 1, 5, 10, 25, 50, 100, 250},
			}),

			// --- Data Source Cache ---
			DsCacheTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "ds",
				Name:      "cache_total",
				Help:      "Data source cache lookups by result (hit/miss)",
			}, []string{"result"}),

			// --- Heartbeat (TKT-2026-0034 explicit mode) ---
			HeartbeatsTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "heartbeat",
				Name:      "total",
				Help:      "POST /api/v1/heartbeat requests processed, labeled by outcome.",
			}, []string{"status"}),

			HeartbeatDuration: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "heartbeat",
				Name:      "duration_seconds",
				Help:      "End-to-end duration of EventService.ProcessHeartbeat in seconds.",
				Buckets:   prometheus.DefBuckets,
			}),
		}

		// Initialize label combinations so metrics appear in /metrics from startup
		for _, at := range []string{"apikey", "basic", "bearer", "none"} {
			m.EventAuthDuration.WithLabelValues(at)
			m.EventAuthFailures.WithLabelValues(at)
			for _, r := range []string{"success", "error"} {
				m.EventAuthTotal.WithLabelValues(at, r)
			}
		}
		for _, s := range []string{"success", "error"} {
			m.EventsProcessed.WithLabelValues(s)
		}
		for _, subj := range []string{"events.device", "events.trigger"} {
			for _, s := range []string{"ok", "error"} {
				m.EventsPublished.WithLabelValues(subj, s)
			}
		}
		for _, op := range []string{"create", "update", "delete", "get", "list"} {
			m.DsOperationDuration.WithLabelValues(op)
			for _, s := range []string{"success", "error"} {
				m.DsOperations.WithLabelValues(op, s)
			}
		}
		for _, r := range []string{"hit", "miss"} {
			m.DsCacheTotal.WithLabelValues(r)
		}
		for _, s := range []string{"success", "error"} {
			m.HeartbeatsTotal.WithLabelValues(s)
		}

		return m
	})

	logger.Info("[APP:BOOTSTRAP] Metrics registry initialized")
}
