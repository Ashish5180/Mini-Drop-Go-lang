package common

import (
	"runtime"
	"sync/atomic"
	"time"
)

// PerformanceMetrics tracks system performance metrics
type PerformanceMetrics struct {
	TotalRequests    atomic.Int64
	TotalBytes       atomic.Int64
	TotalErrors      atomic.Int64
	AverageLatencyMs atomic.Int64
	StartTime        time.Time
}

// NewPerformanceMetrics creates a new performance metrics tracker
func NewPerformanceMetrics() *PerformanceMetrics {
	return &PerformanceMetrics{
		StartTime: time.Now(),
	}
}

// Record updates metrics atomically
func (pm *PerformanceMetrics) Record(bytes int64, latencyMs int64, isError bool) {
	pm.TotalRequests.Add(1)
	if bytes > 0 {
		pm.TotalBytes.Add(bytes)
	}
	if isError {
		pm.TotalErrors.Add(1)
	}
	if latencyMs > 0 {
		pm.AverageLatencyMs.Store(latencyMs)
	}
}

// GetStats returns current statistics
func (pm *PerformanceMetrics) GetStats() map[string]interface{} {
	uptime := time.Since(pm.StartTime)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"uptime_seconds":  uptime.Seconds(),
		"total_requests":  pm.TotalRequests.Load(),
		"total_bytes":     pm.TotalBytes.Load(),
		"total_errors":    pm.TotalErrors.Load(),
		"avg_latency_ms":  pm.AverageLatencyMs.Load(),
		"goroutines":      runtime.NumGoroutine(),
		"memory_alloc_mb": m.Alloc / 1024 / 1024,
		"memory_sys_mb":   m.Sys / 1024 / 1024,
		"gc_runs":         m.NumGC,
	}
}

// MiddlewareTimer tracks request timing
type MiddlewareTimer struct {
	start time.Time
}

// NewTimer creates a new timing context
func NewTimer() *MiddlewareTimer {
	return &MiddlewareTimer{start: time.Now()}
}

// Elapsed returns elapsed time in milliseconds
func (t *MiddlewareTimer) Elapsed() int64 {
	return time.Since(t.start).Milliseconds()
}
