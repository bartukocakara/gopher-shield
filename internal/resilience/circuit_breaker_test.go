package resilience

import (
	"errors"
	"testing"
	"time"
)

func TestCircuitBreaker(t *testing.T) {
	// 1. Setup: 3 failures to trip the circuit
	threshold := 3
	cb := NewCircuitBreaker(threshold)
	
	// Override reset timeout for faster testing
	cb.resetTimeout = 100 * time.Millisecond

	t.Run("Initially Closed", func(t *testing.T) {
		if cb.state != StateClosed {
			t.Errorf("Expected state Closed, got %v", cb.state)
		}
	})

	t.Run("Trips to Open after threshold", func(t *testing.T) {
		dummyErr := errors.New("service failure")
		
		// Fail 3 times
		for i := 0; i < threshold; i++ {
			_ = cb.Execute(func() error { return dummyErr })
		}

		if cb.state != StateOpen {
			t.Errorf("Expected state Open after %d failures, got %v", threshold, cb.state)
		}

		// Verify it stays open and returns ErrCircuitOpen
		err := cb.Execute(func() error { return nil })
		if err != ErrCircuitOpen {
			t.Errorf("Expected ErrCircuitOpen, got %v", err)
		}
	})

	t.Run("Resets after timeout", func(t *testing.T) {
		// Wait for the reset timeout to pass
		time.Sleep(150 * time.Millisecond)

		// The first execution after timeout should be allowed (probed)
		err := cb.Execute(func() error { return nil })
		
		if err != nil {
			t.Errorf("Expected successful probe, got %v", err)
		}

		if cb.state != StateClosed {
			t.Errorf("Expected state to reset to Closed, got %v", cb.state)
		}
	})
}