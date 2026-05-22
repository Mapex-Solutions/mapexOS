package bootstrap

import (
	"github.com/prometheus/client_golang/prometheus"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/metrics"
	"go.uber.org/dig"
)

// RouterMetrics contains all metrics specific to the router service.
// This struct is provided to the DIG container so any module can inject it.
type RouterMetrics struct {
	Registry *metrics.Registry

	// --- Event Processing (P0) ---

	// Events that completed processing, by status
	EventsProcessed *prometheus.CounterVec
	// End-to-end processing time per event
	EventProcessingDuration prometheus.Histogram
	// Number of messages per batch
	EventsBatchSize prometheus.Histogram

	// --- Message Lifecycle (P0) ---

	// Message lifecycle outcomes (ack/nack/reject)
	MessagesTotal *prometheus.CounterVec

	// --- Asset TieredCache (P0) ---

	// Cache lookups by tier (L0_RAM/L1_Disk/L2_MinIO/Fallback_HTTP/MISS)
	AssetCacheTotal *prometheus.CounterVec
	// Cache lookup latency by tier
	AssetCacheDuration *prometheus.HistogramVec
	// Cache invalidation operations by status
	CacheInvalidationsTotal *prometheus.CounterVec

	// --- Match Evaluation (P1) ---

	// Match evaluations by result (matched/unmatched/no_config)
	MatchEvaluationsTotal *prometheus.CounterVec
	// Total individual match rules evaluated
	MatchRulesEvaluatedTotal prometheus.Counter

	// --- NATS Publish (P0) ---

	// NATS publish attempts by kind and status
	EventsPublished *prometheus.CounterVec
	// NATS publish latency by kind
	PublishDuration *prometheus.HistogramVec

	// --- RouteGroup CRUD (P2) ---

	// RouteGroup CRUD operations by type and result
	RouteGroupOperations *prometheus.CounterVec
	// RouteGroup CRUD operation latency
	RouteGroupOperationDuration *prometheus.HistogramVec
	// Items returned per list query
	RouteGroupListResultsCount prometheus.Histogram

	// --- RouteGroup Cache (P1) ---

	// RouteGroup cache hit/miss ratio
	RouteGroupCacheTotal *prometheus.CounterVec
}

// InitMetrics creates the metrics registry, declares all service-specific
// metrics, and provides the metrics struct to the DIG container.
func InitMetrics(c *dig.Container) {
	c.Provide(func() *RouterMetrics {
		reg := metrics.NewRegistry("router")

		// Go runtime and process collectors — always enabled
		reg.EnableGoCollector()
		reg.EnableProcessCollector()

		m := &RouterMetrics{
			Registry: reg,

			// --- Event Processing ---
			EventsProcessed: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "event",
				Name:      "processed_total",
				Help:      "Events that completed processing, by status",
			}, []string{"status"}),

			EventProcessingDuration: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "event",
				Name:      "processing_duration_seconds",
				Help:      "End-to-end processing time per event",
				Buckets:   prometheus.DefBuckets,
			}),

			EventsBatchSize: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "event",
				Name:      "batch_size",
				Help:      "Number of messages per batch",
				Buckets:   prometheus.DefBuckets,
			}),

			// --- Message Lifecycle ---
			MessagesTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "message",
				Name:      "total",
				Help:      "Message lifecycle outcomes (ack/nack/reject)",
			}, []string{"result"}),

			// --- Asset TieredCache ---
			AssetCacheTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "asset",
				Name:      "cache_total",
				Help:      "Cache lookups by tier",
			}, []string{"tier"}),

			AssetCacheDuration: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "asset",
				Name:      "cache_duration_seconds",
				Help:      "Cache lookup latency by tier",
				Buckets:   prometheus.DefBuckets,
			}, []string{"tier"}),

			CacheInvalidationsTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "cache",
				Name:      "invalidations_total",
				Help:      "Cache invalidation operations by status",
			}, []string{"status"}),

			// --- Match Evaluation ---
			MatchEvaluationsTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "match",
				Name:      "evaluations_total",
				Help:      "Match evaluations by result",
			}, []string{"result"}),

			MatchRulesEvaluatedTotal: reg.NewCounter(metrics.CounterOpts{
				Subsystem: "match",
				Name:      "rules_evaluated_total",
				Help:      "Total individual match rules evaluated",
			}),

			// --- NATS Publish ---
			EventsPublished: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "event",
				Name:      "published_total",
				Help:      "NATS publish attempts by kind and status",
			}, []string{"kind", "status"}),

			PublishDuration: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "publish",
				Name:      "duration_seconds",
				Help:      "NATS publish latency by kind",
				Buckets:   prometheus.DefBuckets,
			}, []string{"kind"}),

			// --- RouteGroup CRUD ---
			RouteGroupOperations: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "routegroup",
				Name:      "operations_total",
				Help:      "RouteGroup CRUD operations by type and result",
			}, []string{"operation", "status"}),

			RouteGroupOperationDuration: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "routegroup",
				Name:      "operation_duration_seconds",
				Help:      "RouteGroup CRUD operation latency",
				Buckets:   prometheus.DefBuckets,
			}, []string{"operation"}),

			RouteGroupListResultsCount: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "routegroup",
				Name:      "list_results_count",
				Help:      "Items returned per list query",
				Buckets:   []float64{0, 1, 5, 10, 25, 50, 100, 250},
			}),

			// --- RouteGroup Cache ---
			RouteGroupCacheTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "routegroup",
				Name:      "cache_total",
				Help:      "RouteGroup cache hit/miss ratio",
			}, []string{"result"}),
		}

		// Initialize label combinations so metrics appear in /metrics from startup
		for _, s := range []string{"success", "error", "no_match"} {
			m.EventsProcessed.WithLabelValues(s)
		}
		for _, r := range []string{"ack", "nack", "reject", "dlq"} {
			m.MessagesTotal.WithLabelValues(r)
		}
		for _, t := range []string{"l0", "l1", "l2", "fallback"} {
			m.AssetCacheTotal.WithLabelValues(t)
			m.AssetCacheDuration.WithLabelValues(t)
		}
		for _, s := range []string{"success", "error"} {
			m.CacheInvalidationsTotal.WithLabelValues(s)
		}
		for _, r := range []string{"match", "no_match"} {
			m.MatchEvaluationsTotal.WithLabelValues(r)
		}
		for _, k := range []string{"workflow", "trigger"} {
			m.PublishDuration.WithLabelValues(k)
			for _, s := range []string{"ok", "error"} {
				m.EventsPublished.WithLabelValues(k, s)
			}
		}
		for _, op := range []string{"create", "update", "delete", "get", "list"} {
			m.RouteGroupOperationDuration.WithLabelValues(op)
			for _, s := range []string{"success", "error"} {
				m.RouteGroupOperations.WithLabelValues(op, s)
			}
		}
		for _, r := range []string{"hit", "miss"} {
			m.RouteGroupCacheTotal.WithLabelValues(r)
		}

		return m
	})

	logger.Info("[APP:BOOTSTRAP] Metrics registry initialized")
}
