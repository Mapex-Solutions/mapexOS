package bootstrap

import (
	"github.com/prometheus/client_golang/prometheus"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/metrics"
	"go.uber.org/dig"
)

/**
 * AssetsMetrics contains all metrics specific to the assets service.
 * This struct is provided to the DIG container so any module can inject it.
 */
type AssetsMetrics struct {
	Registry *metrics.Registry

	/* Asset CRUD */

	// Asset CRUD operations by type and result
	AssetOperations *prometheus.CounterVec
	// Asset CRUD operation latency
	AssetOperationDuration *prometheus.HistogramVec
	// Items returned per asset list query
	AssetListResultsCount prometheus.Histogram
	// Asset counter cache hit/miss
	AssetCacheTotal *prometheus.CounterVec

	/* Asset Template CRUD */

	// Template CRUD operations by type and result
	TemplateOperations *prometheus.CounterVec
	// Template CRUD operation latency
	TemplateOperationDuration *prometheus.HistogramVec
	// Items returned per template list query
	TemplateListResultsCount prometheus.Histogram
	// Template counter cache hit/miss
	TemplateCacheTotal *prometheus.CounterVec

	/* Auth Callout */

	// Auth callout attempts by result (success/rejected/error)
	AuthCalloutTotal *prometheus.CounterVec
	// Auth callout processing latency
	AuthCalloutDuration prometheus.Histogram
	// Auth cache hit/miss counter
	AuthCacheTotal *prometheus.CounterVec

	/* Health Monitor */

	// Heartbeats received from ASSET-HEARTBEAT stream
	HealthHeartbeatsReceived prometheus.Counter
	// Scan cycle duration
	HealthScanDuration prometheus.Histogram
	// Sensors evaluated during scan
	HealthSensorsScanned prometheus.Counter
	// Sensors confirmed stale (alert fired)
	HealthSensorsStale prometheus.Counter
	// Alerts published by type (offline/online)
	HealthAlertsPublished *prometheus.CounterVec
	// Miss counter increments
	HealthMissIncrements prometheus.Counter
	// Online transitions (offline → online)
	HealthOnlineTransitions prometheus.Counter
	// Redis operation errors by operation
	HealthRedisErrors *prometheus.CounterVec
	// Active orgs being monitored
	HealthOrgsMonitored prometheus.Gauge

	/* Health Monitor — MQTT Presence */

	// Presence events received by action (connect_from_auth | disconnect_from_sys)
	HealthPresenceReceived *prometheus.CounterVec
	// Presence events filtered before any I/O, by reason
	// (non_mqtt_client | non_client_kind | empty_user | user_unmapped |
	//  asset_lookup_error | asset_not_found | asset_disabled |
	//  never_connected | stale_disconnect)
	HealthPresenceFiltered *prometheus.CounterVec
	// Presence events that mutated state, by outcome (success | error)
	HealthPresenceProcessed *prometheus.CounterVec
	// End-to-end presence handler duration in seconds
	HealthPresenceHandlerDuration prometheus.Histogram
}

// InitMetrics creates the metrics registry, declares all service-specific
// metrics, and provides the metrics struct to the DIG container.
func InitMetrics(c *dig.Container) {
	c.Provide(func() *AssetsMetrics {
		reg := metrics.NewRegistry("assets")

		// Go runtime and process collectors — always enabled
		reg.EnableGoCollector()
		reg.EnableProcessCollector()

		m := &AssetsMetrics{
			Registry: reg,

			/* Asset CRUD */
			AssetOperations: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "asset",
				Name:      "operations_total",
				Help:      "Asset CRUD operations by type and result",
			}, []string{"operation", "status"}),

			AssetOperationDuration: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "asset",
				Name:      "operation_duration_seconds",
				Help:      "Asset CRUD operation latency",
				Buckets:   prometheus.DefBuckets,
			}, []string{"operation"}),

			AssetListResultsCount: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "asset",
				Name:      "list_results_count",
				Help:      "Items returned per asset list query",
				Buckets:   []float64{0, 1, 5, 10, 25, 50, 100, 250},
			}),

			AssetCacheTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "asset",
				Name:      "cache_total",
				Help:      "Asset counter cache lookups by result (hit/miss)",
			}, []string{"result"}),

			/* Asset Template CRUD */
			TemplateOperations: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "template",
				Name:      "operations_total",
				Help:      "Template CRUD operations by type and result",
			}, []string{"operation", "status"}),

			TemplateOperationDuration: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "template",
				Name:      "operation_duration_seconds",
				Help:      "Template CRUD operation latency",
				Buckets:   prometheus.DefBuckets,
			}, []string{"operation"}),

			TemplateListResultsCount: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "template",
				Name:      "list_results_count",
				Help:      "Items returned per template list query",
				Buckets:   []float64{0, 1, 5, 10, 25, 50, 100, 250},
			}),

			TemplateCacheTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "template",
				Name:      "cache_total",
				Help:      "Template counter cache lookups by result (hit/miss)",
			}, []string{"result"}),

			/* Auth Callout */
			AuthCalloutTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "auth",
				Name:      "callout_total",
				Help:      "Auth callout attempts by result (success/rejected/error)",
			}, []string{"result"}),

			AuthCalloutDuration: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "auth",
				Name:      "callout_duration_seconds",
				Help:      "Auth callout processing latency",
				Buckets:   prometheus.DefBuckets,
			}),

			AuthCacheTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "auth",
				Name:      "cache_total",
				Help:      "Auth cache lookups by result (hit/miss)",
			}, []string{"result"}),

			/* Health Monitor */
			HealthHeartbeatsReceived: reg.NewCounter(metrics.CounterOpts{
				Subsystem: "health",
				Name:      "heartbeats_received_total",
				Help:      "Heartbeats received from ASSET-HEARTBEAT stream",
			}),

			HealthScanDuration: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "health",
				Name:      "scan_duration_seconds",
				Help:      "Health monitor scan cycle duration",
				Buckets:   []float64{0.1, 0.5, 1, 5, 10, 30, 60},
			}),

			HealthSensorsScanned: reg.NewCounter(metrics.CounterOpts{
				Subsystem: "health",
				Name:      "sensors_scanned_total",
				Help:      "Total sensors evaluated during scan cycles",
			}),

			HealthSensorsStale: reg.NewCounter(metrics.CounterOpts{
				Subsystem: "health",
				Name:      "sensors_stale_total",
				Help:      "Sensors confirmed stale (alert fired)",
			}),

			HealthAlertsPublished: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "health",
				Name:      "alerts_published_total",
				Help:      "Health alerts published by type (offline/online)",
			}, []string{"type"}),

			HealthMissIncrements: reg.NewCounter(metrics.CounterOpts{
				Subsystem: "health",
				Name:      "miss_increments_total",
				Help:      "Miss counter increments during scan",
			}),

			HealthOnlineTransitions: reg.NewCounter(metrics.CounterOpts{
				Subsystem: "health",
				Name:      "online_transitions_total",
				Help:      "Online transitions (offline → online)",
			}),

			HealthRedisErrors: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "health",
				Name:      "redis_errors_total",
				Help:      "Redis operation errors by operation type",
			}, []string{"operation"}),

			HealthOrgsMonitored: reg.NewGauge(metrics.GaugeOpts{
				Subsystem: "health",
				Name:      "orgs_monitored",
				Help:      "Number of orgs with active health monitoring",
			}),

			HealthPresenceReceived: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "health",
				Name:      "presence_received_total",
				Help:      "MQTT presence events received by action",
			}, []string{"action"}),

			HealthPresenceFiltered: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "health",
				Name:      "presence_filtered_total",
				Help:      "MQTT presence events filtered before any I/O, by reason",
			}, []string{"reason"}),

			HealthPresenceProcessed: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "health",
				Name:      "presence_processed_total",
				Help:      "MQTT presence events that mutated state, by outcome",
			}, []string{"outcome"}),

			HealthPresenceHandlerDuration: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "health",
				Name:      "presence_handler_duration_seconds",
				Help:      "End-to-end MQTT presence handler duration",
				Buckets:   prometheus.DefBuckets,
			}),
		}

		// Initialize label combinations so metrics appear in /metrics from startup
		for _, op := range []string{"create", "update", "delete", "get", "list"} {
			m.AssetOperationDuration.WithLabelValues(op)
			m.TemplateOperationDuration.WithLabelValues(op)
			for _, s := range []string{"success", "error"} {
				m.AssetOperations.WithLabelValues(op, s)
				m.TemplateOperations.WithLabelValues(op, s)
			}
		}
		for _, r := range []string{"hit", "miss"} {
			m.AssetCacheTotal.WithLabelValues(r)
			m.TemplateCacheTotal.WithLabelValues(r)
			m.AuthCacheTotal.WithLabelValues(r)
		}
		for _, r := range []string{"success", "error"} {
			m.AuthCalloutTotal.WithLabelValues(r)
		}
		for _, t := range []string{"offline", "online"} {
			m.HealthAlertsPublished.WithLabelValues(t)
		}
		for _, op := range []string{"zadd", "zrangebyscore", "hincrby", "sismember", "smismember", "zmscore", "sadd"} {
			m.HealthRedisErrors.WithLabelValues(op)
		}
		for _, a := range []string{"connect_from_auth", "disconnect_from_sys"} {
			m.HealthPresenceReceived.WithLabelValues(a)
		}
		for _, r := range []string{
			"non_mqtt_client", "non_client_kind", "empty_user", "user_unmapped",
			"asset_lookup_error", "asset_not_found", "asset_disabled",
			"never_connected", "stale_disconnect",
		} {
			m.HealthPresenceFiltered.WithLabelValues(r)
		}
		for _, o := range []string{"success", "error"} {
			m.HealthPresenceProcessed.WithLabelValues(o)
		}

		return m
	})

	logger.Info("[INFRA:Metrics] Registry initialized")
}
