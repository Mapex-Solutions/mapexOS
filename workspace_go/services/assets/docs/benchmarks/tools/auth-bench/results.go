package main

import (
	"fmt"
	"time"
)

// PrintResults outputs both machine-parseable KEY=VALUE lines and a
// human-readable summary box.
//
// The KEY=VALUE lines are designed for shell scripts to parse via grep:
//
//	AUTH_BENCH_RPS=$(./auth-bench ... | grep AUTH_BENCH_RPS | cut -d= -f2)
func PrintResults(r *BenchResult) {
	sorted := SortLatencies(r.Latencies)

	rps := float64(r.TotalRequests) / r.Duration.Seconds()
	avgMs := avgLatencyMs(sorted)
	p50Ms := latencyMs(Percentile(sorted, 50))
	p95Ms := latencyMs(Percentile(sorted, 95))
	p99Ms := latencyMs(Percentile(sorted, 99))

	// ── Machine-parseable output ─────────────────────────────
	fmt.Println()
	fmt.Println("# --- AUTH BENCH RESULTS (machine-parseable) ---")
	fmt.Printf("AUTH_BENCH_RPS=%.0f\n", rps)
	fmt.Printf("AUTH_BENCH_AVG_MS=%.2f\n", avgMs)
	fmt.Printf("AUTH_BENCH_P50_MS=%.2f\n", p50Ms)
	fmt.Printf("AUTH_BENCH_P95_MS=%.2f\n", p95Ms)
	fmt.Printf("AUTH_BENCH_P99_MS=%.2f\n", p99Ms)
	fmt.Printf("AUTH_BENCH_OK=%d\n", r.OkCount)
	fmt.Printf("AUTH_BENCH_FAIL=%d\n", r.FailCount)
	fmt.Printf("AUTH_BENCH_TIMEOUT=%d\n", r.TimeoutCount)
	fmt.Printf("AUTH_BENCH_DURATION_S=%.1f\n", r.Duration.Seconds())

	// ── Human-readable summary ───────────────────────────────
	fmt.Println()
	fmt.Println("  ┌──────────────────────────────────────────────────────────┐")
	fmt.Printf("  │ Auth Callout Benchmark — %s\n", r.Scenario)
	fmt.Println("  ├──────────────────────────────────────────────────────────┤")
	fmt.Printf("  │ Throughput:  %.0f req/s\n", rps)
	fmt.Printf("  │ Duration:   %.1fs\n", r.Duration.Seconds())
	fmt.Printf("  │ Concurrency: %d\n", r.Concurrency)
	fmt.Println("  ├──────────────────────────────────────────────────────────┤")
	fmt.Printf("  │ Latency avg: %.2f ms\n", avgMs)
	fmt.Printf("  │ Latency p50: %.2f ms\n", p50Ms)
	fmt.Printf("  │ Latency p95: %.2f ms\n", p95Ms)
	fmt.Printf("  │ Latency p99: %.2f ms\n", p99Ms)
	fmt.Println("  ├──────────────────────────────────────────────────────────┤")
	fmt.Printf("  │ OK:       %d\n", r.OkCount)
	fmt.Printf("  │ Failed:   %d\n", r.FailCount)
	fmt.Printf("  │ Timeout:  %d\n", r.TimeoutCount)
	fmt.Printf("  │ Total:    %d\n", r.TotalRequests)
	fmt.Println("  └──────────────────────────────────────────────────────────┘")
	fmt.Println()
}

// latencyMs converts a duration to milliseconds as float64.
func latencyMs(d time.Duration) float64 {
	return float64(d.Microseconds()) / 1000.0
}

// avgLatencyMs computes the average latency in milliseconds.
func avgLatencyMs(sorted []time.Duration) float64 {
	if len(sorted) == 0 {
		return 0
	}
	var total time.Duration
	for _, d := range sorted {
		total += d
	}
	return float64(total.Microseconds()) / float64(len(sorted)) / 1000.0
}
