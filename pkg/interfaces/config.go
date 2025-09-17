// Package interfaces defines the interfaces for the application.
package interfaces

import "context"

// ConfigWatcher defines the interface for configuration watching functionality.
type ConfigWatcher interface {
	// Start begins watching for configuration changes.
	Start(ctx context.Context) error

	// Stop stops the configuration watcher.
	Stop() error

	// IsWatching returns true if the watcher is currently active.
	IsWatching() bool
}

// ConfigLoader defines the interface for configuration loading functionality.
type ConfigLoader interface {
	// LoadConfig loads configuration from a file path.
	LoadConfig(path string) (interface{}, error)

	// LoadConfigFromString loads configuration from a base64 encoded string.
	LoadConfigFromString(configString string) (interface{}, error)
}
