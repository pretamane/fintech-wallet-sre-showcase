package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
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

	// Return defensive copy
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

// Application encapsulates the HTTP server and routes
type Application struct {
	config Config
	ledger *LedgerEngine
	logger *log.Logger
}

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

	app := &Application{
		config: cfg,
		ledger: NewLedgerEngine(),
		logger: logger,
	}

	mux := http.NewServeMux()

	// Financial API Routes
	mux.HandleFunc("POST /api/v1/wallets/transfer", app.handleTransfer)
	mux.HandleFunc("GET /api/v1/wallets/{account_id}/balance", app.handleBalance)
	mux.HandleFunc("GET /api/v1/wallets", app.handleListAccounts)

	// SRE Observability & Health Probes (ECS / Kubernetes / Fargate)
	mux.HandleFunc("GET /healthz", app.handleHealthz)
	mux.HandleFunc("GET /readyz", app.handleReadyz)
	mux.HandleFunc("GET /metrics", app.handleMetrics)

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

func (app *Application) handleHealthz(w http.ResponseWriter, r *http.Request) {
	app.writeJSON(w, http.StatusOK, map[string]string{
		"status":  "UP",
		"service": "a-bank-wallet-service",
		"version": "1.0.0",
	})
}

func (app *Application) handleReadyz(w http.ResponseWriter, r *http.Request) {
	// In production, verifies database pool connectivity, redis lock manager, etc.
	app.writeJSON(w, http.StatusOK, map[string]string{
		"status": "READY",
		"ledger": "ONLINE",
	})
}

func (app *Application) handleMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := app.ledger.MetricsSummary()
	metrics["timestamp"] = time.Now().UTC()
	app.writeJSON(w, http.StatusOK, metrics)
}

func (app *Application) handleListAccounts(w http.ResponseWriter, r *http.Request) {
	accounts := app.ledger.GetAccounts()
	app.writeJSON(w, http.StatusOK, accounts)
}

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
				"root_dashboard": "/",
				"accounts":       "GET /api/v1/wallets",
				"balance":        "GET /api/v1/wallets/{account_id}/balance",
				"transfer":       "POST /api/v1/wallets/transfer",
				"health":         "GET /healthz",
				"readiness":      "GET /readyz",
				"metrics":        "GET /metrics",
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

