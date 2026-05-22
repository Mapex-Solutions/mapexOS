package bootstrap

import (
	"github.com/prometheus/client_golang/prometheus"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/metrics"
	"go.uber.org/dig"
)

// TriggerMetrics contains all metrics specific to the triggers service.
// This struct is provided to the DIG container so any module can inject it.
type TriggerMetrics struct {
	Registry *metrics.Registry

	// --- Trigger Execution (P0) ---

	// Triggers that completed processing, by status (success/error/disabled/no_executor)
	TriggersProcessed *prometheus.CounterVec
	// End-to-end processing time per trigger
	TriggerProcessingDuration prometheus.Histogram
	// Number of messages per batch
	TriggersBatchSize prometheus.Histogram

	// --- Message Lifecycle (P0) ---

	// Message lifecycle outcomes (ack/nack/reject/dlq)
	MessagesTotal *prometheus.CounterVec

	// --- Executor (P0) ---

	// Executor execution duration by type (http/email/mqtt/...)
	ExecutorDuration *prometheus.HistogramVec
	// Executor invocations by type and status
	ExecutorTotal *prometheus.CounterVec

	// --- NATS Publish (P1) ---

	// NATS publish attempts to events.trigger by status (ok/error)
	EventsPublished *prometheus.CounterVec
	// NATS publish latency to events.trigger
	PublishDuration prometheus.Histogram

	// --- Trigger Config (P1) ---

	// Trigger cache lookups by result (hit/miss)
	TriggerCacheTotal *prometheus.CounterVec
	// Placeholder resolutions by status (success/error)
	PlaceholderResolutions *prometheus.CounterVec

	// --- Workflow Execution (P0) ---

	// Workflow plugin/trigger requests processed, by mode and status
	WorkflowProcessed *prometheus.CounterVec
	// End-to-end processing time per workflow request
	WorkflowProcessingDuration *prometheus.HistogramVec
	// Batch size for workflow execution consumer
	WorkflowBatchSize prometheus.Histogram
	// Workflow message lifecycle outcomes (ack/nack/reject)
	WorkflowMessagesTotal *prometheus.CounterVec
	// Workflow executor duration by action type
	WorkflowExecutorDuration *prometheus.HistogramVec
	// Workflow executor invocations by action type and status
	WorkflowExecutorTotal *prometheus.CounterVec
	// Workflow resume callback publish by status
	WorkflowResumePublished *prometheus.CounterVec
}

// InitMetrics creates the metrics registry, declares all service-specific
// metrics, and provides the metrics struct to the DIG container.
func InitMetrics(c *dig.Container) {
	c.Provide(func() *TriggerMetrics {
		reg := metrics.NewRegistry("triggers")

		// Go runtime and process collectors — always enabled
		reg.EnableGoCollector()
		reg.EnableProcessCollector()

		m := &TriggerMetrics{
			Registry: reg,

			// --- Trigger Execution ---
			TriggersProcessed: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "trigger",
				Name:      "processed_total",
				Help:      "Triggers that completed processing, by status",
			}, []string{"status"}),

			TriggerProcessingDuration: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "trigger",
				Name:      "processing_duration_seconds",
				Help:      "End-to-end processing time per trigger",
				Buckets:   prometheus.DefBuckets,
			}),

			TriggersBatchSize: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "trigger",
				Name:      "batch_size",
				Help:      "Number of messages per batch",
				Buckets:   prometheus.DefBuckets,
			}),

			// --- Message Lifecycle ---
			MessagesTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "message",
				Name:      "total",
				Help:      "Message lifecycle outcomes (ack/nack/reject/dlq)",
			}, []string{"result"}),

			// --- Executor ---
			ExecutorDuration: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "executor",
				Name:      "duration_seconds",
				Help:      "Executor execution duration by type",
				Buckets:   prometheus.DefBuckets,
			}, []string{"type"}),

			ExecutorTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "executor",
				Name:      "total",
				Help:      "Executor invocations by type and status",
			}, []string{"type", "status"}),

			// --- NATS Publish ---
			EventsPublished: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "event",
				Name:      "published_total",
				Help:      "NATS publish attempts to events.trigger by status",
			}, []string{"status"}),

			PublishDuration: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "publish",
				Name:      "duration_seconds",
				Help:      "NATS publish latency to events.trigger",
				Buckets:   prometheus.DefBuckets,
			}),

			// --- Trigger Config ---
			TriggerCacheTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "trigger",
				Name:      "cache_total",
				Help:      "Trigger cache lookups by result",
			}, []string{"result"}),

			PlaceholderResolutions: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "placeholder",
				Name:      "resolutions_total",
				Help:      "Placeholder resolutions by status",
			}, []string{"status"}),

			// --- Workflow Execution ---
			WorkflowProcessed: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "workflow",
				Name:      "processed_total",
				Help:      "Workflow plugin/trigger requests processed, by mode and status",
			}, []string{"mode", "status"}),

			WorkflowProcessingDuration: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "workflow",
				Name:      "processing_duration_seconds",
				Help:      "End-to-end processing time per workflow request",
				Buckets:   prometheus.DefBuckets,
			}, []string{"mode"}),

			WorkflowBatchSize: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "workflow",
				Name:      "batch_size",
				Help:      "Batch size for workflow execution consumer",
				Buckets:   prometheus.DefBuckets,
			}),

			WorkflowMessagesTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "workflow_message",
				Name:      "total",
				Help:      "Workflow message lifecycle outcomes (ack/nack/reject)",
			}, []string{"result"}),

			WorkflowExecutorDuration: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "workflow_executor",
				Name:      "duration_seconds",
				Help:      "Workflow executor duration by action type",
				Buckets:   prometheus.DefBuckets,
			}, []string{"type"}),

			WorkflowExecutorTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "workflow_executor",
				Name:      "total",
				Help:      "Workflow executor invocations by action type and status",
			}, []string{"type", "status"}),

			WorkflowResumePublished: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "workflow_resume",
				Name:      "published_total",
				Help:      "Workflow resume callback publish by status",
			}, []string{"status"}),
		}

		// Initialize label combinations so metrics appear in /metrics from startup
		for _, s := range []string{"success", "error", "disabled", "no_executor"} {
			m.TriggersProcessed.WithLabelValues(s)
		}
		for _, r := range []string{"ack", "nack", "reject", "dlq"} {
			m.MessagesTotal.WithLabelValues(r)
		}
		for _, s := range []string{"ok", "error"} {
			m.EventsPublished.WithLabelValues(s)
		}
		for _, r := range []string{"hit", "miss"} {
			m.TriggerCacheTotal.WithLabelValues(r)
		}
		for _, s := range []string{"success", "error"} {
			m.PlaceholderResolutions.WithLabelValues(s)
		}
		for _, t := range []string{"http", "email", "mqtt", "nats", "websocket", "sms", "telegram", "workflow"} {
			for _, s := range []string{"success", "error"} {
				m.ExecutorTotal.WithLabelValues(t, s)
			}
			m.ExecutorDuration.WithLabelValues(t)
		}

		// Workflow execution label pre-init
		for _, mode := range []string{"plugin", "trigger"} {
			for _, s := range []string{"success", "error"} {
				m.WorkflowProcessed.WithLabelValues(mode, s)
			}
			m.WorkflowProcessingDuration.WithLabelValues(mode)
		}
		for _, r := range []string{"ack", "nack", "reject"} {
			m.WorkflowMessagesTotal.WithLabelValues(r)
		}
		for _, t := range []string{"http", "mqtt", "nats", "email", "rabbitmq", "websocket"} {
			for _, s := range []string{"success", "error"} {
				m.WorkflowExecutorTotal.WithLabelValues(t, s)
			}
			m.WorkflowExecutorDuration.WithLabelValues(t)
		}
		for _, s := range []string{"ok", "error"} {
			m.WorkflowResumePublished.WithLabelValues(s)
		}

		return m
	})

	logger.Info("[APP:BOOTSTRAP] Metrics registry initialized")
}
