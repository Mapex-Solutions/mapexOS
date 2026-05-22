package bootstrap

import (
	"github.com/prometheus/client_golang/prometheus"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/metrics"
	"go.uber.org/dig"
)

// WorkflowMetrics contains all metrics specific to the workflow service.
type WorkflowMetrics struct {
	Registry *metrics.Registry

	/* Definition CRUD */

	// Definition CRUD operations by type and result
	DefinitionOperations *prometheus.CounterVec
	// Definition CRUD operation latency
	DefinitionOperationDuration *prometheus.HistogramVec
	// Items returned per definition list query
	DefinitionListResultsCount prometheus.Histogram
	// Definition cache hit/miss
	DefinitionCacheTotal *prometheus.CounterVec

	/* Plugin CRUD */

	// Plugin CRUD operations by type and result
	PluginOperations *prometheus.CounterVec
	// Plugin CRUD operation latency
	PluginOperationDuration *prometheus.HistogramVec
	// Items returned per plugin list query
	PluginListResultsCount prometheus.Histogram

	/* Cache */

	// Cache invalidation operations by status
	CacheInvalidationsTotal *prometheus.CounterVec

	/* Execution Lifecycle */

	// Completed executions by trigger source
	ExecutionCompletedTotal *prometheus.CounterVec
	// Failed executions by trigger source
	ExecutionFailedTotal *prometheus.CounterVec
	// Started executions by trigger source
	ExecutionStartedTotal *prometheus.CounterVec
	// Total execution time from start to terminal state
	ExecutionDuration *prometheus.HistogramVec
	// Steps executed per execution
	ExecutionSteps prometheus.Histogram
	// Currently in-flight executions
	ExecutionActive prometheus.Gauge

	/* Checkpoint */

	// KV checkpoint (Put/Update) latency
	CheckpointDuration prometheus.Histogram

	/* Dispatch */

	// Dispatch calls by type and result
	DispatchTotal *prometheus.CounterVec

	/* Resilience */

	// CAS conflict retry count in HandleResume
	CASRetriesTotal prometheus.Counter
	// Stale callback token rejections
	TokenRejectionsTotal prometheus.Counter
	// Walker stops due to graceful shutdown
	ShutdownStopsTotal prometheus.Counter
}

// InitMetrics creates the metrics registry, declares all service-specific
// metrics, and provides the metrics struct to the DIG container.
func InitMetrics(c *dig.Container) {
	c.Provide(func() *WorkflowMetrics {
		reg := metrics.NewRegistry("workflow")

		// Go runtime and process collectors — always enabled
		reg.EnableGoCollector()
		reg.EnableProcessCollector()

		m := &WorkflowMetrics{
			Registry: reg,

			/* Definition CRUD */
			DefinitionOperations: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "definition",
				Name:      "operations_total",
				Help:      "Definition CRUD operations by type and result",
			}, []string{"operation", "status"}),

			DefinitionOperationDuration: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "definition",
				Name:      "operation_duration_seconds",
				Help:      "Definition CRUD operation latency",
				Buckets:   prometheus.DefBuckets,
			}, []string{"operation"}),

			DefinitionListResultsCount: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "definition",
				Name:      "list_results_count",
				Help:      "Items returned per definition list query",
				Buckets:   []float64{0, 1, 5, 10, 25, 50, 100, 250},
			}),

			DefinitionCacheTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "definition",
				Name:      "cache_total",
				Help:      "Definition cache lookups by result (hit/miss)",
			}, []string{"result"}),

			/* Plugin CRUD */
			PluginOperations: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "plugin",
				Name:      "operations_total",
				Help:      "Plugin CRUD operations by type and result",
			}, []string{"operation", "status"}),

			PluginOperationDuration: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "plugin",
				Name:      "operation_duration_seconds",
				Help:      "Plugin CRUD operation latency",
				Buckets:   prometheus.DefBuckets,
			}, []string{"operation"}),

			PluginListResultsCount: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "plugin",
				Name:      "list_results_count",
				Help:      "Items returned per plugin list query",
				Buckets:   []float64{0, 1, 5, 10, 25, 50, 100, 250, 500},
			}),

			/* Cache */
			CacheInvalidationsTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "cache",
				Name:      "invalidations_total",
				Help:      "Cache invalidation operations by status",
			}, []string{"status"}),

			/* Execution Lifecycle */
			ExecutionCompletedTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "execution",
				Name:      "completed_total",
				Help:      "Completed executions by trigger source",
			}, []string{"trigger"}),

			ExecutionFailedTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "execution",
				Name:      "failed_total",
				Help:      "Failed executions by trigger source",
			}, []string{"trigger"}),

			ExecutionStartedTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "execution",
				Name:      "started_total",
				Help:      "Started executions by trigger source",
			}, []string{"trigger"}),

			ExecutionDuration: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "execution",
				Name:      "duration_seconds",
				Help:      "Total execution time from start to terminal state",
				Buckets:   []float64{0.01, 0.05, 0.1, 0.5, 1, 5, 10, 30, 60},
			}, []string{"trigger"}),

			ExecutionSteps: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "execution",
				Name:      "steps_total",
				Help:      "Steps executed per execution",
				Buckets:   []float64{1, 5, 10, 25, 50, 100, 200, 300},
			}),

			ExecutionActive: reg.NewGauge(metrics.GaugeOpts{
				Subsystem: "execution",
				Name:      "active",
				Help:      "Currently in-flight executions",
			}),

			/* Checkpoint */
			CheckpointDuration: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "checkpoint",
				Name:      "duration_seconds",
				Help:      "KV checkpoint (Put/Update) latency",
				Buckets:   []float64{0.001, 0.002, 0.005, 0.01, 0.025, 0.05, 0.1},
			}),

			/* Dispatch */
			DispatchTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "dispatch",
				Name:      "total",
				Help:      "Dispatch calls by type and result",
			}, []string{"type", "result"}),

			/* Resilience */
			CASRetriesTotal: reg.NewCounter(metrics.CounterOpts{
				Subsystem: "resilience",
				Name:      "cas_retries_total",
				Help:      "CAS conflict retry count in HandleResume",
			}),

			TokenRejectionsTotal: reg.NewCounter(metrics.CounterOpts{
				Subsystem: "resilience",
				Name:      "token_rejections_total",
				Help:      "Stale callback token rejections",
			}),

			ShutdownStopsTotal: reg.NewCounter(metrics.CounterOpts{
				Subsystem: "resilience",
				Name:      "shutdown_stops_total",
				Help:      "Walker stops due to graceful shutdown",
			}),
		}

		// Initialize label combinations so metrics appear in /metrics from startup
		for _, op := range []string{"create", "update", "delete", "get", "list"} {
			for _, s := range []string{"success", "error"} {
				m.DefinitionOperations.WithLabelValues(op, s)
				m.PluginOperations.WithLabelValues(op, s)
			}
			m.DefinitionOperationDuration.WithLabelValues(op)
			m.PluginOperationDuration.WithLabelValues(op)
		}
		for _, r := range []string{"hit", "miss"} {
			m.DefinitionCacheTotal.WithLabelValues(r)
		}
		for _, s := range []string{"success", "error"} {
			m.CacheInvalidationsTotal.WithLabelValues(s)
		}
		for _, t := range []string{"workflow", "subworkflow", "http"} {
			m.ExecutionCompletedTotal.WithLabelValues(t)
			m.ExecutionFailedTotal.WithLabelValues(t)
			m.ExecutionStartedTotal.WithLabelValues(t)
			m.ExecutionDuration.WithLabelValues(t)
		}
		for _, t := range []string{"subworkflow", "schedule", "resume"} {
			for _, r := range []string{"ok", "error"} {
				m.DispatchTotal.WithLabelValues(t, r)
			}
		}

		return m
	})

	logger.Info("[APP:BOOTSTRAP] Metrics registry initialized")
}
