package bootstrap

import (
	"github.com/prometheus/client_golang/prometheus"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/metrics"
	"go.uber.org/dig"
)

// EventsMetrics contains all metrics specific to the events service.
// This struct is provided to the DIG container so any module can inject it.
type EventsMetrics struct {
	Registry *metrics.Registry

	// --- Event Processing (P0) ---

	// Events that completed processing, by consumer and status
	EventsProcessed *prometheus.CounterVec
	// End-to-end processing time per batch, by consumer
	EventProcessingDuration *prometheus.HistogramVec
	// Number of messages per batch, by consumer
	EventsBatchSize *prometheus.HistogramVec

	// --- Message Lifecycle (P0) ---

	// Message lifecycle outcomes by consumer (ack/nack/reject/dlq)
	MessagesTotal *prometheus.CounterVec

	// --- ClickHouse (P0) ---

	// ClickHouse bulk insert duration by table
	ClickHouseInsertDuration *prometheus.HistogramVec
	// ClickHouse bulk insert attempts by table and status
	ClickHouseInsertTotal *prometheus.CounterVec
	// ClickHouse bulk insert batch size by table
	ClickHouseInsertBatchSize *prometheus.HistogramVec

	// --- Template Cache / EVA (P1) ---

	// Template cache lookups by result (hit/miss/error)
	TemplateCacheTotal *prometheus.CounterVec
	// Total EVA fields resolved
	EvaFieldsMapped prometheus.Counter

	// --- Retention (P2) ---

	// Retention cache lookups by result (hit/miss)
	RetentionCacheTotal *prometheus.CounterVec
}

// InitMetrics creates the metrics registry, declares all service-specific
// metrics, and provides the metrics struct to the DIG container.
func InitMetrics(c *dig.Container) {
	c.Provide(func() *EventsMetrics {
		reg := metrics.NewRegistry("events")

		// Go runtime and process collectors — always enabled
		reg.EnableGoCollector()
		reg.EnableProcessCollector()

		m := &EventsMetrics{
			Registry: reg,

			// --- Event Processing ---
			EventsProcessed: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "event",
				Name:      "processed_total",
				Help:      "Events that completed processing, by consumer and status",
			}, []string{"consumer", "status"}),

			EventProcessingDuration: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "event",
				Name:      "processing_duration_seconds",
				Help:      "End-to-end processing time per batch, by consumer",
				Buckets:   prometheus.DefBuckets,
			}, []string{"consumer"}),

			EventsBatchSize: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "event",
				Name:      "batch_size",
				Help:      "Number of messages per batch, by consumer",
				Buckets:   prometheus.DefBuckets,
			}, []string{"consumer"}),

			// --- Message Lifecycle ---
			MessagesTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "message",
				Name:      "total",
				Help:      "Message lifecycle outcomes by consumer (ack/nack/reject/dlq)",
			}, []string{"consumer", "result"}),

			// --- ClickHouse ---
			ClickHouseInsertDuration: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "clickhouse",
				Name:      "insert_duration_seconds",
				Help:      "ClickHouse bulk insert duration by table",
				Buckets:   prometheus.DefBuckets,
			}, []string{"table"}),

			ClickHouseInsertTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "clickhouse",
				Name:      "insert_total",
				Help:      "ClickHouse bulk insert attempts by table and status",
			}, []string{"table", "status"}),

			ClickHouseInsertBatchSize: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "clickhouse",
				Name:      "insert_batch_size",
				Help:      "ClickHouse bulk insert batch size by table",
				Buckets:   prometheus.DefBuckets,
			}, []string{"table"}),

			// --- Template Cache / EVA ---
			TemplateCacheTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "template",
				Name:      "cache_total",
				Help:      "Template cache lookups by result",
			}, []string{"result"}),

			EvaFieldsMapped: reg.NewCounter(metrics.CounterOpts{
				Subsystem: "eva",
				Name:      "fields_mapped_total",
				Help:      "Total EVA fields resolved from templates",
			}),

			// --- Retention ---
			RetentionCacheTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "retention",
				Name:      "cache_total",
				Help:      "Retention cache lookups by result",
			}, []string{"result"}),
		}

		// Initialize label combinations so metrics appear in /metrics from startup
		consumers := []string{"device", "trigger", "workflow", "dlq"}
		tables := []string{"events_device", "events_trigger", "events_workflow", "events_dlq"}

		for _, c := range consumers {
			m.EventProcessingDuration.WithLabelValues(c)
			m.EventsBatchSize.WithLabelValues(c)
			for _, s := range []string{"success", "error"} {
				m.EventsProcessed.WithLabelValues(c, s)
			}
			for _, r := range []string{"ack", "nack", "reject", "dlq"} {
				m.MessagesTotal.WithLabelValues(c, r)
			}
		}
		for _, t := range tables {
			m.ClickHouseInsertDuration.WithLabelValues(t)
			m.ClickHouseInsertBatchSize.WithLabelValues(t)
			for _, s := range []string{"ok", "error"} {
				m.ClickHouseInsertTotal.WithLabelValues(t, s)
			}
		}
		for _, r := range []string{"hit", "miss"} {
			m.TemplateCacheTotal.WithLabelValues(r)
			m.RetentionCacheTotal.WithLabelValues(r)
		}

		return m
	})

	logger.Info("[APP:BOOTSTRAP] Metrics registry initialized")
}
