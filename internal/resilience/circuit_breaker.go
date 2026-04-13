package resilience

import (
	"errors"
	"log"
	"sync"
	"time"
)

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

var ErrCircuitOpen = errors.New("circuit breaker is open")

type CircuitBreaker struct {
	mu               sync.RWMutex
	state            State
	failureThreshold int
	failureCount     int
	lastFailureTime  time.Time
	resetTimeout     time.Duration
	OnStateChange    func(string) // Hook for metrics
}

func NewCircuitBreaker(threshold int) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: threshold,
		state:            StateClosed,
		resetTimeout:     10 * time.Second,
	}
}

func (cb *CircuitBreaker) Execute(op func() error) error {
	if !cb.canExecute() {
		return ErrCircuitOpen
	}

	err := op()
	if err != nil {
		cb.recordFailure()
		return err
	}

	cb.recordSuccess()
	return nil
}

func (cb *CircuitBreaker) canExecute() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	if cb.state == StateClosed {
		return true
	}
	if cb.state == StateOpen {
		if time.Since(cb.lastFailureTime) > cb.resetTimeout {
			log.Println("DEBUG: Circuit entering Half-Open probe state...") // <--- BU LOGU EKLE
			return true
		}
	}
	return cb.state == StateHalfOpen
}

func (cb *CircuitBreaker) recordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failureCount++
	cb.lastFailureTime = time.Now()

	if cb.failureCount >= cb.failureThreshold && cb.state != StateOpen {
		cb.state = StateOpen
		if cb.OnStateChange != nil {
			cb.OnStateChange("open")
		}
	}
}

func (cb *CircuitBreaker) recordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failureCount = 0
	if cb.state != StateClosed {
		cb.state = StateClosed
		if cb.OnStateChange != nil {
			cb.OnStateChange("closed")
		}
	}
}