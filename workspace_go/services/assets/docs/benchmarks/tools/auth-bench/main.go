// Package main provides a standalone benchmark tool for NATS Auth Callout
// via real MQTT connections (port 1883).
//
// The benchmark measures end-to-end authentication latency: MQTT CONNECT →
// NATS auth callout → Assets service → CONNACK. This is the real device
// authentication flow.
//
// Usage:
//
//	./auth-bench --count 10000 --concurrency 50 --scenario cache-hit
//	./auth-bench --count 10000 --scenario cache-miss --mqtt-broker localhost:1883
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

// Config holds all CLI configuration for the benchmark.
type Config struct {
	Count       int
	Concurrency int
	Scenario    string
	MqttBroker  string
	Warmup      int
	Timeout     time.Duration
	UserCount   int
}

func main() {
	cfg := Config{}

	flag.IntVar(&cfg.Count, "count", 10000, "Number of auth requests")
	flag.IntVar(&cfg.Concurrency, "concurrency", 50, "Parallel workers")
	flag.StringVar(&cfg.Scenario, "scenario", "cache-hit", "Scenario: cache-hit | cache-miss")
	flag.StringVar(&cfg.MqttBroker, "mqtt-broker", "localhost:1883", "MQTT broker address (host:port)")
	flag.IntVar(&cfg.Warmup, "warmup", 100, "Warmup requests (discarded)")
	flag.DurationVar(&cfg.Timeout, "timeout", 5*time.Second, "Per-request timeout")
	flag.IntVar(&cfg.UserCount, "user-count", 10000, "Number of seeded MQTT users to rotate through")

	flag.Parse()

	// Validate scenario
	switch cfg.Scenario {
	case "cache-hit", "cache-miss":
		// ok
	default:
		fmt.Fprintf(os.Stderr, "ERROR: invalid scenario %q (use cache-hit or cache-miss)\n", cfg.Scenario)
		os.Exit(1)
	}

	// For cache-miss, ensure user-count >= count for 100% DB hits.
	// If user-count < count, requests will wrap around and some will be cache hits.
	if cfg.Scenario == "cache-miss" && cfg.UserCount < cfg.Count {
		fmt.Fprintf(os.Stderr, "WARN: --user-count (%d) < --count (%d) — after first %d requests, subsequent ones will hit Redis cache\n",
			cfg.UserCount, cfg.Count, cfg.UserCount)
	}

	// Print config
	fmt.Println("================================================================")
	fmt.Println("  Auth Callout Benchmark (MQTT Connect)")
	fmt.Println("================================================================")
	fmt.Printf("  Scenario:     %s\n", cfg.Scenario)
	fmt.Printf("  Count:        %d\n", cfg.Count)
	fmt.Printf("  Concurrency:  %d\n", cfg.Concurrency)
	fmt.Printf("  Warmup:       %d\n", cfg.Warmup)
	fmt.Printf("  MQTT Broker:  %s\n", cfg.MqttBroker)
	fmt.Printf("  Timeout:      %s\n", cfg.Timeout)
	fmt.Printf("  User pool:    %d users (rotating)\n", cfg.UserCount)
	fmt.Println("================================================================")
	fmt.Println()

	// Run benchmark
	result, err := RunBenchmark(cfg)
	if err != nil {
		log.Fatalf("Benchmark failed: %v", err)
	}

	// Print results
	PrintResults(result)
}
