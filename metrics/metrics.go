// Package metrics provides performance monitoring and statistics collection for the Discord bot.
package metrics

import (
	"errors"
	"sync"
	"time"

	botErrors "github.com/dunamismax/discogo/errors"
)

// Metrics holds all application metrics.
type Metrics struct {
	// Command metrics.
	CommandsTotal      int64
	CommandsSuccessful int64
	CommandsFailed     int64
	CommandsPerSecond  float64

	// API metrics.
	APIRequestsTotal      int64
	APIRequestsSuccessful int64
	APIRequestsFailed     int64
	APIRequestsPerSecond  float64
	APIResponseTimeSum    int64 // in milliseconds.
	APIResponseCount      int64

	// Error metrics by type.
	ErrorsByType map[botErrors.ErrorType]int64

	// Bot metrics.
	BotStartTime time.Time

	// Rate tracking.
	commandWindow *RateWindow
	apiWindow     *RateWindow
	mutex         sync.RWMutex
}

// RateWindow tracks events within a time window for rate calculations.
type RateWindow struct {
	events []time.Time
	window time.Duration
	mutex  sync.Mutex
}

// NewRateWindow creates a new rate tracking window.
func NewRateWindow(window time.Duration) *RateWindow {
	return &RateWindow{
		events: make([]time.Time, 0),
		window: window,
	}
}

// Add records an event timestamp.
func (rw *RateWindow) Add(timestamp time.Time) {
	rw.mutex.Lock()
	defer rw.mutex.Unlock()

	// Add new event.
	rw.events = append(rw.events, timestamp)

	// Remove events outside the window.
	cutoff := timestamp.Add(-rw.window)
	validEvents := make([]time.Time, 0, len(rw.events))

	for _, event := range rw.events {
		if event.After(cutoff) {
			validEvents = append(validEvents, event)
		}
	}

	rw.events = validEvents
}

// Rate calculates the current rate per second.
func (rw *RateWindow) Rate() float64 {
	rw.mutex.Lock()
	defer rw.mutex.Unlock()

	if len(rw.events) == 0 {
		return 0.0
	}

	// Remove expired events.
	now := time.Now()
	cutoff := now.Add(-rw.window)
	validEvents := 0

	for _, event := range rw.events {
		if event.After(cutoff) {
			validEvents++
		}
	}

	// Calculate rate per second.
	windowSeconds := rw.window.Seconds()

	return float64(validEvents) / windowSeconds
}

var globalMetrics *Metrics
var once sync.Once

// Initialize sets up the global metrics instance.
func Initialize() *Metrics {
	once.Do(func() {
		globalMetrics = &Metrics{
			ErrorsByType:  make(map[botErrors.ErrorType]int64),
			BotStartTime:  time.Now(),
			commandWindow: NewRateWindow(60 * time.Second), // 1-minute window.
			apiWindow:     NewRateWindow(60 * time.Second), // 1-minute window.
		}
	})

	return globalMetrics
}

// Get returns the global metrics instance.
func Get() *Metrics {
	if globalMetrics == nil {
		return Initialize()
	}

	return globalMetrics
}

// IncrementCommands increments command counters.
func (m *Metrics) IncrementCommands(successful bool) {
	now := time.Now()

	m.mutex.Lock()

	m.CommandsTotal++
	if successful {
		m.CommandsSuccessful++
	} else {
		m.CommandsFailed++
	}

	m.mutex.Unlock()

	m.commandWindow.Add(now)
	m.CommandsPerSecond = m.commandWindow.Rate()
}

// IncrementAPIRequests increments API request counters.
func (m *Metrics) IncrementAPIRequests(successful bool, responseTimeMs int64) {
	now := time.Now()

	m.mutex.Lock()

	m.APIRequestsTotal++
	if successful {
		m.APIRequestsSuccessful++
	} else {
		m.APIRequestsFailed++
	}

	m.APIResponseTimeSum += responseTimeMs
	m.APIResponseCount++
	m.mutex.Unlock()

	m.apiWindow.Add(now)
	m.APIRequestsPerSecond = m.apiWindow.Rate()
}

// IncrementError increments error counter by type.
func (m *Metrics) IncrementError(errorType botErrors.ErrorType) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.ErrorsByType[errorType]++
}

// GetAverageResponseTime calculates the average API response time.
func (m *Metrics) GetAverageResponseTime() float64 {
	m.mutex.RLock()
	responseTimeSum := m.APIResponseTimeSum
	responseCount := m.APIResponseCount
	m.mutex.RUnlock()

	if responseCount == 0 {
		return 0.0
	}

	return float64(responseTimeSum) / float64(responseCount)
}

// GetUptime returns the bot uptime.
func (m *Metrics) GetUptime() time.Duration {
	return time.Since(m.BotStartTime)
}

// GetSuccessRate calculates the command success rate as a percentage.
func (m *Metrics) GetSuccessRate() float64 {
	m.mutex.RLock()
	total := m.CommandsTotal
	successful := m.CommandsSuccessful
	m.mutex.RUnlock()

	if total == 0 {
		return 0.0
	}

	return (float64(successful) / float64(total)) * 100.0
}

// GetAPISuccessRate calculates the API success rate as a percentage.
func (m *Metrics) GetAPISuccessRate() float64 {
	m.mutex.RLock()
	total := m.APIRequestsTotal
	successful := m.APIRequestsSuccessful
	m.mutex.RUnlock()

	if total == 0 {
		return 0.0
	}

	return (float64(successful) / float64(total)) * 100.0
}

// Summary returns a comprehensive metrics summary.
type Summary struct {
	// Command statistics.
	CommandsTotal      int64   `json:"commands_total"`
	CommandsSuccessful int64   `json:"commands_successful"`
	CommandsFailed     int64   `json:"commands_failed"`
	CommandsPerSecond  float64 `json:"commands_per_second"`
	CommandSuccessRate float64 `json:"command_success_rate_percent"`

	// API statistics.
	APIRequestsTotal      int64   `json:"api_requests_total"`
	APIRequestsSuccessful int64   `json:"api_requests_successful"`
	APIRequestsFailed     int64   `json:"api_requests_failed"`
	APIRequestsPerSecond  float64 `json:"api_requests_per_second"`
	APISuccessRate        float64 `json:"api_success_rate_percent"`
	AverageResponseTime   float64 `json:"average_response_time_ms"`

	// Error statistics.
	ErrorsByType map[botErrors.ErrorType]int64 `json:"errors_by_type"`

	// System statistics.
	UptimeSeconds float64 `json:"uptime_seconds"`
	BotStartTime  string  `json:"bot_start_time"`
}

// GetSummary returns a comprehensive metrics summary.
func (m *Metrics) GetSummary() Summary {
	m.mutex.RLock()

	errorsByType := make(map[botErrors.ErrorType]int64)
	for k, v := range m.ErrorsByType {
		errorsByType[k] = v
	}

	m.mutex.RUnlock()

	m.mutex.RLock()
	summary := Summary{
		CommandsTotal:         m.CommandsTotal,
		CommandsSuccessful:    m.CommandsSuccessful,
		CommandsFailed:        m.CommandsFailed,
		CommandsPerSecond:     m.CommandsPerSecond,
		APIRequestsTotal:      m.APIRequestsTotal,
		APIRequestsSuccessful: m.APIRequestsSuccessful,
		APIRequestsFailed:     m.APIRequestsFailed,
		APIRequestsPerSecond:  m.APIRequestsPerSecond,
		UptimeSeconds:         m.GetUptime().Seconds(),
		BotStartTime:          m.BotStartTime.Format(time.RFC3339),
		ErrorsByType:          errorsByType,
	}

	commandSuccessRate := float64(0)
	if summary.CommandsTotal > 0 {
		commandSuccessRate = (float64(summary.CommandsSuccessful) / float64(summary.CommandsTotal)) * 100.0
	}

	apiSuccessRate := float64(0)
	if summary.APIRequestsTotal > 0 {
		apiSuccessRate = (float64(summary.APIRequestsSuccessful) / float64(summary.APIRequestsTotal)) * 100.0
	}

	averageResponseTime := float64(0)
	if m.APIResponseCount > 0 {
		averageResponseTime = float64(m.APIResponseTimeSum) / float64(m.APIResponseCount)
	}

	m.mutex.RUnlock()

	summary.CommandSuccessRate = commandSuccessRate
	summary.APISuccessRate = apiSuccessRate
	summary.AverageResponseTime = averageResponseTime

	return summary
}

// RecordCommand is a convenience function to record command execution.
func RecordCommand(successful bool) {
	Get().IncrementCommands(successful)
}

// RecordAPIRequest is a convenience function to record API requests.
func RecordAPIRequest(successful bool, responseTimeMs int64) {
	Get().IncrementAPIRequests(successful, responseTimeMs)
}

// RecordError is a convenience function to record errors.
func RecordError(err error) {
	var botErr *botErrors.BotError
	if errors.As(err, &botErr) {
		Get().IncrementError(botErr.Type)
	} else {
		Get().IncrementError(botErrors.ErrorTypeInternal)
	}
}
