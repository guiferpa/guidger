package server

import (
	"net/http"
	"sync"
	"sync/atomic"
)

// Telemetry tracks metrics for the HTTP server
type Telemetry struct {
	mu sync.RWMutex

	// Request counters (using atomic for thread-safety and performance)
	totalRequests   int64
	successRequests int64
	errorRequests   int64
	webhookRequests int64
	balanceRequests int64

	// Error counts by status code (protected by mutex)
	errorCounts map[int]int64
}

// NewTelemetry creates a new telemetry instance
func NewTelemetry() *Telemetry {
	return &Telemetry{
		errorCounts: make(map[int]int64),
	}
}

// IncrementTotalRequests increments the total request counter
func (t *Telemetry) IncrementTotalRequests() {
	atomic.AddInt64(&t.totalRequests, 1)
}

// IncrementSuccessRequests increments the success request counter
func (t *Telemetry) IncrementSuccessRequests() {
	atomic.AddInt64(&t.successRequests, 1)
}

// IncrementErrorRequests increments the error request counter
func (t *Telemetry) IncrementErrorRequests(statusCode int) {
	atomic.AddInt64(&t.errorRequests, 1)
	t.mu.Lock()
	t.errorCounts[statusCode]++
	t.mu.Unlock()
}

// IncrementWebhookRequests increments the webhook request counter
func (t *Telemetry) IncrementWebhookRequests() {
	atomic.AddInt64(&t.webhookRequests, 1)
}

// IncrementBalanceRequests increments the balance request counter
func (t *Telemetry) IncrementBalanceRequests() {
	atomic.AddInt64(&t.balanceRequests, 1)
}

// GetMetrics returns a snapshot of current metrics
func (t *Telemetry) GetMetrics() map[string]interface{} {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// Create a copy of errorCounts to avoid race conditions
	errorCountsCopy := make(map[int]int64)
	for k, v := range t.errorCounts {
		errorCountsCopy[k] = v
	}

	return map[string]interface{}{
		"total_requests":   atomic.LoadInt64(&t.totalRequests),
		"success_requests": atomic.LoadInt64(&t.successRequests),
		"error_requests":   atomic.LoadInt64(&t.errorRequests),
		"webhook_requests": atomic.LoadInt64(&t.webhookRequests),
		"balance_requests": atomic.LoadInt64(&t.balanceRequests),
		"error_counts":     errorCountsCopy,
	}
}

// TelemetryMiddleware tracks request metrics
func TelemetryMiddleware(telemetry *Telemetry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			telemetry.IncrementTotalRequests()

			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rw, r)

			// Track success/error based on status code
			if rw.statusCode >= 200 && rw.statusCode < 400 {
				telemetry.IncrementSuccessRequests()
			} else {
				telemetry.IncrementErrorRequests(rw.statusCode)
			}

			// Track endpoint-specific metrics
			path := r.URL.Path
			if path == "/webhook" || path == "/webhook/" {
				telemetry.IncrementWebhookRequests()
			} else if len(path) >= 8 && path[:8] == "/balance" && path != "/metrics" {
				telemetry.IncrementBalanceRequests()
			}
		})
	}
}
