package health

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Status represents the health status of a component
type Status string

const (
	StatusUp      Status = "UP"
	StatusDown    Status = "DOWN"
	StatusDegraded Status = "DEGRADED"
)

// Check represents a single health check result
type Check struct {
	Name     string                 `json:"name"`
	Status   Status                 `json:"status"`
	Message  string                 `json:"message,omitempty"`
	Duration time.Duration          `json:"duration_ms"`
	Details  map[string]interface{} `json:"details,omitempty"`
}

// Report represents the overall health report
type Report struct {
	Status    Status           `json:"status"`
	Timestamp time.Time        `json:"timestamp"`
	Duration  time.Duration    `json:"duration_ms"`
	Checks    map[string]Check `json:"checks"`
}

// Checker is an interface for health checkers
type Checker interface {
	Check(ctx context.Context) Check
}

// Manager manages all health checks
type Manager struct {
	checkers map[string]Checker
	mu       sync.RWMutex
}

// NewManager creates a new health check manager
func NewManager() *Manager {
	return &Manager{
		checkers: make(map[string]Checker),
	}
}

// Register registers a health checker
func (m *Manager) Register(name string, checker Checker) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.checkers[name] = checker
}

// Check runs all health checks
func (m *Manager) Check(ctx context.Context) Report {
	start := time.Now()
	m.mu.RLock()
	defer m.mu.RUnlock()

	checks := make(map[string]Check)
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Run all checks in parallel
	for name, checker := range m.checkers {
		wg.Add(1)
		go func(n string, c Checker) {
			defer wg.Done()
			check := c.Check(ctx)
			mu.Lock()
			checks[n] = check
			mu.Unlock()
		}(name, checker)
	}

	wg.Wait()

	// Determine overall status
	overallStatus := StatusUp
	for _, check := range checks {
		if check.Status == StatusDown {
			overallStatus = StatusDown
			break
		} else if check.Status == StatusDegraded && overallStatus == StatusUp {
			overallStatus = StatusDegraded
		}
	}

	return Report{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Duration:  time.Since(start),
		Checks:    checks,
	}
}

// Handler returns an HTTP handler for health checks
func (m *Manager) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		report := m.Check(ctx)

		w.Header().Set("Content-Type", "application/json")

		// Set status code based on health
		statusCode := http.StatusOK
		if report.Status == StatusDown {
			statusCode = http.StatusServiceUnavailable
		} else if report.Status == StatusDegraded {
			statusCode = http.StatusOK // 200 but degraded
		}
		w.WriteHeader(statusCode)

		json.NewEncoder(w).Encode(report)
	}
}

// LivenessHandler returns a simple liveness probe (always returns 200 if app is running)
func LivenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

// ReadinessHandler returns a readiness probe (checks if app is ready to serve traffic)
func (m *Manager) ReadinessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		report := m.Check(ctx)

		if report.Status == StatusDown {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("NOT READY"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("READY"))
	}
}

// DatabaseChecker checks database connectivity
type DatabaseChecker struct {
	DB *sql.DB
}

func (c *DatabaseChecker) Check(ctx context.Context) Check {
	start := time.Now()
	check := Check{
		Name:    "database",
		Details: make(map[string]interface{}),
	}

	// Ping database
	if err := c.DB.PingContext(ctx); err != nil {
		check.Status = StatusDown
		check.Message = fmt.Sprintf("Failed to ping database: %v", err)
		check.Duration = time.Since(start)
		return check
	}

	// Get stats
	stats := c.DB.Stats()
	check.Details["open_connections"] = stats.OpenConnections
	check.Details["in_use"] = stats.InUse
	check.Details["idle"] = stats.Idle
	check.Details["max_open_connections"] = stats.MaxOpenConnections

	// Check if we're running out of connections
	if stats.OpenConnections >= stats.MaxOpenConnections-5 {
		check.Status = StatusDegraded
		check.Message = "Database connection pool is nearly exhausted"
	} else {
		check.Status = StatusUp
	}

	check.Duration = time.Since(start)
	return check
}

// RedisChecker checks Redis connectivity
type RedisChecker struct {
	Client *redis.Client
}

func (c *RedisChecker) Check(ctx context.Context) Check {
	start := time.Now()
	check := Check{
		Name:    "redis",
		Details: make(map[string]interface{}),
	}

	// Ping Redis
	if err := c.Client.Ping(ctx).Err(); err != nil {
		check.Status = StatusDown
		check.Message = fmt.Sprintf("Failed to ping Redis: %v", err)
		check.Duration = time.Since(start)
		return check
	}

	// Get info
	info, err := c.Client.Info(ctx, "stats").Result()
	if err == nil {
		check.Details["info"] = "available"
	}

	// Get memory stats
	memory, err := c.Client.Info(ctx, "memory").Result()
	if err == nil {
		check.Details["memory"] = "available"
	}

	// Simple check - if we got here, Redis is up
	_ = info
	_ = memory
	check.Status = StatusUp
	check.Duration = time.Since(start)
	return check
}

// DiskSpaceChecker checks available disk space
type DiskSpaceChecker struct {
	Path             string
	MinFreeBytes     int64
	MinFreePercent   float64
}

func (c *DiskSpaceChecker) Check(ctx context.Context) Check {
	start := time.Now()
	check := Check{
		Name:    "disk_space",
		Status:  StatusUp,
		Details: make(map[string]interface{}),
	}

	// Note: This is a placeholder. In production, you'd use syscall.Statfs
	// or a cross-platform library to get actual disk stats
	check.Details["path"] = c.Path
	check.Details["note"] = "Disk space check not implemented (platform-specific)"

	check.Duration = time.Since(start)
	return check
}

// CustomChecker allows creating custom health checks
type CustomChecker struct {
	Name     string
	CheckFn  func(ctx context.Context) (Status, string, map[string]interface{})
}

func (c *CustomChecker) Check(ctx context.Context) Check {
	start := time.Now()
	status, message, details := c.CheckFn(ctx)

	return Check{
		Name:     c.Name,
		Status:   status,
		Message:  message,
		Details:  details,
		Duration: time.Since(start),
	}
}