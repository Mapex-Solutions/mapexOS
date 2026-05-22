// loadgen is a custom HTTP load generator for MapexOS auth benchmarks.
//
// Unlike `hey`, it rotates through N different seeded users per request,
// providing realistic cache diversity and bcrypt path testing.
//
// Modes:
//   - login:    POST /auth/login with rotating email (client1..clientN@test.com)
//   - coverage: Phase 1 = login all users (collect JWTs), Phase 2 = GET with rotating Bearer tokens
//
// Output format is hey-compatible so existing parse_hey_output() works unchanged.
//
// Usage:
//
//	go run loadgen.go -mode login -url http://localhost:5000/api/v1/auth/login -n 10000 -c 50 -users 1000
//	go run loadgen.go -mode coverage -url http://localhost:5000/api/v1/auth/users/me/coverage \
//	  -login-url http://localhost:5000/api/v1/auth/login -n 10000 -c 50 -users 1000
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// result holds the outcome of a single HTTP request.
type result struct {
	statusCode int
	latency    time.Duration
	err        error
}

func main() {
	mode := flag.String("mode", "login", "Benchmark mode: login or coverage")
	url := flag.String("url", "", "Target URL for benchmark requests")
	loginURL := flag.String("login-url", "", "Login endpoint URL (coverage mode only; defaults to -url)")
	totalReqs := flag.Int("n", 10000, "Total number of requests")
	concurrency := flag.Int("c", 50, "Number of concurrent workers")
	userCount := flag.Int("users", 1000, "Number of seeded users (client1..clientN@test.com)")
	password := flag.String("password", "test@123", "Password for all seeded users")

	flag.Parse()

	if *url == "" {
		fmt.Fprintln(os.Stderr, "error: -url is required")
		flag.Usage()
		os.Exit(1)
	}
	if *mode != "login" && *mode != "coverage" {
		fmt.Fprintf(os.Stderr, "error: -mode must be 'login' or 'coverage', got '%s'\n", *mode)
		os.Exit(1)
	}
	if *loginURL == "" {
		*loginURL = *url
	}

	transport := &http.Transport{
		MaxIdleConns:        *concurrency + 10,
		MaxIdleConnsPerHost: *concurrency + 10,
		IdleConnTimeout:     90 * time.Second,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	// Phase 1 (coverage only): login all users to collect JWTs.
	var tokens []string
	if *mode == "coverage" {
		var err error
		tokens, err = loginAllUsers(client, *loginURL, *userCount, *password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: login phase failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "\n  Login phase complete: %d tokens collected\n\n", len(tokens))
	}

	// Phase 2: Benchmark.
	workChan := make(chan int, *concurrency*2)
	resultsChan := make(chan result, *concurrency*2)

	var wg sync.WaitGroup

	// Start workers.
	for w := 0; w < *concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range workChan {
				var r result
				switch *mode {
				case "login":
					r = doLogin(client, *url, i, *userCount, *password)
				case "coverage":
					r = doCoverage(client, *url, i, tokens)
				}
				resultsChan <- r
			}
		}()
	}

	// Feed work and collect results concurrently.
	results := make([]result, 0, *totalReqs)
	done := make(chan struct{})
	go func() {
		for r := range resultsChan {
			results = append(results, r)
		}
		close(done)
	}()

	benchStart := time.Now()

	// Feed indices.
	for i := 0; i < *totalReqs; i++ {
		workChan <- i
	}
	close(workChan)

	// Wait for all workers to finish.
	wg.Wait()
	close(resultsChan)

	// Wait for result collector.
	<-done

	benchDuration := time.Since(benchStart)

	printReport(results, benchDuration, *totalReqs)
}

// doLogin sends a POST login request for user at index i.
func doLogin(client *http.Client, url string, i, userCount int, password string) result {
	userIndex := (i % userCount) + 1
	body := fmt.Sprintf(`{"email":"client%d@test.com","password":"%s","keepConnected":false}`, userIndex, password)

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return result{err: err}
	}
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := client.Do(req)
	latency := time.Since(start)
	if err != nil {
		return result{latency: latency, err: err}
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	return result{statusCode: resp.StatusCode, latency: latency}
}

// doCoverage sends a GET request with a rotating Bearer token.
func doCoverage(client *http.Client, url string, i int, tokens []string) result {
	token := tokens[i%len(tokens)]

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return result{err: err}
	}
	req.Header.Set("Authorization", "Bearer "+token)

	start := time.Now()
	resp, err := client.Do(req)
	latency := time.Since(start)
	if err != nil {
		return result{latency: latency, err: err}
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	return result{statusCode: resp.StatusCode, latency: latency}
}

// loginResponse matches the MapexOS auth login response.
type loginResponse struct {
	Status int `json:"status"`
	Data   struct {
		AccessToken string `json:"access_token"`
	} `json:"data"`
}

// loginAllUsers sequentially logs in all seeded users and returns their JWT tokens.
func loginAllUsers(client *http.Client, loginURL string, userCount int, password string) ([]string, error) {
	tokens := make([]string, 0, userCount)
	fmt.Fprintf(os.Stderr, "  Logging in %d users to collect JWTs...\n", userCount)

	for i := 1; i <= userCount; i++ {
		body := fmt.Sprintf(`{"email":"client%d@test.com","password":"%s","keepConnected":false}`, i, password)

		req, err := http.NewRequest(http.MethodPost, loginURL, strings.NewReader(body))
		if err != nil {
			return nil, fmt.Errorf("user %d: request build failed: %w", i, err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("user %d: request failed: %w", i, err)
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("user %d: read body failed: %w", i, err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("user %d: status %d, body: %s", i, resp.StatusCode, string(respBody))
		}

		var lr loginResponse
		if err := json.Unmarshal(respBody, &lr); err != nil {
			return nil, fmt.Errorf("user %d: json decode failed: %w", i, err)
		}
		if lr.Data.AccessToken == "" {
			return nil, fmt.Errorf("user %d: empty access_token in response", i)
		}

		tokens = append(tokens, lr.Data.AccessToken)

		// Progress indicator every 100 users.
		if i%100 == 0 || i == userCount {
			fmt.Fprintf(os.Stderr, "    [%d/%d] tokens collected\n", i, userCount)
		}
	}

	return tokens, nil
}

// printReport outputs hey-compatible summary to stdout.
func printReport(results []result, totalDuration time.Duration, totalReqs int) {
	if len(results) == 0 {
		fmt.Println("No results collected.")
		return
	}

	// Separate successful latencies and count status codes.
	statusCounts := make(map[int]int)
	latencies := make([]float64, 0, len(results))
	var errCount int

	for _, r := range results {
		if r.err != nil {
			errCount++
			continue
		}
		statusCounts[r.statusCode]++
		latencies = append(latencies, r.latency.Seconds())
	}

	sort.Float64s(latencies)

	n := len(latencies)
	if n == 0 {
		fmt.Printf("All %d requests failed with errors.\n", errCount)
		return
	}

	// Compute stats.
	var sum float64
	for _, l := range latencies {
		sum += l
	}
	avg := sum / float64(n)
	fastest := latencies[0]
	slowest := latencies[n-1]
	rps := float64(totalReqs) / totalDuration.Seconds()

	// Percentile helper.
	pct := func(p float64) float64 {
		idx := int(float64(n-1) * p)
		if idx >= n {
			idx = n - 1
		}
		return latencies[idx]
	}

	// Print hey-compatible format.
	fmt.Println()
	fmt.Println("Summary:")
	fmt.Printf("  Total:\t%.4f secs\n", totalDuration.Seconds())
	fmt.Printf("  Slowest:\t%.4f secs\n", slowest)
	fmt.Printf("  Fastest:\t%.4f secs\n", fastest)
	fmt.Printf("  Average:\t%.4f secs\n", avg)
	fmt.Printf("  Requests/sec:\t%.4f\n", rps)
	fmt.Println()

	fmt.Println()
	fmt.Println("Latency distribution:")
	for _, p := range []struct {
		label string
		val   float64
	}{
		{"10%", 0.10},
		{"25%", 0.25},
		{"50%", 0.50},
		{"75%", 0.75},
		{"90%", 0.90},
		{"95%", 0.95},
		{"99%", 0.99},
	} {
		fmt.Printf("  %s in %.4f secs\n", p.label, pct(p.val))
	}
	fmt.Println()

	fmt.Println()
	fmt.Println("Status code distribution:")
	// Sort status codes for deterministic output.
	codes := make([]int, 0, len(statusCounts))
	for code := range statusCounts {
		codes = append(codes, code)
	}
	sort.Ints(codes)
	for _, code := range codes {
		fmt.Printf("  [%d]\t%d responses\n", code, statusCounts[code])
	}

	if errCount > 0 {
		fmt.Println()
		fmt.Printf("Error distribution:\n")
		fmt.Printf("  [connection errors]\t%d\n", errCount)
	}
	fmt.Println()
}
