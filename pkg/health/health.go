package health

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Details   map[string]string `json:"details"`
}

type HealthChecker struct {
	mu       sync.RWMutex
	status   HealthStatus
	checks   map[string]func() error
	interval time.Duration
}

func NewHealthChecker(interval time.Duration) *HealthChecker {
	hc := &HealthChecker{
		status: HealthStatus{
			Status:    "healthy",
			Timestamp: time.Now(),
			Details:   make(map[string]string),
		},
		checks:   make(map[string]func() error),
		interval: interval,
	}

	go hc.startHealthCheck()
	return hc
}

func (hc *HealthChecker) AddCheck(name string, check func() error) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.checks[name] = check
}

func (hc *HealthChecker) startHealthCheck() {
	ticker := time.NewTicker(hc.interval)
	defer ticker.Stop()

	for range ticker.C {
		hc.runChecks()
	}
}

func (hc *HealthChecker) runChecks() {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	status := "healthy"
	details := make(map[string]string)

	for name, check := range hc.checks {
		if err := check(); err != nil {
			status = "unhealthy"
			details[name] = err.Error()
		} else {
			details[name] = "ok"
		}
	}

	hc.status = HealthStatus{
		Status:    status,
		Timestamp: time.Now(),
		Details:   details,
	}
}

func (hc *HealthChecker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hc.mu.RLock()
	status := hc.status
	hc.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	if status.Status == "unhealthy" {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	json.NewEncoder(w).Encode(status)
} 