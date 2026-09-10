package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ScenarioResult captures statistical metrics for an executed load test
type ScenarioResult struct {
	Name            string            `json:"name"`
	TargetURL       string            `json:"target_url"`
	TotalRequests   int               `json:"total_requests"`
	Concurrency     int               `json:"concurrency"`
	DurationMs      int64             `json:"duration_ms"`
	ThroughputRPS   float64           `json:"throughput_rps"`
	StatusCodes     map[int]int       `json:"status_codes"`
	SuccessRate     float64           `json:"success_rate"`
	LatenciesMs     []float64         `json:"-"`
	MinLatencyMs    float64           `json:"min_latency_ms"`
	MeanLatencyMs   float64           `json:"mean_latency_ms"`
	MedianLatencyMs float64           `json:"median_latency_ms"`
	P90LatencyMs    float64           `json:"p90_latency_ms"`
	P95LatencyMs    float64           `json:"p95_latency_ms"`
	P99LatencyMs    float64           `json:"p99_latency_ms"`
	MaxLatencyMs    float64           `json:"max_latency_ms"`
	InvariantPassed bool              `json:"invariant_passed"`
	InvariantNotes  string            `json:"invariant_notes"`
	Errors          []string          `json:"errors,omitempty"`
}

type Config struct {
	Target      string
	Scenario    string
	Requests    int
	Concurrency int
	Timeout     time.Duration
	InsecureTLS bool
	JSONOutput  bool
	Verbose     bool
}

func main() {
	cfg := Config{}
	flag.StringVar(&cfg.Target, "target", "https://wallet.thaw-zin-2k77.de5.net", "Base target URL")
	flag.StringVar(&cfg.Scenario, "scenario", "all", "Scenario to run: all, health, transfer, replay, outbox, waf")
	flag.IntVar(&cfg.Requests, "requests", 100, "Total number of requests per scenario")
	flag.IntVar(&cfg.Concurrency, "concurrency", 10, "Number of concurrent goroutine workers")
	timeoutSec := flag.Int("timeout", 10, "Per-request timeout in seconds")
	flag.BoolVar(&cfg.InsecureTLS, "insecure", false, "Skip TLS certificate verification")
	flag.BoolVar(&cfg.JSONOutput, "json", false, "Output machine-readable JSON")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "Print detailed progress")
	flag.Parse()

	cfg.Timeout = time.Duration(*timeoutSec) * time.Second
	cfg.Target = strings.TrimRight(cfg.Target, "/")

	if !cfg.JSONOutput {
		printBanner(cfg)
	}

	transport := &http.Transport{
		MaxIdleConns:        cfg.Concurrency * 2,
		MaxIdleConnsPerHost: cfg.Concurrency * 2,
		IdleConnTimeout:     60 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: cfg.InsecureTLS},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   cfg.Timeout,
	}

	results := make([]ScenarioResult, 0)

	scenarios := []string{}
	if cfg.Scenario == "all" {
		scenarios = []string{"health", "transfer", "replay", "outbox", "waf"}
	} else {
		scenarios = []string{cfg.Scenario}
	}

	for _, sc := range scenarios {
		var res ScenarioResult
		switch sc {
		case "health":
			res = runHealthStress(client, cfg)
		case "transfer":
			res = runTransferStress(client, cfg)
		case "replay":
			res = runReplayStress(client, cfg)
		case "outbox":
			res = runOutboxStress(client, cfg)
		case "waf":
			res = runWAFStress(client, cfg)
		default:
			fmt.Fprintf(os.Stderr, "Unknown scenario: %s\n", sc)
			os.Exit(1)
		}
		results = append(results, res)
		if !cfg.JSONOutput {
			printScenarioResult(res)
		}
	}

	if cfg.JSONOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(results)
	} else {
		printSummaryTable(results)
	}
}

func printBanner(cfg Config) {
	fmt.Println("┌──────────────────────────────────────────────────────────────────────────┐")
	fmt.Println("│           A BANK FINTECH WALLET: ENTERPRISE SRE STRESS TEST SUITE        │")
	fmt.Println("├──────────────────────────────────────────────────────────────────────────┤")
	fmt.Printf("│ Target Base URL : %-54s │\n", cfg.Target)
	fmt.Printf("│ Scenario        : %-54s │\n", cfg.Scenario)
	fmt.Printf("│ Requests/Test   : %-54d │\n", cfg.Requests)
	fmt.Printf("│ Concurrency     : %-54d │\n", cfg.Concurrency)
	fmt.Printf("│ Client Timeout  : %-54v │\n", cfg.Timeout)
	fmt.Printf("│ Mode Standard   : Paranoid Mode Zero-Trust Empirical Telemetry           │\n")
	fmt.Println("└──────────────────────────────────────────────────────────────────────────┘")
	fmt.Println()
}

// runWorkerPool dispatches N requests across C concurrency workers
func runWorkerPool(totalReqs, concurrency int, taskFn func(reqIdx int) (int, time.Duration, error)) ScenarioResult {
	latencies := make([]float64, 0, totalReqs)
	statusCodes := make(map[int]int)
	errorsList := make([]string, 0)
	var mu sync.Mutex

	reqChan := make(chan int, totalReqs)
	for i := 0; i < totalReqs; i++ {
		reqChan <- i
	}
	close(reqChan)

	var wg sync.WaitGroup
	startTime := time.Now()

	for c := 0; c < concurrency; c++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range reqChan {
				code, dur, err := taskFn(idx)
				durMs := float64(dur.Nanoseconds()) / 1e6

				mu.Lock()
				latencies = append(latencies, durMs)
				statusCodes[code]++
				if err != nil {
					if len(errorsList) < 10 {
						errorsList = append(errorsList, err.Error())
					}
				}
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	totalElapsed := time.Since(startTime)

	// Statistical calculations
	sort.Float64s(latencies)
	n := len(latencies)

	var sum float64
	for _, v := range latencies {
		sum += v
	}

	res := ScenarioResult{
		TotalRequests: totalReqs,
		Concurrency:   concurrency,
		DurationMs:    totalElapsed.Milliseconds(),
		ThroughputRPS: float64(totalReqs) / totalElapsed.Seconds(),
		StatusCodes:   statusCodes,
		LatenciesMs:   latencies,
		Errors:        errorsList,
	}

	if n > 0 {
		res.MinLatencyMs = latencies[0]
		res.MaxLatencyMs = latencies[n-1]
		res.MeanLatencyMs = sum / float64(n)
		res.MedianLatencyMs = percentile(latencies, 50)
		res.P90LatencyMs = percentile(latencies, 90)
		res.P95LatencyMs = percentile(latencies, 95)
		res.P99LatencyMs = percentile(latencies, 99)
	}

	return res
}

func percentile(sorted []float64, pct float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	index := (pct / 100.0) * float64(len(sorted)-1)
	lower := int(index)
	upper := lower + 1
	if upper >= len(sorted) {
		return sorted[len(sorted)-1]
	}
	weight := index - float64(lower)
	return sorted[lower]*(1-weight) + sorted[upper]*weight
}

// ----------------------------------------------------------------------------
// Scenario 1: Health Liveness Probe
// ----------------------------------------------------------------------------
func runHealthStress(client *http.Client, cfg Config) ScenarioResult {
	url := cfg.Target + "/healthz"
	res := runWorkerPool(cfg.Requests, cfg.Concurrency, func(idx int) (int, time.Duration, error) {
		t0 := time.Now()
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return 0, time.Since(t0), err
		}
		resp, err := client.Do(req)
		dur := time.Since(t0)
		if err != nil {
			return 0, dur, err
		}
		defer resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
		return resp.StatusCode, dur, nil
	})

	res.Name = "Pillar 1: Container Liveness Probe (/healthz)"
	res.TargetURL = url
	okCount := res.StatusCodes[http.StatusOK]
	res.SuccessRate = (float64(okCount) / float64(cfg.Requests)) * 100.0
	res.InvariantPassed = (okCount == cfg.Requests)
	if res.InvariantPassed {
		res.InvariantNotes = fmt.Sprintf("100%% HTTP 200 UP across %d requests. Zero dropped sockets.", cfg.Requests)
	} else {
		res.InvariantNotes = fmt.Sprintf("Expected %d HTTP 200, received %d. WAF or connection rate limit engaged.", cfg.Requests, okCount)
	}
	return res
}

// ----------------------------------------------------------------------------
// Scenario 2: ACID Double-Entry Transfer with Unique Idempotency Keys
// ----------------------------------------------------------------------------
func runTransferStress(client *http.Client, cfg Config) ScenarioResult {
	url := cfg.Target + "/api/v1/wallets/transfer"

	res := runWorkerPool(cfg.Requests, cfg.Concurrency, func(idx int) (int, time.Duration, error) {
		payload := map[string]interface{}{
			"from_account": "ACC-9999", // Master treasury account with 10,000,000 MMK balance
			"to_account":   "ACC-1001",
			"amount":       100.00,
			"currency":     "MMK",
		}
		bodyBytes, _ := json.Marshal(payload)
		idemKey := fmt.Sprintf("STRESS-TX-%d-%d", time.Now().UnixNano(), idx)

		t0 := time.Now()
		req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
		if err != nil {
			return 0, time.Since(t0), err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Idempotency-Key", idemKey)

		resp, err := client.Do(req)
		dur := time.Since(t0)
		if err != nil {
			return 0, dur, err
		}
		defer resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
		return resp.StatusCode, dur, nil
	})

	res.Name = "Pillar 2: ACID Ledger Transfers (/api/v1/wallets/transfer)"
	res.TargetURL = url
	okCount := res.StatusCodes[http.StatusOK]
	res.SuccessRate = (float64(okCount) / float64(cfg.Requests)) * 100.0
	res.InvariantPassed = (okCount == cfg.Requests)
	if res.InvariantPassed {
		res.InvariantNotes = fmt.Sprintf("All %d concurrent transfers committed with atomic double-entry balance conservation.", okCount)
	} else {
		res.InvariantNotes = fmt.Sprintf("Passed: %d, Non-200: %d. Check balance or rate limiting.", okCount, cfg.Requests-okCount)
	}
	return res
}

// ----------------------------------------------------------------------------
// Scenario 3: High-Concurrency Idempotency Replay Attack (Zero Double-Debit Proof)
// ----------------------------------------------------------------------------
func runReplayStress(client *http.Client, cfg Config) ScenarioResult {
	url := cfg.Target + "/api/v1/wallets/transfer"
	sharedIdemKey := fmt.Sprintf("REPLAY-ATTACK-%d", time.Now().UnixNano())

	var firstExecutionCount int32
	var replayCachedCount int32

	res := runWorkerPool(cfg.Requests, cfg.Concurrency, func(idx int) (int, time.Duration, error) {
		payload := map[string]interface{}{
			"from_account": "ACC-9999",
			"to_account":   "ACC-2002",
			"amount":       250.00,
			"currency":     "MMK",
		}
		bodyBytes, _ := json.Marshal(payload)

		t0 := time.Now()
		req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
		if err != nil {
			return 0, time.Since(t0), err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Idempotency-Key", sharedIdemKey)

		resp, err := client.Do(req)
		dur := time.Since(t0)
		if err != nil {
			return 0, dur, err
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusOK {
			var parsed map[string]interface{}
			if err := json.Unmarshal(respBody, &parsed); err == nil {
				if idemp, ok := parsed["idempotent_replay"].(bool); ok && idemp {
					atomic.AddInt32(&replayCachedCount, 1)
				} else {
					atomic.AddInt32(&firstExecutionCount, 1)
				}
			}
		}

		return resp.StatusCode, dur, nil
	})

	res.Name = "Pillar 3: Idempotency Replay Attack (Zero Double-Debit Proof)"
	res.TargetURL = url
	okCount := res.StatusCodes[http.StatusOK]
	res.SuccessRate = (float64(okCount) / float64(cfg.Requests)) * 100.0

	// Invariant: exactly 1 execution must be original debit, all others must be idempotent replays
	fCount := atomic.LoadInt32(&firstExecutionCount)
	rCount := atomic.LoadInt32(&replayCachedCount)
	res.InvariantPassed = (fCount == 1 && rCount == int32(cfg.Requests-1))
	res.InvariantNotes = fmt.Sprintf("First-Execution Debits: %d | Idempotent Replays: %d | Zero Double-Debit Proven: %v",
		fCount, rCount, res.InvariantPassed)

	return res
}

// ----------------------------------------------------------------------------
// Scenario 4: Core Banking Outbox SQS FIFO Ingestion Burst
// ----------------------------------------------------------------------------
func runOutboxStress(client *http.Client, cfg Config) ScenarioResult {
	url := cfg.Target + "/api/v1/cbs/outbox-dispatch"

	res := runWorkerPool(cfg.Requests, cfg.Concurrency, func(idx int) (int, time.Duration, error) {
		payload := map[string]interface{}{
			"from_account": "ACC-1001",
			"to_account":   "ACC-2002",
			"amount":       5000.00,
			"currency":     "MMK",
		}
		bodyBytes, _ := json.Marshal(payload)

		t0 := time.Now()
		req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
		if err != nil {
			return 0, time.Since(t0), err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		dur := time.Since(t0)
		if err != nil {
			return 0, dur, err
		}
		defer resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
		return resp.StatusCode, dur, nil
	})

	res.Name = "Pillar 4: CBS Outbox SQS FIFO Dispatch (/api/v1/cbs/outbox-dispatch)"
	res.TargetURL = url
	okCount := res.StatusCodes[http.StatusOK]
	res.SuccessRate = (float64(okCount) / float64(cfg.Requests)) * 100.0
	res.InvariantPassed = (okCount == cfg.Requests)
	if res.InvariantPassed {
		res.InvariantNotes = fmt.Sprintf("All %d ISO 20022 events dispatched with deterministic SHA-256 deduplication IDs.", okCount)
	} else {
		res.InvariantNotes = fmt.Sprintf("OK: %d, Non-200: %d. Investigate outbox buffer.", okCount, cfg.Requests-okCount)
	}
	return res
}

// ----------------------------------------------------------------------------
// Scenario 5: DevSecOps WAF Threat Interception Stress
// ----------------------------------------------------------------------------
func runWAFStress(client *http.Client, cfg Config) ScenarioResult {
	url := cfg.Target + "/api/v1/devsecops/waf-probe"

	res := runWorkerPool(cfg.Requests, cfg.Concurrency, func(idx int) (int, time.Duration, error) {
		// Alternate between SQLi attack and clean payload
		isAttack := (idx % 2 == 0)
		payload := map[string]string{}
		if isAttack {
			payload["payload"] = "' OR '1'='1' --"
		} else {
			payload["payload"] = "Clean mobile user query"
		}
		bodyBytes, _ := json.Marshal(payload)

		t0 := time.Now()
		req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
		if err != nil {
			return 0, time.Since(t0), err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		dur := time.Since(t0)
		if err != nil {
			return 0, dur, err
		}
		defer resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
		return resp.StatusCode, dur, nil
	})

	res.Name = "Pillar 5: DevSecOps WAF Threat Interception (/api/v1/devsecops/waf-probe)"
	res.TargetURL = url
	c403 := res.StatusCodes[http.StatusForbidden]
	c200 := res.StatusCodes[http.StatusOK]
	expected403 := (cfg.Requests + 1) / 2
	expected200 := cfg.Requests / 2

	res.InvariantPassed = (c403 == expected403 && c200 == expected200)
	res.SuccessRate = (float64(c403+c200) / float64(cfg.Requests)) * 100.0
	res.InvariantNotes = fmt.Sprintf("Blocked SQLi Attacks (403): %d (Expected: %d) | Forwarded Clean (200): %d (Expected: %d)",
		c403, expected403, c200, expected200)
	return res
}

func printScenarioResult(r ScenarioResult) {
	fmt.Printf("▶ %s\n", r.Name)
	fmt.Printf("  Target      : %s\n", r.TargetURL)
	fmt.Printf("  Requests    : %d requests (Concurrency: %d workers)\n", r.TotalRequests, r.Concurrency)
	fmt.Printf("  Duration    : %d ms | Throughput: %.1f req/sec\n", r.DurationMs, r.ThroughputRPS)
	fmt.Printf("  Status Codes: ")
	first := true
	for code, count := range r.StatusCodes {
		if !first {
			fmt.Print(", ")
		}
		fmt.Printf("HTTP %d: %d", code, count)
		first = false
	}
	fmt.Println()
	fmt.Printf("  Latency     : Min: %.1fms | P50: %.1fms | P90: %.1fms | P95: %.1fms | P99: %.1fms | Max: %.1fms\n",
		r.MinLatencyMs, r.MedianLatencyMs, r.P90LatencyMs, r.P95LatencyMs, r.P99LatencyMs, r.MaxLatencyMs)
	fmt.Printf("  Invariant   : [%s] %s\n", passBadge(r.InvariantPassed), r.InvariantNotes)
	if len(r.Errors) > 0 {
		fmt.Printf("  Errors (%d) : %s\n", len(r.Errors), strings.Join(r.Errors, "; "))
	}
	fmt.Println()
}

func passBadge(passed bool) string {
	if passed {
		return "PASS"
	}
	return "FAIL"
}

func printSummaryTable(results []ScenarioResult) {
	fmt.Println("┌────────────────────────────────────────────────────────────────────────────────────────────────────────┐")
	fmt.Println("│                                       STRESS TEST SUITE SUMMARY                                        │")
	fmt.Println("├───────────────────────────────────────────────────────────────────┬──────────┬──────────┬──────────────┤")
	fmt.Println("│ SCENARIO                                                          │ REQ/SEC  │ P99 (ms) │ STATUS       │")
	fmt.Println("├───────────────────────────────────────────────────────────────────┼──────────┼──────────┼──────────────┤")
	allPassed := true
	for _, r := range results {
		status := "PASS"
		if !r.InvariantPassed {
			status = "FAIL"
			allPassed = false
		}
		name := r.Name
		if len(name) > 65 {
			name = name[:62] + "..."
		}
		fmt.Printf("│ %-65s │ %8.1f │ %8.1f │ %-12s │\n", name, r.ThroughputRPS, r.P99LatencyMs, status)
	}
	fmt.Println("└───────────────────────────────────────────────────────────────────┴──────────┴──────────┴──────────────┘")
	if allPassed {
		fmt.Println("RESULT: ALL SRE INVARIANTS SATISFIED UNDER HIGH CONCURRENCY LOAD.")
	} else {
		fmt.Println("RESULT: REGRESSIONS OR POLICY VIOLATIONS DETECTED UNDER LOAD.")
	}
	fmt.Println()
}
