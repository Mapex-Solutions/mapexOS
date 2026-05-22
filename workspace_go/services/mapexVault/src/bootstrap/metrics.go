package bootstrap

import (
	"github.com/prometheus/client_golang/prometheus"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/metrics"
	"go.uber.org/dig"
)

// VaultMetrics contains all metrics specific to the vault service.
type VaultMetrics struct {
	Registry *metrics.Registry

	// Credential CRUD operations by type and result
	CredentialOperations *prometheus.CounterVec
	// Credential operation latency
	CredentialOperationDuration *prometheus.HistogramVec
	// Decrypt requests (internal API)
	DecryptRequests *prometheus.CounterVec
	// Token refresh operations by result
	TokenRefreshTotal *prometheus.CounterVec
}

// InitMetrics creates the metrics registry and provides it to DIG.
func InitMetrics(c *dig.Container) {
	c.Provide(func() *VaultMetrics {
		reg := metrics.NewRegistry("vault")

		// Go runtime and process collectors — always enabled
		reg.EnableGoCollector()
		reg.EnableProcessCollector()

		m := &VaultMetrics{
			Registry: reg,

			CredentialOperations: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "credential",
				Name:      "operations_total",
				Help:      "Credential CRUD operations by type and result",
			}, []string{"operation", "status"}),

			CredentialOperationDuration: reg.NewHistogramVec(metrics.HistogramOpts{
				Subsystem: "credential",
				Name:      "operation_duration_seconds",
				Help:      "Credential operation latency",
				Buckets:   prometheus.DefBuckets,
			}, []string{"operation"}),

			DecryptRequests: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "credential",
				Name:      "decrypt_requests_total",
				Help:      "Internal API decrypt requests by result",
			}, []string{"status"}),

			TokenRefreshTotal: reg.NewCounterVec(metrics.CounterOpts{
				Subsystem: "credential",
				Name:      "token_refresh_total",
				Help:      "Token refresh operations by result",
			}, []string{"type", "status"}),
		}

		// Initialize label combinations so metrics appear in /metrics from startup
		for _, op := range []string{"create", "update", "delete", "get", "list", "test"} {
			for _, s := range []string{"success", "error"} {
				m.CredentialOperations.WithLabelValues(op, s)
			}
			m.CredentialOperationDuration.WithLabelValues(op)
		}
		for _, s := range []string{"success", "error"} {
			m.DecryptRequests.WithLabelValues(s)
		}
		for _, t := range []string{"oauth2", "userAndPass"} {
			for _, s := range []string{"success", "error"} {
				m.TokenRefreshTotal.WithLabelValues(t, s)
			}
		}

		return m
	})

	logger.Info("[APP:BOOTSTRAP] Metrics registry initialized")
}
