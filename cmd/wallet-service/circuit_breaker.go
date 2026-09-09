package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// CircuitState represents the current operating state of the circuit breaker
type CircuitState int

const (
	StateClosed CircuitState = iota // Normal operations: requests flow through
	StateHalfOpen                   // Recovery probe: testing upstream health
	StateOpen                       // Outage tripped: fail fast without network I/O
)

func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateHalfOpen:
		return "HALF-OPEN"
	case StateOpen:
		return "OPEN"
	default:
		return "UNKNOWN"
	}
}

var (
	// ErrCircuitOpen is returned immediately when the circuit breaker is tripped
	ErrCircuitOpen = errors.New("circuit breaker is OPEN: downstream banking rail unavailable, failing fast")
)

// CircuitBreakerConfig defines threshold parameters for external banking rails
type CircuitBreakerConfig struct {
	FailureThreshold uint32        // Number of consecutive failures before tripping OPEN (e.g. 5)
	SuccessThreshold uint32        // Number of consecutive successes in HALF-OPEN to reset to CLOSED (e.g. 2)
	CooldownWindow   time.Duration // Duration to stay OPEN before testing recovery in HALF-OPEN (e.g. 10s)
}

// CircuitBreaker provides thread-safe failure protection against slow external banking APIs
type CircuitBreaker struct {
	name            string
	config          CircuitBreakerConfig
	mu              sync.RWMutex
	state           CircuitState
	consecutiveFail uint32
	consecutiveSucc uint32
	lastTrippedAt   time.Time
}

// NewCircuitBreaker initializes a circuit breaker for external payment rails (CBM-Net, MPU, Visa)
func NewCircuitBreaker(name string, cfg CircuitBreakerConfig) *CircuitBreaker {
	if cfg.FailureThreshold == 0 {
		cfg.FailureThreshold = 5
	}
	if cfg.SuccessThreshold == 0 {
		cfg.SuccessThreshold = 2
	}
	if cfg.CooldownWindow == 0 {
		cfg.CooldownWindow = 10 * time.Second
	}

	return &CircuitBreaker{
		name:   name,
		config: cfg,
		state:  StateClosed,
	}
}

// Execute wraps an external network operation with circuit breaker fail-fast mechanics
func (cb *CircuitBreaker) Execute(operation func() error) error {
	cb.mu.Lock()

	// Check if OPEN cooldown has elapsed, transitioning to HALF-OPEN
	if cb.state == StateOpen {
		if time.Since(cb.lastTrippedAt) > cb.config.CooldownWindow {
			cb.state = StateHalfOpen
			cb.consecutiveSucc = 0
		} else {
			cb.mu.Unlock()
			return ErrCircuitOpen
		}
	}

	cb.mu.Unlock()

	// Execute external operation
	err := operation()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.handleFailure()
		return err
	}

	cb.handleSuccess()
	return nil
}

func (cb *CircuitBreaker) handleFailure() {
	cb.consecutiveSucc = 0
	cb.consecutiveFail++

	if cb.state == StateHalfOpen || cb.consecutiveFail >= cb.config.FailureThreshold {
		cb.state = StateOpen
		cb.lastTrippedAt = time.Now()
	}
}

func (cb *CircuitBreaker) handleSuccess() {
	if cb.state == StateHalfOpen {
		cb.consecutiveSucc++
		if cb.consecutiveSucc >= cb.config.SuccessThreshold {
			cb.state = StateClosed
			cb.consecutiveFail = 0
			cb.consecutiveSucc = 0
		}
	} else if cb.state == StateClosed {
		cb.consecutiveFail = 0
	}
}

// State returns the current circuit state
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	if cb.state == StateOpen && time.Since(cb.lastTrippedAt) > cb.config.CooldownWindow {
		return StateHalfOpen
	}
	return cb.state
}

// Summary returns telemetry for observability endpoints
func (cb *CircuitBreaker) Summary() map[string]interface{} {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	currentState := cb.state
	if currentState == StateOpen && time.Since(cb.lastTrippedAt) > cb.config.CooldownWindow {
		currentState = StateHalfOpen
	}

	return map[string]interface{}{
		"rail_name":            cb.name,
		"state":                currentState.String(),
		"consecutive_failures": cb.consecutiveFail,
		"failure_threshold":    cb.config.FailureThreshold,
		"cooldown_seconds":     cb.config.CooldownWindow.Seconds(),
	}
}

// ForceTrip manually opens the circuit (used for chaos engineering and automated tests)
func (cb *CircuitBreaker) ForceTrip() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = StateOpen
	cb.lastTrippedAt = time.Now()
}

// ForceReset manually closes the circuit
func (cb *CircuitBreaker) ForceReset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = StateClosed
	cb.consecutiveFail = 0
	cb.consecutiveSucc = 0
}
