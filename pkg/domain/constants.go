package domain

import "time"

// Constants for timeouts and intervals.
const (
	// ConfigCheckInterval is the interval for checking client configuration changes.
	ConfigCheckInterval = 60 * time.Second
	// DefaultTimeout is the default timeout for API calls.
	DefaultTimeout = 30 * time.Second
	// RetryDelay is the delay between retries for failed operations.
	RetryDelay = 5 * time.Second
)
