package main

import (
	"testing"
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
