package main

import (
	"errors"
	"testing"
	"time"
)

func TestDoubleEntryTransfer(t *testing.T) {
	le := NewLedgerEngine()

	req := TransferRequest{
		FromAccount: "ACC-1001",
		ToAccount:   "ACC-2002",
		Amount:      5000.00,
		Currency:    "MMK",
		Reference:   "Test Transfer 1",
	}

	idempotencyKey := "IDEM-TEST-001"

	res, err := le.Transfer(idempotencyKey, req)
	if err != nil {
		t.Fatalf("Transfer failed unexpectedly: %v", err)
	}

	if res.Status != "COMPLETED" {
		t.Errorf("Expected status COMPLETED, got %s", res.Status)
	}

	if res.Idempotent {
		t.Errorf("Expected first transfer to be non-idempotent, got Idempotent=true")
	}

	// Verify FromAccount Balance (500,000 - 5,000 = 495,000)
	fromAcc, _ := le.GetBalance("ACC-1001")
	if fromAcc.Balance != 495000.00 {
		t.Errorf("Expected fromAcc balance 495000.00, got %.2f", fromAcc.Balance)
	}

	// Verify ToAccount Balance (150,000 + 5,000 = 155,000)
	toAcc, _ := le.GetBalance("ACC-2002")
	if toAcc.Balance != 155000.00 {
		t.Errorf("Expected toAcc balance 155000.00, got %.2f", toAcc.Balance)
	}
}

func TestIdempotencyReplay(t *testing.T) {
	le := NewLedgerEngine()

	req := TransferRequest{
		FromAccount: "ACC-1001",
		ToAccount:   "ACC-2002",
		Amount:      2000.00,
		Currency:    "MMK",
		Reference:   "Idempotency Test",
	}

	idempotencyKey := "IDEM-TEST-REPLAY"

	// First execution
	res1, err := le.Transfer(idempotencyKey, req)
	if err != nil {
		t.Fatalf("First transfer failed: %v", err)
	}

	// Second execution (Simulating mobile network retry)
	res2, err := le.Transfer(idempotencyKey, req)
	if err != nil {
		t.Fatalf("Second transfer replay failed: %v", err)
	}

	if !res2.Idempotent {
		t.Errorf("Expected res2 to have Idempotent=true")
	}

	if res1.TransactionID != res2.TransactionID {
		t.Errorf("Expected matching TransactionID on replay, got %s vs %s", res1.TransactionID, res2.TransactionID)
	}

	// Ensure balance was only debited ONCE (500,000 - 2,000 = 498,000)
	fromAcc, _ := le.GetBalance("ACC-1001")
	if fromAcc.Balance != 498000.00 {
		t.Errorf("Duplicate debit detected! Expected balance 498000.00, got %.2f", fromAcc.Balance)
	}
}

func TestInsufficientFunds(t *testing.T) {
	le := NewLedgerEngine()

	req := TransferRequest{
		FromAccount: "ACC-1001",
		ToAccount:   "ACC-2002",
		Amount:      999999999.00, // Exceeds balance
		Currency:    "MMK",
		Reference:   "Overdraft Test",
	}

	_, err := le.Transfer("IDEM-OVERDRAFT", req)
	if err == nil {
		t.Fatalf("Expected error for insufficient funds, got nil")
	}
}

func TestCircuitBreakerTripping(t *testing.T) {
	cb := NewCircuitBreaker("test-rail", CircuitBreakerConfig{
		FailureThreshold: 3,
		SuccessThreshold: 2,
		CooldownWindow:   100 * time.Millisecond,
	})

	if cb.State() != StateClosed {
		t.Fatalf("Expected initial state CLOSED, got %s", cb.State())
	}

	testErr := errors.New("upstream timeout")

	// Trigger 3 consecutive failures
	for i := 0; i < 3; i++ {
		err := cb.Execute(func() error {
			return testErr
		})
		if !errors.Is(err, testErr) {
			t.Fatalf("Expected testErr, got %v", err)
		}
	}

	// 4th call must fail fast with ErrCircuitOpen without executing the inner func
	invoked := false
	err := cb.Execute(func() error {
		invoked = true
		return nil
	})

	if !errors.Is(err, ErrCircuitOpen) {
		t.Fatalf("Expected ErrCircuitOpen, got %v", err)
	}
	if invoked {
		t.Fatalf("Expected operation to NOT be executed when circuit is OPEN")
	}
}

func TestCircuitBreakerRecovery(t *testing.T) {
	cb := NewCircuitBreaker("test-recovery-rail", CircuitBreakerConfig{
		FailureThreshold: 2,
		SuccessThreshold: 2,
		CooldownWindow:   50 * time.Millisecond,
	})

	cb.ForceTrip()
	if cb.State() != StateOpen {
		t.Fatalf("Expected state OPEN after ForceTrip, got %s", cb.State())
	}

	// Wait for cooldown to expire
	time.Sleep(60 * time.Millisecond)

	// In HALF-OPEN, first success
	err := cb.Execute(func() error {
		return nil
	})
	if err != nil {
		t.Fatalf("Expected probe to succeed, got %v", err)
	}

	// Second success closes the circuit
	err = cb.Execute(func() error {
		return nil
	})
	if err != nil {
		t.Fatalf("Expected second probe to succeed, got %v", err)
	}

	if cb.State() != StateClosed {
		t.Fatalf("Expected circuit to recover to CLOSED, got %s", cb.State())
	}
}

func TestAlertStoreLifecycle(t *testing.T) {
	store := NewAlertStore()

	alert := AlertItem{
		Status: "firing",
		Labels: map[string]string{
			"alertname": "CBMNetPaymentRailDegraded",
			"severity":  "critical",
		},
		Annotations: map[string]string{
			"summary": "Clearing socket timeout",
		},
		StartsAt: time.Now().UTC(),
	}

	rec := store.RecordAlert(alert)
	if rec.Status != "FIRING" {
		t.Fatalf("Expected status FIRING, got %s", rec.Status)
	}

	incidents := store.GetIncidents()
	if len(incidents) != 1 {
		t.Fatalf("Expected 1 incident, got %d", len(incidents))
	}

	// Acknowledge
	ack, err := store.AcknowledgeLatest(rec.IncidentID)
	if err != nil {
		t.Fatalf("Failed to acknowledge: %v", err)
	}
	if ack.Status != "ACKNOWLEDGED" || ack.AcknowledgedAt == nil {
		t.Fatalf("Expected ACKNOWLEDGED with timestamp, got %s", ack.Status)
	}

	// Resolve
	res, err := store.ResolveLatest(rec.IncidentID)
	if err != nil {
		t.Fatalf("Failed to resolve: %v", err)
	}
	if res.Status != "RESOLVED" || res.ResolvedAt == nil {
		t.Fatalf("Expected RESOLVED with timestamp, got %s", res.Status)
	}
}

func TestWAFProbeDetection(t *testing.T) {
	maliciousSQLi := []string{
		"' UNION SELECT * FROM accounts--",
		"1; DROP TABLE accounts;",
		"admin' OR 1=1--",
	}

	for _, payload := range maliciousSQLi {
		if !sqliRegex.MatchString(payload) {
			t.Errorf("Expected SQLi pattern match for '%s'", payload)
		}
	}

	maliciousXSS := []string{
		"<script>alert(1)</script>",
		"<img src=x onerror=alert('pci')>",
		"javascript:stealTokens()",
	}

	for _, payload := range maliciousXSS {
		if !xssRegex.MatchString(payload) {
			t.Errorf("Expected XSS pattern match for '%s'", payload)
		}
	}

	cleanPayload := "Retail Payment for Groceries"
	if sqliRegex.MatchString(cleanPayload) || xssRegex.MatchString(cleanPayload) {
		t.Errorf("Clean payload falsely flagged as threat: '%s'", cleanPayload)
	}
}
