package bootstrap

import (
	"github.com/prometheus/client_golang/prometheus"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/metrics"
	"go.uber.org/dig"
)

// MapexosMetrics contains all metrics specific to the mapexos service.
// This struct is provided to the DIG container so any module can inject it.
type MapexosMetrics struct {
	Registry *metrics.Registry

	// --- Auth (P0) ---

	// Authentication attempts by method (login/refresh/logout) and status (success/failure)
	AuthAttemptsTotal *prometheus.CounterVec
	// Authentication operation latency by method
	AuthDurationSeconds *prometheus.HistogramVec
	// Session operations by operation (store/get/invalidate) and status (success/failure)
	SessionOperationsTotal *prometheus.CounterVec

	// --- User CRUD (P0) ---

	// User CRUD operations by operation (create/get/get_by_id/update/delete/list/count) and status (success/failure)
	UserOperationsTotal *prometheus.CounterVec
	// User CRUD operation latency by operation
	UserOperationDurationSeconds *prometheus.HistogramVec
	// Items returned per user list query
	UserListResultsCount prometheus.Histogram

	// --- Group CRUD (P1) ---

	// Group CRUD operations by operation (create/get/get_by_id/update/delete/list/add_member/remove_member) and status
	GroupOperationsTotal *prometheus.CounterVec
	// Group CRUD operation latency by operation
	GroupOperationDurationSeconds *prometheus.HistogramVec
	// Items returned per group list query
	GroupListResultsCount prometheus.Histogram

	// --- Role CRUD (P1) ---

	// Role CRUD operations by operation (create/get/get_by_id/update/delete/list) and status
	RoleOperationsTotal *prometheus.CounterVec
	// Role CRUD operation latency by operation
	RoleOperationDurationSeconds *prometheus.HistogramVec
	// Items returned per role list query
	RoleListResultsCount prometheus.Histogram

	// --- Membership CRUD (P1) ---

	// Membership CRUD operations by operation (create/get/get_by_id/update/delete/list) and status
	MembershipOperationsTotal *prometheus.CounterVec
	// Membership CRUD operation latency by operation
	MembershipOperationDurationSeconds *prometheus.HistogramVec
	// Items returned per membership list query
	MembershipListResultsCount prometheus.Histogram

	// --- Organization CRUD (P1) ---

	// Organization CRUD operations by operation and status
	OrganizationOperationsTotal *prometheus.CounterVec
	// Organization CRUD operation latency by operation
	OrganizationOperationDurationSeconds *prometheus.HistogramVec
	// Items returned per organization list query
	OrganizationListResultsCount prometheus.Histogram

	// --- Cache (P1) ---

	// Cache lookups by type (authorization/coverage/counter) and result (hit/miss)
	CacheTotal *prometheus.CounterVec
}

// InitMetrics creates the metrics registry, declares all service-specific
// metrics, and provides the metrics struct to the DIG container.
func InitMetrics(c *dig.Container) {
	c.Provide(func() *MapexosMetrics {
		reg := metrics.NewRegistry("mapexos")

		// Go runtime and process collectors — always enabled
		reg.EnableGoCollector()
		reg.EnableProcessCollector()

		m := &MapexosMetrics{
			Registry: reg,

			// --- Auth ---
			AuthAttemptsTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "auth",
				Name:      "attempts_total",
				Help:      "Authentication attempts by method and status",
			}, []string{"method", "status"}),

			AuthDurationSeconds: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "auth",
				Name:      "duration_seconds",
				Help:      "Authentication operation latency by method",
				Buckets:   prometheus.DefBuckets,
			}, []string{"method"}),

			SessionOperationsTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "session",
				Name:      "operations_total",
				Help:      "Session operations by operation and status",
			}, []string{"operation", "status"}),

			// --- User CRUD ---
			UserOperationsTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "user",
				Name:      "operations_total",
				Help:      "User CRUD operations by operation and status",
			}, []string{"operation", "status"}),

			UserOperationDurationSeconds: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "user",
				Name:      "operation_duration_seconds",
				Help:      "User CRUD operation latency by operation",
				Buckets:   prometheus.DefBuckets,
			}, []string{"operation"}),

			UserListResultsCount: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "user",
				Name:      "list_results_count",
				Help:      "Items returned per user list query",
				Buckets:   []float64{0, 1, 5, 10, 25, 50, 100, 250},
			}),

			// --- Group CRUD ---
			GroupOperationsTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "group",
				Name:      "operations_total",
				Help:      "Group CRUD operations by operation and status",
			}, []string{"operation", "status"}),

			GroupOperationDurationSeconds: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "group",
				Name:      "operation_duration_seconds",
				Help:      "Group CRUD operation latency by operation",
				Buckets:   prometheus.DefBuckets,
			}, []string{"operation"}),

			GroupListResultsCount: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "group",
				Name:      "list_results_count",
				Help:      "Items returned per group list query",
				Buckets:   []float64{0, 1, 5, 10, 25, 50, 100, 250},
			}),

			// --- Role CRUD ---
			RoleOperationsTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "role",
				Name:      "operations_total",
				Help:      "Role CRUD operations by operation and status",
			}, []string{"operation", "status"}),

			RoleOperationDurationSeconds: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "role",
				Name:      "operation_duration_seconds",
				Help:      "Role CRUD operation latency by operation",
				Buckets:   prometheus.DefBuckets,
			}, []string{"operation"}),

			RoleListResultsCount: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "role",
				Name:      "list_results_count",
				Help:      "Items returned per role list query",
				Buckets:   []float64{0, 1, 5, 10, 25, 50, 100, 250},
			}),

			// --- Membership CRUD ---
			MembershipOperationsTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "membership",
				Name:      "operations_total",
				Help:      "Membership CRUD operations by operation and status",
			}, []string{"operation", "status"}),

			MembershipOperationDurationSeconds: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "membership",
				Name:      "operation_duration_seconds",
				Help:      "Membership CRUD operation latency by operation",
				Buckets:   prometheus.DefBuckets,
			}, []string{"operation"}),

			MembershipListResultsCount: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "membership",
				Name:      "list_results_count",
				Help:      "Items returned per membership list query",
				Buckets:   []float64{0, 1, 5, 10, 25, 50, 100, 250},
			}),

			// --- Organization CRUD ---
			OrganizationOperationsTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "organization",
				Name:      "operations_total",
				Help:      "Organization CRUD operations by operation and status",
			}, []string{"operation", "status"}),

			OrganizationOperationDurationSeconds: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "organization",
				Name:      "operation_duration_seconds",
				Help:      "Organization CRUD operation latency by operation",
				Buckets:   prometheus.DefBuckets,
			}, []string{"operation"}),

			OrganizationListResultsCount: reg.NewHistogram(metrics.HistogramOpts{
				Subsystem: "organization",
				Name:      "list_results_count",
				Help:      "Items returned per organization list query",
				Buckets:   []float64{0, 1, 5, 10, 25, 50, 100, 250},
			}),

			// --- Cache ---
			CacheTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "cache",
				Name:      "total",
				Help:      "Cache lookups by type and result",
			}, []string{"type", "result"}),
		}

		// Initialize label combinations so metrics appear in /metrics from startup
		for _, method := range []string{"password", "token", "apikey"} {
			m.AuthDurationSeconds.WithLabelValues(method)
			for _, s := range []string{"success", "error"} {
				m.AuthAttemptsTotal.WithLabelValues(method, s)
			}
		}
		for _, op := range []string{"create", "delete", "refresh"} {
			for _, s := range []string{"success", "error"} {
				m.SessionOperationsTotal.WithLabelValues(op, s)
			}
		}
		crudOps := []string{"create", "update", "delete", "get", "list"}
		for _, op := range crudOps {
			m.UserOperationDurationSeconds.WithLabelValues(op)
			m.GroupOperationDurationSeconds.WithLabelValues(op)
			m.RoleOperationDurationSeconds.WithLabelValues(op)
			m.MembershipOperationDurationSeconds.WithLabelValues(op)
			m.OrganizationOperationDurationSeconds.WithLabelValues(op)
			for _, s := range []string{"success", "error"} {
				m.UserOperationsTotal.WithLabelValues(op, s)
				m.GroupOperationsTotal.WithLabelValues(op, s)
				m.RoleOperationsTotal.WithLabelValues(op, s)
				m.MembershipOperationsTotal.WithLabelValues(op, s)
				m.OrganizationOperationsTotal.WithLabelValues(op, s)
			}
		}
		for _, t := range []string{"authorization", "coverage", "counter"} {
			for _, r := range []string{"hit", "miss"} {
				m.CacheTotal.WithLabelValues(t, r)
			}
		}

		return m
	})

	logger.Info("[APP:BOOTSTRAP] Metrics registry initialized")
}
