package main

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// BenchResult holds the aggregated benchmark results.
type BenchResult struct {
	TotalRequests int
	OkCount       int64
	FailCount     int64
	TimeoutCount  int64
	Duration      time.Duration
	Latencies     []time.Duration
	Scenario      string
	Concurrency   int
}

// credential holds pre-built MQTT credentials for a single device.
type credential struct {
	username string
	password string
}

// buildCredentialPool creates the rotating pool of MQTT device credentials.
// Username format: bench-mqtt-device-NNNNN (5-digit zero-padded)
// Password format: bench-mqtt-secret-NNNNN (5-digit zero-padded)
func buildCredentialPool(count int) []credential {
	pool := make([]credential, count)
	for i := 0; i < count; i++ {
		pool[i] = credential{
			username: fmt.Sprintf("bench-mqtt-device-%05d", i+1),
			password: fmt.Sprintf("bench-mqtt-secret-%05d", i+1),
		}
	}
	return pool
}

// RunBenchmark executes the auth callout benchmark using real MQTT connections.
//
// Both cache-hit and cache-miss scenarios pre-build a pool of credentials for
// UserCount unique MQTT users. Requests rotate through the pool.
//
// For cache-hit: Redis is primed with all users during warmup. All benchmark
// requests hit Redis cache.
// For cache-miss: Redis is flushed. First request per user hits MongoDB,
// subsequent requests for the same user hit Redis. Use user-count >= count
// for 100% DB hits.
//
// Flow:
//  1. Pre-build credential pool (one per unique user)
//  2. Warmup: sequential MQTT Connect/Disconnect (populates Redis cache)
//  3. Main loop: semaphore-limited goroutines, each does
//     MQTT Connect → measure CONNACK latency → Disconnect
//  4. Collect per-request latencies
//  5. Return aggregated results
func RunBenchmark(cfg Config) (*BenchResult, error) {
	// 1. Build credential pool
	creds := buildCredentialPool(cfg.UserCount)
	fmt.Printf("[1/4] Built %d credential pairs (user pool)\n", cfg.UserCount)

	// 2. Warmup — sequential connects to populate Redis cache
	if cfg.Warmup > 0 {
		fmt.Printf("[2/4] Warmup: %d sequential MQTT connections...\n", cfg.Warmup)
		for i := 0; i < cfg.Warmup; i++ {
			cred := creds[i%len(creds)]
			clientID := fmt.Sprintf("auth-bench-warmup-%d", i)

			opts := mqtt.NewClientOptions().
				AddBroker("tcp://" + cfg.MqttBroker).
				SetClientID(clientID).
				SetUsername(cred.username).
				SetPassword(cred.password).
				SetCleanSession(true).
				SetAutoReconnect(false).
				SetConnectTimeout(cfg.Timeout)

			c := mqtt.NewClient(opts)
			token := c.Connect()
			if token.WaitTimeout(cfg.Timeout) && token.Error() == nil {
				c.Disconnect(0)
			}
		}
		fmt.Println("      Warmup complete.")
	} else {
		fmt.Println("[2/4] Warmup skipped.")
	}

	// 3. Main benchmark loop
	fmt.Printf("[3/4] Running %d MQTT connections @ %d concurrency...\n", cfg.Count, cfg.Concurrency)

	latencies := make([]time.Duration, cfg.Count)
	var okCount, failCount, timeoutCount atomic.Int64

	sem := make(chan struct{}, cfg.Concurrency)
	var wg sync.WaitGroup

	start := time.Now()

	for i := 0; i < cfg.Count; i++ {
		sem <- struct{}{}
		wg.Add(1)

		go func(idx int) {
			defer func() {
				<-sem
				wg.Done()
			}()

			cred := creds[idx%len(creds)]
			clientID := fmt.Sprintf("auth-bench-%d-%d", idx%cfg.Concurrency, idx)

			opts := mqtt.NewClientOptions().
				AddBroker("tcp://" + cfg.MqttBroker).
				SetClientID(clientID).
				SetUsername(cred.username).
				SetPassword(cred.password).
				SetCleanSession(true).
				SetAutoReconnect(false).
				SetConnectTimeout(cfg.Timeout)

			c := mqtt.NewClient(opts)

			reqStart := time.Now()
			token := c.Connect()
			ok := token.WaitTimeout(cfg.Timeout)
			elapsed := time.Since(reqStart)

			latencies[idx] = elapsed

			if !ok {
				// WaitTimeout returned false → timed out
				timeoutCount.Add(1)
				return
			}

			if token.Error() != nil {
				failCount.Add(1)
				return
			}

			okCount.Add(1)
			c.Disconnect(0)
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	fmt.Println("[4/4] Benchmark complete.")

	return &BenchResult{
		TotalRequests: cfg.Count,
		OkCount:       okCount.Load(),
		FailCount:     failCount.Load(),
		TimeoutCount:  timeoutCount.Load(),
		Duration:      duration,
		Latencies:     latencies,
		Scenario:      cfg.Scenario,
		Concurrency:   cfg.Concurrency,
	}, nil
}

// Percentile computes the p-th percentile from a sorted latency slice.
func Percentile(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(float64(len(sorted)) * p / 100.0)
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}

// SortLatencies returns a sorted copy of the latency slice.
func SortLatencies(latencies []time.Duration) []time.Duration {
	sorted := make([]time.Duration, len(latencies))
	copy(sorted, latencies)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})
	return sorted
}
