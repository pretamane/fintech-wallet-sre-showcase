package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Config holds runtime configuration passed via environment variables
type Config struct {
	Port         string
	Environment  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// TransferRequest defines the atomic transfer contract
type TransferRequest struct {
	FromAccount string  `json:"from_account"`
	ToAccount   string  `json:"to_account"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Reference   string  `json:"reference"`
}

// TransferResponse defines the successful transaction receipt
type TransferResponse struct {
	TransactionID string    `json:"transaction_id"`
	FromAccount   string    `json:"from_account"`
	ToAccount     string    `json:"to_account"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"`
	Timestamp     time.Time `json:"timestamp"`
	Idempotent    bool      `json:"idempotent_replay"`
}

// Account represents a financial ledger account
type Account struct {
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
	Currency  string  `json:"currency"`
}

// IdempotencyRecord stores past transaction executions
type IdempotencyRecord struct {
	RequestHash string
	Response    TransferResponse
	CreatedAt   time.Time
}

// LedgerEngine encapsulates the thread-safe double-entry ledger
type LedgerEngine struct {
	mu           sync.RWMutex
	accounts     map[string]*Account
	idempotency  map[string]IdempotencyRecord
	txCounter    uint64
	totalDebited float64
}

// NewLedgerEngine initializes sample accounts for testing
func NewLedgerEngine() *LedgerEngine {
	le := &LedgerEngine{
		accounts:    make(map[string]*Account),
		idempotency: make(map[string]IdempotencyRecord),
	}

	// Bootstrap demo accounts
	le.accounts["ACC-1001"] = &Account{AccountID: "ACC-1001", Balance: 500000.00, Currency: "MMK"}
	le.accounts["ACC-2002"] = &Account{AccountID: "ACC-2002", Balance: 150000.00, Currency: "MMK"}
	le.accounts["ACC-9999"] = &Account{AccountID: "ACC-9999", Balance: 10000000.00, Currency: "MMK"} // Central Reserve

	return le
}

// Transfer executes an atomic double-entry balance transfer with strict idempotency
func (le *LedgerEngine) Transfer(idempotencyKey string, req TransferRequest) (*TransferResponse, error) {
	if idempotencyKey == "" {
		return nil, errors.New("missing mandatory X-Idempotency-Key header")
	}

	if req.Amount <= 0 {
		return nil, errors.New("transfer amount must be strictly greater than zero")
	}

	if req.FromAccount == req.ToAccount {
		return nil, errors.New("source and destination accounts must be distinct")
	}

	reqHash := fmt.Sprintf("%s:%s:%.2f:%s", req.FromAccount, req.ToAccount, req.Amount, req.Currency)

	le.mu.Lock()
	defer le.mu.Unlock()

	// 1. Idempotency Check (Prevent duplicate charging on network retry)
	if existing, found := le.idempotency[idempotencyKey]; found {
		if existing.RequestHash != reqHash {
			return nil, fmt.Errorf("idempotency key conflict: key '%s' was already used with a different transaction payload", idempotencyKey)
		}
		// Return identical cached response marked as idempotent replay
		res := existing.Response
		res.Idempotent = true
		return &res, nil
	}

	// 2. Validate Accounts
	fromAcc, exists := le.accounts[req.FromAccount]
	if !exists {
		return nil, fmt.Errorf("source account '%s' not found", req.FromAccount)
	}

	toAcc, exists := le.accounts[req.ToAccount]
	if !exists {
		return nil, fmt.Errorf("destination account '%s' not found", req.ToAccount)
	}

	// 3. Balance Check
	if fromAcc.Balance < req.Amount {
		return nil, fmt.Errorf("insufficient funds: account '%s' has balance %.2f %s, requested %.2f %s",
			fromAcc.AccountID, fromAcc.Balance, fromAcc.Currency, req.Amount, req.Currency)
	}

	// 4. Atomic Execution (Double-Entry Debit & Credit)
	fromAcc.Balance -= req.Amount
	toAcc.Balance += req.Amount
	le.txCounter++
	le.totalDebited += req.Amount

	txID := fmt.Sprintf("TXN-%d-%d", time.Now().Unix(), le.txCounter)

	response := TransferResponse{
		TransactionID: txID,
		FromAccount:   req.FromAccount,
		ToAccount:     req.ToAccount,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Status:        "COMPLETED",
		Timestamp:     time.Now().UTC(),
		Idempotent:    false,
	}

	// Store idempotency record
	le.idempotency[idempotencyKey] = IdempotencyRecord{
		RequestHash: reqHash,
		Response:    response,
		CreatedAt:   time.Now().UTC(),
	}

	return &response, nil
}

// GetBalance retrieves account balance
func (le *LedgerEngine) GetBalance(accountID string) (*Account, error) {
	le.mu.RLock()
	defer le.mu.RUnlock()

	acc, exists := le.accounts[accountID]
	if !exists {
		return nil, fmt.Errorf("account '%s' not found", accountID)
	}

	return &Account{
		AccountID: acc.AccountID,
		Balance:   acc.Balance,
		Currency:  acc.Currency,
	}, nil
}

// GetAccounts returns defensive copies of all registered accounts
func (le *LedgerEngine) GetAccounts() []*Account {
	le.mu.RLock()
	defer le.mu.RUnlock()

	accounts := make([]*Account, 0, len(le.accounts))
	for _, acc := range le.accounts {
		accounts = append(accounts, &Account{
			AccountID: acc.AccountID,
			Balance:   acc.Balance,
			Currency:  acc.Currency,
		})
	}
	return accounts
}

// MetricsSummary returns internal SRE telemetry
func (le *LedgerEngine) MetricsSummary() map[string]interface{} {
	le.mu.RLock()
	defer le.mu.RUnlock()

	return map[string]interface{}{
		"total_transactions": le.txCounter,
		"total_volume":       le.totalDebited,
		"total_accounts":     len(le.accounts),
		"idempotency_keys":   len(le.idempotency),
	}
}

// ============================================================================
// SRE Alertmanager Webhook & Incident Management Engine
// ============================================================================

// AlertItem represents an alert payload item from Prometheus Alertmanager
type AlertItem struct {
	Status       string            `json:"status"` // "firing" | "resolved"
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

// AlertManagerWebhookPayload represents Prometheus Alertmanager v4 JSON format
type AlertManagerWebhookPayload struct {
	Receiver          string            `json:"receiver"`
	Status            string            `json:"status"` // "firing" | "resolved"
	Alerts            []AlertItem       `json:"alerts"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Version           string            `json:"version"`
}

// IncidentRecord stores active and historical operational incident states
type IncidentRecord struct {
	IncidentID     string     `json:"incident_id"`
	AlertName      string     `json:"alertname"`
	Severity       string     `json:"severity"`
	Status         string     `json:"status"` // "FIRING", "ACKNOWLEDGED", "RESOLVED"
	Summary        string     `json:"summary"`
	Description    string     `json:"description"`
	RunbookURL     string     `json:"runbook_url"`
	EscalationTier string     `json:"escalation_tier"`
	FiredAt        time.Time  `json:"fired_at"`
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty"`
	ResolvedAt     *time.Time `json:"resolved_at,omitempty"`
}

// AlertStore manages incident lifecycle thread-safely
type AlertStore struct {
	mu        sync.RWMutex
	incidents []*IncidentRecord
	counter   uint64
}

// NewAlertStore initializes an in-memory incident ledger
func NewAlertStore() *AlertStore {
	return &AlertStore{
		incidents: make([]*IncidentRecord, 0),
	}
}

// RecordAlert records an incident into the store
func (s *AlertStore) RecordAlert(alert AlertItem) *IncidentRecord {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	incidentID := fmt.Sprintf("INC-ABANK-%d-%03d", time.Now().Unix(), s.counter)
	alertName := alert.Labels["alertname"]
	if alertName == "" {
		alertName = "CBMNetPaymentRailDegraded"
	}
	severity := alert.Labels["severity"]
	if severity == "" {
		severity = "critical"
	}
	summary := alert.Annotations["summary"]
	if summary == "" {
		summary = "Downstream CBM-Net 2 payment clearing rail experiencing socket timeout"
	}
	desc := alert.Annotations["description"]
	if desc == "" {
		desc = "Circuit breaker tripped to OPEN; consecutive failures exceeded threshold on central bank gateway"
	}
	runbook := alert.Annotations["runbook"]
	if runbook == "" {
		runbook = "https://github.com/pretamane/fintech-wallet-sre-showcase/blob/main/README.md#runbook"
	}

	record := &IncidentRecord{
		IncidentID:     incidentID,
		AlertName:      alertName,
		Severity:       severity,
		Status:         "FIRING",
		Summary:        summary,
		Description:    desc,
		RunbookURL:     runbook,
		EscalationTier: "Tier 1 SRE On-Call (PagerDuty) -> Primary DevOps Lead",
		FiredAt:        time.Now().UTC(),
	}

	// Prepend newest first, bounded at 50 records
	s.incidents = append([]*IncidentRecord{record}, s.incidents...)
	if len(s.incidents) > 50 {
		s.incidents = s.incidents[:50]
	}
	return record
}

// GetIncidents returns defensive copies of all recorded incidents
func (s *AlertStore) GetIncidents() []*IncidentRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*IncidentRecord, len(s.incidents))
	for i, inc := range s.incidents {
		copyInc := *inc
		res[i] = &copyInc
	}
	return res
}

// AcknowledgeLatest marks the latest firing incident or specified ID as ACKNOWLEDGED
func (s *AlertStore) AcknowledgeLatest(incidentID string) (*IncidentRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	for _, inc := range s.incidents {
		if (incidentID == "" && inc.Status == "FIRING") || inc.IncidentID == incidentID {
			inc.Status = "ACKNOWLEDGED"
			inc.AcknowledgedAt = &now
			return inc, nil
		}
	}
	return nil, errors.New("no active firing incident found to acknowledge")
}

// ResolveLatest marks the incident as RESOLVED
func (s *AlertStore) ResolveLatest(incidentID string) (*IncidentRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	for _, inc := range s.incidents {
		if (incidentID == "" && (inc.Status == "FIRING" || inc.Status == "ACKNOWLEDGED")) || inc.IncidentID == incidentID {
			inc.Status = "RESOLVED"
			inc.ResolvedAt = &now
			return inc, nil
		}
	}
	return nil, errors.New("no active incident found to resolve")
}

// ============================================================================
// DevSecOps & FinTech Protocol Payload Structures
// ============================================================================

// WAFProbeRequest represents a security injection simulation
type WAFProbeRequest struct {
	Payload string `json:"payload"`
}

// WAFProbeResponse provides diagnostic feedback from edge security
type WAFProbeResponse struct {
	Status               string `json:"status"`
	HTTPStatus           int    `json:"http_status"`
	ThreatDetected       bool   `json:"threat_detected"`
	ThreatType           string `json:"threat_type,omitempty"`
	MatchedSignature     string `json:"matched_signature,omitempty"`
	WAFRuleID            string `json:"waf_rule_id,omitempty"`
	EdgeAction           string `json:"edge_action"`
	PCIComplianceControl string `json:"pci_compliance_control"`
}

// KMSEncryptRequest represents sensitive data to envelope encrypt
type KMSEncryptRequest struct {
	PlaintextPANOrNRC string `json:"plaintext_pan_or_nrc"`
}

// KMSEncryptResponse represents the AWS KMS CMK envelope-encrypted output
type KMSEncryptResponse struct {
	Status               string `json:"status"`
	KMSKeyARN            string `json:"kms_key_arn"`
	KMSKeyAlias          string `json:"kms_key_alias"`
	Algorithm            string `json:"algorithm"`
	CiphertextBlob       string `json:"ciphertext_blob"`
	MaskedDisplay        string `json:"masked_display"`
	KeyRotationPolicy    string `json:"key_rotation_policy"`
	PCIComplianceControl string `json:"pci_compliance_control"`
}

// OutboxDispatchRequest simulates staging an event into SQS FIFO
type OutboxDispatchRequest struct {
	FromAccount string  `json:"from_account"`
	ToAccount   string  `json:"to_account"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
}

// OutboxDispatchResponse represents the SQS FIFO queue dispatch receipt
type OutboxDispatchResponse struct {
	Status                 string    `json:"status"`
	QueueURL               string    `json:"queue_url"`
	MessageGroupID         string    `json:"message_group_id"`
	MessageDeduplicationID string    `json:"message_deduplication_id"`
	FinancialStandard      string    `json:"financial_standard"`
	CBSDecouplingGuarantee string    `json:"cbs_decoupling_guarantee"`
	DLQRetentionPeriod     string    `json:"dlq_retention_period"`
	Timestamp              time.Time `json:"timestamp"`
	PayloadSummary         string    `json:"payload_summary"`
}

// ============================================================================
// Application & Route Initialization
// ============================================================================

// Application encapsulates the HTTP server, state, and routes
type Application struct {
	config     Config
	ledger     *LedgerEngine
	logger     *log.Logger
	cbmBreaker *CircuitBreaker
	alertStore *AlertStore
}

var (
	sqliRegex      = regexp.MustCompile(`(?i)(union\s+select|select\s+.*from|drop\s+table|--|;\s*drop|or\s+1\s*=\s*1|and\s+1\s*=\s*1|'\s*or\s*')`)
	xssRegex       = regexp.MustCompile(`(?i)(<script|javascript:|onerror\s*=|onload\s*=|alert\(|<iframe|<img\s+src)`)
	traversalRegex = regexp.MustCompile(`(?i)(\.\./|/etc/passwd|/etc/shadow)`)
)

func main() {
	logger := log.New(os.Stdout, "[WALLET-SRE] ", log.LstdFlags|log.Lmicroseconds|log.LUTC)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "production"
	}

	cfg := Config{
		Port:         port,
		Environment:  env,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	breaker := NewCircuitBreaker("CBM-Net-Central-Bank-Rail", CircuitBreakerConfig{
		FailureThreshold: 5,
		SuccessThreshold: 2,
		CooldownWindow:   10 * time.Second,
	})

	app := &Application{
		config:     cfg,
		ledger:     NewLedgerEngine(),
		logger:     logger,
		cbmBreaker: breaker,
		alertStore: NewAlertStore(),
	}

	mux := http.NewServeMux()

	// Financial API Routes
	mux.HandleFunc("POST /api/v1/wallets/transfer", app.handleTransfer)
	mux.HandleFunc("GET /api/v1/wallets/{account_id}/balance", app.handleBalance)
	mux.HandleFunc("GET /api/v1/wallets", app.handleListAccounts)

	// Resilience Engineering & Chaos Testing (Circuit Breaker Controls)
	mux.HandleFunc("GET /api/v1/resilience/circuit-breaker", app.handleCircuitBreakerStatus)
	mux.HandleFunc("POST /api/v1/resilience/circuit-breaker/trip", app.handleCircuitBreakerTrip)
	mux.HandleFunc("POST /api/v1/resilience/circuit-breaker/reset", app.handleCircuitBreakerReset)

	// Production Alertmanager Webhook Integration & Incident Escalation
	mux.HandleFunc("POST /api/v1/alerts/webhook", app.handleAlertWebhook)
	mux.HandleFunc("GET /api/v1/alerts", app.handleAlertsList)
	mux.HandleFunc("POST /api/v1/alerts/simulate", app.handleAlertSimulate)
	mux.HandleFunc("POST /api/v1/alerts/acknowledge", app.handleAlertAcknowledge)
	mux.HandleFunc("POST /api/v1/alerts/resolve", app.handleAlertResolve)

	// DevSecOps & FinTech Protocol Sandboxes
	mux.HandleFunc("POST /api/v1/devsecops/waf-probe", app.handleWAFProbe)
	mux.HandleFunc("POST /api/v1/devsecops/kms-encrypt", app.handleKMSEncrypt)
	mux.HandleFunc("POST /api/v1/cbs/outbox-dispatch", app.handleOutboxDispatch)

	// SRE Observability & Health Probes (ECS / Kubernetes / Fargate)
	mux.HandleFunc("GET /healthz", app.handleHealthz)
	mux.HandleFunc("GET /readyz", app.handleReadyz)
	mux.HandleFunc("GET /metrics", app.handleMetrics)
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		_, _ = w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32"><rect width="32" height="32" rx="6" fill="#1e3a8a"/><text x="16" y="23" font-family="-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,sans-serif" font-size="20" font-weight="800" fill="#ffffff" text-anchor="middle">A</text></svg>`))
	})

	// Root Route: Interactive FinTech SRE Console (HTML / JSON)
	mux.HandleFunc("GET /", app.handleRoot)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      app.loggingMiddleware(mux),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	// Graceful shutdown handling
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Printf("Starting A Bank Wallet Service on port %s [env=%s]", cfg.Port, cfg.Environment)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("HTTP server fatal error: %v", err)
		}
	}()

	<-stopChan
	logger.Println("Received termination signal (SIGTERM/SIGINT). Commencing graceful draining...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server drained successfully. Shutdown completed.")
}

func (app *Application) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		app.logger.Printf("%s %s %s %s", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}

// ============================================================================
// Financial Transaction Handlers
// ============================================================================

func (app *Application) handleTransfer(w http.ResponseWriter, r *http.Request) {
	idempotencyKey := strings.TrimSpace(r.Header.Get("X-Idempotency-Key"))
	if idempotencyKey == "" {
		app.writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Missing mandatory 'X-Idempotency-Key' HTTP header",
		})
		return
	}

	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("Malformed JSON payload: %v", err),
		})
		return
	}

	res, err := app.ledger.Transfer(idempotencyKey, req)
	if err != nil {
		if strings.Contains(err.Error(), "conflict") {
			app.writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "insufficient funds") {
			app.writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
			return
		}
		app.writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	app.writeJSON(w, http.StatusOK, res)
}

func (app *Application) handleBalance(w http.ResponseWriter, r *http.Request) {
	accountID := r.PathValue("account_id")
	if accountID == "" {
		accountID = strings.TrimPrefix(r.URL.Path, "/api/v1/wallets/")
		accountID = strings.TrimSuffix(accountID, "/balance")
	}

	acc, err := app.ledger.GetBalance(accountID)
	if err != nil {
		app.writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	app.writeJSON(w, http.StatusOK, acc)
}

func (app *Application) handleListAccounts(w http.ResponseWriter, r *http.Request) {
	accounts := app.ledger.GetAccounts()
	app.writeJSON(w, http.StatusOK, accounts)
}

// ============================================================================
// Resilience & Chaos Handlers
// ============================================================================

func (app *Application) handleHealthz(w http.ResponseWriter, r *http.Request) {
	app.writeJSON(w, http.StatusOK, map[string]string{
		"status":  "UP",
		"service": "a-bank-wallet-service",
		"version": "1.0.0",
	})
}

func (app *Application) handleReadyz(w http.ResponseWriter, r *http.Request) {
	status := "READY"
	httpStatus := http.StatusOK
	if app.cbmBreaker.State() == StateOpen {
		status = "DEGRADED_DOWNSTREAM_RAIL_OPEN"
		httpStatus = http.StatusServiceUnavailable
	}

	app.writeJSON(w, httpStatus, map[string]interface{}{
		"status":          status,
		"ledger":          "ONLINE",
		"circuit_breaker": app.cbmBreaker.State().String(),
	})
}

func (app *Application) handleMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := app.ledger.MetricsSummary()
	metrics["timestamp"] = time.Now().UTC()
	metrics["circuit_breaker"] = app.cbmBreaker.Summary()
	metrics["active_incidents"] = len(app.alertStore.GetIncidents())
	app.writeJSON(w, http.StatusOK, metrics)
}

func (app *Application) handleCircuitBreakerStatus(w http.ResponseWriter, r *http.Request) {
	app.writeJSON(w, http.StatusOK, app.cbmBreaker.Summary())
}

func (app *Application) handleCircuitBreakerTrip(w http.ResponseWriter, r *http.Request) {
	app.cbmBreaker.ForceTrip()
	app.logger.Println("Circuit breaker manually TRIPPED to OPEN for chaos testing")

	// Automatically record an incident in AlertStore
	inc := app.alertStore.RecordAlert(AlertItem{
		Status: "firing",
		Labels: map[string]string{
			"alertname": "CBMNetPaymentRailDegraded",
			"severity":  "critical",
			"service":   "a-bank-wallet-service",
			"tier":      "core-clearing-rail",
		},
		Annotations: map[string]string{
			"summary":     "Downstream CBM-Net 2 payment clearing rail experiencing socket timeout",
			"description": "Circuit breaker forced to OPEN via SRE Chaos Injection; HTTP 503 fail-fast active",
			"runbook":     "https://github.com/pretamane/fintech-wallet-sre-showcase/blob/main/README.md#runbook",
		},
		StartsAt: time.Now().UTC(),
	})

	app.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message":         "Circuit breaker forced to OPEN state for chaos testing",
		"state":           app.cbmBreaker.State().String(),
		"incident_queued": inc,
	})
}

func (app *Application) handleCircuitBreakerReset(w http.ResponseWriter, r *http.Request) {
	app.cbmBreaker.ForceReset()
	app.logger.Println("Circuit breaker manually RESET to CLOSED")
	_, _ = app.alertStore.ResolveLatest("")

	app.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Circuit breaker reset to CLOSED state",
		"state":   app.cbmBreaker.State().String(),
	})
}

// ============================================================================
// Alertmanager Webhook & Pager Escalation Handlers
// ============================================================================

func (app *Application) handleAlertWebhook(w http.ResponseWriter, r *http.Request) {
	var payload AlertManagerWebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		app.writeJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Malformed Alertmanager payload: %v", err)})
		return
	}

	recorded := make([]*IncidentRecord, 0, len(payload.Alerts))
	for _, item := range payload.Alerts {
		if strings.ToLower(item.Status) == "firing" {
			rec := app.alertStore.RecordAlert(item)
			recorded = append(recorded, rec)
			app.logger.Printf("[P1 ESCALATION] Fired Alert: %s (Severity: %s, Incident: %s)", rec.AlertName, rec.Severity, rec.IncidentID)
		} else if strings.ToLower(item.Status) == "resolved" {
			_, _ = app.alertStore.ResolveLatest("")
			app.logger.Printf("[INCIDENT RESOLVED] Alert %s cleared", item.Labels["alertname"])
		}
	}

	app.writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":          "PROCESSED",
		"alerts_received": len(payload.Alerts),
		"incidents":       recorded,
	})
}

func (app *Application) handleAlertsList(w http.ResponseWriter, r *http.Request) {
	incidents := app.alertStore.GetIncidents()
	activeCount := 0
	for _, inc := range incidents {
		if inc.Status == "FIRING" || inc.Status == "ACKNOWLEDGED" {
			activeCount++
		}
	}

	app.writeJSON(w, http.StatusOK, map[string]interface{}{
		"active_count": activeCount,
		"total_count":  len(incidents),
		"incidents":    incidents,
	})
}

func (app *Application) handleAlertSimulate(w http.ResponseWriter, r *http.Request) {
	rec := app.alertStore.RecordAlert(AlertItem{
		Status: "firing",
		Labels: map[string]string{
			"alertname": "SLOErrorBudgetBurnHigh",
			"severity":  "critical",
			"service":   "a-bank-wallet-service",
			"burn_rate": "14.4x (P1 Pager Escalation)",
		},
		Annotations: map[string]string{
			"summary":     "Multi-window multi-burn-rate 14.4x threshold breached (1 hour / 5% budget consumed)",
			"description": "Payment error budget burn rate requires immediate Level 1 On-Call intervention",
			"runbook":     "https://github.com/pretamane/fintech-wallet-sre-showcase/blob/main/README.md#runbook",
		},
		StartsAt: time.Now().UTC(),
	})

	app.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message":  "Prometheus Alertmanager payload simulated successfully",
		"incident": rec,
	})
}

func (app *Application) handleAlertAcknowledge(w http.ResponseWriter, r *http.Request) {
	var body struct {
		IncidentID string `json:"incident_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	inc, err := app.alertStore.AcknowledgeLatest(body.IncidentID)
	if err != nil {
		app.writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	app.logger.Printf("[INCIDENT ACK] Incident %s acknowledged by SRE Lead", inc.IncidentID)
	app.writeJSON(w, http.StatusOK, inc)
}

func (app *Application) handleAlertResolve(w http.ResponseWriter, r *http.Request) {
	var body struct {
		IncidentID string `json:"incident_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	inc, err := app.alertStore.ResolveLatest(body.IncidentID)
	if err != nil {
		app.writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	// Also recover circuit breaker if it was tripped
	app.cbmBreaker.ForceReset()
	app.logger.Printf("[INCIDENT RESOLVED] Incident %s resolved. Rail recovered to CLOSED.", inc.IncidentID)
	app.writeJSON(w, http.StatusOK, inc)
}

// ============================================================================
// DevSecOps & FinTech Protocol Handlers
// ============================================================================

func (app *Application) handleWAFProbe(w http.ResponseWriter, r *http.Request) {
	var req WAFProbeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Malformed probe request"})
		return
	}

	payload := strings.TrimSpace(req.Payload)
	threatDetected := false
	threatType := ""
	matchedSig := ""
	wafRuleID := ""

	if sqliRegex.MatchString(payload) {
		threatDetected = true
		threatType = "SQL_INJECTION"
		matchedSig = "SQLi Pattern [UNION SELECT / OR 1=1 / DROP TABLE]"
		wafRuleID = "AWS-AWSManagedRulesCommonRuleSet-SQLiRule"
	} else if xssRegex.MatchString(payload) {
		threatDetected = true
		threatType = "CROSS_SITE_SCRIPTING"
		matchedSig = "XSS Pattern [<script> / onerror= / javascript:]"
		wafRuleID = "AWS-AWSManagedRulesCommonRuleSet-XSSRule"
	} else if traversalRegex.MatchString(payload) {
		threatDetected = true
		threatType = "PATH_TRAVERSAL"
		matchedSig = "Directory Traversal [../ / /etc/passwd]"
		wafRuleID = "AWS-AWSManagedRulesKnownBadInputsRuleSet"
	}

	if threatDetected {
		w.Header().Set("X-Security-Policy", "PCI-DSS-v4.0-Requirement-6.4")
		w.Header().Set("X-WAF-Protection", "AWS-WAFv2-RuleSet-Blocked")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		app.writeJSON(w, http.StatusForbidden, WAFProbeResponse{
			Status:               "BLOCKED",
			HTTPStatus:           http.StatusForbidden,
			ThreatDetected:       true,
			ThreatType:           threatType,
			MatchedSignature:     matchedSig,
			WAFRuleID:            wafRuleID,
			EdgeAction:           "TERMINATED_AT_WAF_INSPECTION_LAYER",
			PCIComplianceControl: "PCI-DSS v4.0 Req 6.4.2 Automated Technical Solution Protecting Public Facing Web Applications",
		})
		return
	}

	app.writeJSON(w, http.StatusOK, WAFProbeResponse{
		Status:               "ALLOWED",
		HTTPStatus:           http.StatusOK,
		ThreatDetected:       false,
		EdgeAction:           "FORWARDED_TO_UPSTREAM_TARGET_GROUP",
		PCIComplianceControl: "PCI-DSS v4.0 Clean Input Inspection",
	})
}

func (app *Application) handleKMSEncrypt(w http.ResponseWriter, r *http.Request) {
	var req KMSEncryptRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	val := strings.TrimSpace(req.PlaintextPANOrNRC)
	if val == "" {
		val = "12/DAGAMA(N)012345"
	}

	// Masking: Show first 4 and last 2, mask middle
	masked := val
	if len(val) > 6 {
		prefix := val[:4]
		suffix := val[len(val)-2:]
		masked = fmt.Sprintf("%s%s%s", prefix, strings.Repeat("*", len(val)-6), suffix)
	}

	hash := sha256.Sum256([]byte(val + "-a-bank-kms-salt-2026"))
	cipherBlob := fmt.Sprintf("kms:v1:aes256gcm:%s", hex.EncodeToString(hash[:]))

	app.writeJSON(w, http.StatusOK, KMSEncryptResponse{
		Status:               "ENCRYPTED",
		KMSKeyARN:            "arn:aws:kms:us-east-1:464868388812:key/5788a5cb-e914-4553-af87-6ba7bcb802b2",
		KMSKeyAlias:          "alias/a-bank-wallet-cmk",
		Algorithm:            "AES-256-GCM-256BIT-ENVELOPE",
		CiphertextBlob:       cipherBlob,
		MaskedDisplay:        masked,
		KeyRotationPolicy:    "AWS-KMS-365-DAY-AUTOMATIC-ANNUAL-ROTATION",
		PCIComplianceControl: "PCI-DSS v4.0 Req 3.5.1 Cryptographic Key Protection & Dual Control Separation",
	})
}

func (app *Application) handleOutboxDispatch(w http.ResponseWriter, r *http.Request) {
	var req OutboxDispatchRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	fromAcc := req.FromAccount
	if fromAcc == "" {
		fromAcc = "ACC-1001"
	}
	toAcc := req.ToAccount
	if toAcc == "" {
		toAcc = "ACC-2002"
	}
	amount := req.Amount
	if amount <= 0 {
		amount = 10000.00
	}
	currency := req.Currency
	if currency == "" {
		currency = "MMK"
	}

	now := time.Now().UTC()
	entropy := fmt.Sprintf("%s:%s:%.2f:%d", fromAcc, toAcc, amount, now.UnixNano())
	hash := sha256.Sum256([]byte(entropy))
	dedupID := hex.EncodeToString(hash[:])
	groupID := fmt.Sprintf("%s-SETTLEMENT-GROUP", fromAcc)

	summary := fmt.Sprintf("ISO 20022 pacs.008.001.08 [From: %s, To: %s, Amount: %.2f %s, SettlementTime: %s]",
		fromAcc, toAcc, amount, currency, now.Format(time.RFC3339))

	app.writeJSON(w, http.StatusOK, OutboxDispatchResponse{
		Status:                 "DISPATCHED_TO_SQS_FIFO",
		QueueURL:               "https://sqs.us-east-1.amazonaws.com/464868388812/a-bank-transaction-outbox.fifo",
		MessageGroupID:         groupID,
		MessageDeduplicationID: dedupID,
		FinancialStandard:      "ISO 20022 pacs.008.001.08 (Financial Institutional Customer Credit Transfer)",
		CBSDecouplingGuarantee: "ZERO_ROW_LOCKS_ON_CORE_BANKING_ORACLE_FLEXCUBE_ACCOUNTS",
		DLQRetentionPeriod:     "14 Days (1,209,600 seconds audit window)",
		Timestamp:              now,
		PayloadSummary:         summary,
	})
}

// ============================================================================
// Root & JSON Directory Handler
// ============================================================================

func (app *Application) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		app.writeJSON(w, http.StatusOK, map[string]interface{}{
			"service":              "A Bank FinTech Mobile Wallet Microservice",
			"division":             "Consumer Mobile Wallet & Digital Payments SRE",
			"status":               "OPERATIONAL",
			"runtime":              "AWS ECS Fargate (Serverless Container Tier)",
			"edge_gateway":         "Cloudflare Anycast Global CDN & WAF",
			"architecture_pattern": "Hybrid Cloud-Native FinTech Parity (CCE/Fargate + GaussDB/DynamoDB)",
			"endpoints": map[string]string{
				"root_dashboard":     "/",
				"accounts":           "GET /api/v1/wallets",
				"balance":            "GET /api/v1/wallets/{account_id}/balance",
				"transfer":           "POST /api/v1/wallets/transfer",
				"circuit_breaker":    "GET /api/v1/resilience/circuit-breaker",
				"circuit_trip":       "POST /api/v1/resilience/circuit-breaker/trip",
				"circuit_reset":      "POST /api/v1/resilience/circuit-breaker/reset",
				"alerts_webhook":     "POST /api/v1/alerts/webhook",
				"alerts_list":        "GET /api/v1/alerts",
				"alerts_simulate":    "POST /api/v1/alerts/simulate",
				"alerts_acknowledge": "POST /api/v1/alerts/acknowledge",
				"alerts_resolve":     "POST /api/v1/alerts/resolve",
				"devsecops_waf":      "POST /api/v1/devsecops/waf-probe",
				"devsecops_kms":      "POST /api/v1/devsecops/kms-encrypt",
				"cbs_outbox":         "POST /api/v1/cbs/outbox-dispatch",
				"health":             "GET /healthz",
				"readiness":          "GET /readyz",
				"metrics":            "GET /metrics",
			},
			"version": "1.0.0",
		})
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(rootDashboardHTML))
}

func (app *Application) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
