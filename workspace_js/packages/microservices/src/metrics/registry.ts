/**
 * MetricsRegistry — prom-client wrapper following the Go metrics.Registry pattern.
 *
 * Provides a consistent API for creating Prometheus metrics with automatic
 * namespace prefixing: {namespace}_{subsystem}_{name}
 *
 * @example
 * ```typescript
 * const reg = new MetricsRegistry('jsexec');
 * reg.enableDefaultMetrics();
 *
 * const counter = reg.newCounterVec({ name: 'events_processed_total', help: '...' }, ['consumer', 'status']);
 * counter.inc({ consumer: 'mqtt_data', status: 'success' });
 *
 * reg.registerEndpoint(app); // GET /metrics
 * ```
 */
import {
	Registry,
	Counter,
	Histogram,
	Gauge,
	collectDefaultMetrics,
} from 'prom-client';

import type { RequestHandler, Application } from 'express';
import type { CounterOpts, HistogramOpts, GaugeOpts } from './types';

/**
 * MetricsRegistry wraps prom-client Registry with namespace support.
 *
 * Mirrors the Go MetricsRegistry API for cross-language consistency.
 */
export class MetricsRegistry {
	private readonly registry: Registry;
	private readonly namespace: string;

	/**
	 * Create a new MetricsRegistry.
	 *
	 * @param namespace - Metric namespace prefix (e.g., "jsexec", "httpgw")
	 */
	constructor(namespace: string) {
		this.registry = new Registry();
		this.namespace = namespace;
	}

	/**
	 * Build the full metric name: {namespace}_{subsystem}_{name} or {namespace}_{name}
	 *
	 * @param subsystem - Optional subsystem
	 * @param name - Metric name
	 * @returns Full metric name
	 */
	private buildName(subsystem: string | undefined, name: string): string {
		if (subsystem) {
			return `${this.namespace}_${subsystem}_${name}`;
		}
		return `${this.namespace}_${name}`;
	}

	/**
	 * Enable Node.js default metrics (heap, GC, event loop lag, etc.)
	 *
	 * @param prefix - Optional prefix for default metrics (defaults to empty)
	 */
	enableDefaultMetrics(prefix?: string): void {
		collectDefaultMetrics({
			register: this.registry,
			prefix: prefix || '',
		});
	}

	/**
	 * Create a Counter (no labels).
	 *
	 * @param opts - Counter options
	 * @returns prom-client Counter
	 */
	newCounter(opts: CounterOpts): Counter {
		const counter = new Counter({
			name: this.buildName(opts.subsystem, opts.name),
			help: opts.help,
			registers: [this.registry],
		});
		return counter;
	}

	/**
	 * Create a Counter with labels.
	 *
	 * @param opts - Counter options
	 * @param labels - Label names
	 * @returns prom-client Counter with labels
	 */
	newCounterVec(opts: CounterOpts, labels: string[]): Counter {
		const counter = new Counter({
			name: this.buildName(opts.subsystem, opts.name),
			help: opts.help,
			labelNames: labels,
			registers: [this.registry],
		});
		return counter;
	}

	/**
	 * Create a Histogram (no labels).
	 *
	 * @param opts - Histogram options
	 * @returns prom-client Histogram
	 */
	newHistogram(opts: HistogramOpts): Histogram {
		const histogram = new Histogram({
			name: this.buildName(opts.subsystem, opts.name),
			help: opts.help,
			buckets: opts.buckets,
			registers: [this.registry],
		});
		return histogram;
	}

	/**
	 * Create a Histogram with labels.
	 *
	 * @param opts - Histogram options
	 * @param labels - Label names
	 * @returns prom-client Histogram with labels
	 */
	newHistogramVec(opts: HistogramOpts, labels: string[]): Histogram {
		const histogram = new Histogram({
			name: this.buildName(opts.subsystem, opts.name),
			help: opts.help,
			buckets: opts.buckets,
			labelNames: labels,
			registers: [this.registry],
		});
		return histogram;
	}

	/**
	 * Create a Gauge (no labels).
	 *
	 * @param opts - Gauge options
	 * @returns prom-client Gauge
	 */
	newGauge(opts: GaugeOpts): Gauge {
		const gauge = new Gauge({
			name: this.buildName(opts.subsystem, opts.name),
			help: opts.help,
			registers: [this.registry],
		});
		return gauge;
	}

	/**
	 * Create a Gauge with labels.
	 *
	 * @param opts - Gauge options
	 * @param labels - Label names
	 * @returns prom-client Gauge with labels
	 */
	newGaugeVec(opts: GaugeOpts, labels: string[]): Gauge {
		const gauge = new Gauge({
			name: this.buildName(opts.subsystem, opts.name),
			help: opts.help,
			labelNames: labels,
			registers: [this.registry],
		});
		return gauge;
	}

	/**
	 * Returns an Express RequestHandler that serves metrics in Prometheus text format.
	 *
	 * @returns Express middleware for GET /metrics
	 */
	handler(): RequestHandler {
		return async (_req, res) => {
			try {
				res.set('Content-Type', this.registry.contentType);
				const metrics = await this.registry.metrics();
				res.end(metrics);
			} catch (error) {
				res.status(500).end(String(error));
			}
		};
	}

	/**
	 * Convenience method to register GET /metrics on an Express app.
	 *
	 * @param app - Express application
	 */
	registerEndpoint(app: Application): void {
		app.get('/metrics', this.handler());
	}

	/**
	 * Get the underlying prom-client Registry (for advanced usage).
	 *
	 * @returns The prom-client Registry instance
	 */
	getRegistry(): Registry {
		return this.registry;
	}
}
