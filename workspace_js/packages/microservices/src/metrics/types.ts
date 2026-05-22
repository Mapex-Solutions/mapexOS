/**
 * Metrics type definitions for prom-client wrappers.
 *
 * These interfaces mirror the Go metrics package pattern
 * (workspace_go/packages/microservices/metrics/) for consistency
 * across Go and Node.js services.
 */

/**
 * Options for creating a Counter metric.
 */
export interface CounterOpts {
	/** Optional subsystem name (middle part of metric name) */
	subsystem?: string;
	/** Metric name (last part of metric name) */
	name: string;
	/** Help text describing this metric */
	help: string;
}

/**
 * Options for creating a Histogram metric.
 */
export interface HistogramOpts {
	/** Optional subsystem name (middle part of metric name) */
	subsystem?: string;
	/** Metric name (last part of metric name) */
	name: string;
	/** Help text describing this metric */
	help: string;
	/** Custom bucket boundaries (defaults to prom-client defaults) */
	buckets?: number[];
}

/**
 * Options for creating a Gauge metric.
 */
export interface GaugeOpts {
	/** Optional subsystem name (middle part of metric name) */
	subsystem?: string;
	/** Metric name (last part of metric name) */
	name: string;
	/** Help text describing this metric */
	help: string;
}
