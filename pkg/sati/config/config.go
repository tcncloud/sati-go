// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//
// Copyright 2024 TCN Inc

// Package config provides configuration management functionality for the sati-go project.
// It handles loading, parsing, and watching configuration files that contain base64-encoded JSON data.
package config

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/tcncloud/sati-go/pkg/ports"
)

// Ensure ConfigWatcher implements the ports.ConfigWatcher interface
var _ ports.ConfigWatcher = (*ConfigWatcher)(nil)

// Ensure configLoader implements the ports.ConfigLoader interface
var _ ports.ConfigLoader = (*configLoader)(nil)

// Error constants for configuration operations.
var (
	ErrConfigPathsRequired = errors.New("config paths are required")
	ErrLoaderRequired      = errors.New("loader is required")
	ErrInvalidBase64       = errors.New("invalid base64 encoding")
	ErrInvalidJSON         = errors.New("invalid JSON format")
	ErrEmptyConfig         = errors.New("empty configuration")
	ErrRequiredField       = errors.New("required field is missing")
)

// Config represents the application configuration structure.
type Config struct {
	CACertificate           string `json:"ca_certificate"`
	Certificate             string `json:"certificate"`
	PrivateKey              string `json:"private_key"`
	FingerprintSHA256       string `json:"fingerprint_sha256"`
	FingerprintSHA256String string `json:"fingerprint_sha256_string"`
	APIEndpoint             string `json:"api_endpoint"`
	CertificateName         string `json:"certificate_name"`
	CertificateDescription  string `json:"certificate_description"`
}

// Validate checks if the configuration has all required fields.
func (c *Config) Validate() error {
	if c.APIEndpoint == "" {
		return fmt.Errorf("%w: api_endpoint", ErrRequiredField)
	}
	if c.CACertificate == "" {
		return fmt.Errorf("%w: ca_certificate", ErrRequiredField)
	}
	if c.Certificate == "" {
		return fmt.Errorf("%w: certificate", ErrRequiredField)
	}
	if c.PrivateKey == "" {
		return fmt.Errorf("%w: private_key", ErrRequiredField)
	}
	return nil
}

// decodeConfigFromBytes decodes base64-encoded JSON data into a Config struct.
func decodeConfigFromBytes(data []byte) (*Config, error) {
	if len(data) == 0 {
		return nil, ErrEmptyConfig
	}

	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(decoded, data)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidBase64, err)
	}

	var config Config
	if err := json.Unmarshal(decoded[:n], &config); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}

	return &config, nil
}

// LoadConfig loads configuration from a file path.
func LoadConfig(path string) (*Config, error) {
	// Note: path is validated to be a configuration file path, not user input
	//nolint:gosec // Configuration file path is controlled, not user input
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config, err := decodeConfigFromBytes(data)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// NewConfigFromString creates a Config object from a base64 encoded JSON string.
func NewConfigFromString(configString string) (*Config, error) {
	return decodeConfigFromBytes([]byte(configString))
}

// ConfigLoaderFunc defines the signature for configuration loader functions.
type ConfigLoaderFunc func(path string) error

// ConfigWatcher manages file watching for configuration changes.
// It implements the ports.ConfigWatcher interface.
type ConfigWatcher struct {
	watcher     *fsnotify.Watcher
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
	done        chan struct{}
	configPaths []string
	loader      ConfigLoaderFunc
	watching    bool
}

var (
	globalWatcher *ConfigWatcher
	globalMu      sync.RWMutex
)

// NewConfigWatcher creates a new configuration watcher.
func NewConfigWatcher(configPaths []string, loader ConfigLoaderFunc) (*ConfigWatcher, error) {
	if len(configPaths) == 0 {
		return nil, ErrConfigPathsRequired
	}
	if loader == nil {
		return nil, ErrLoaderRequired
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &ConfigWatcher{
		watcher:     watcher,
		ctx:         ctx,
		cancel:      cancel,
		done:        make(chan struct{}),
		configPaths: configPaths,
		loader:      loader,
		watching:    false,
	}, nil
}

// Start begins watching for configuration changes.
// It also reads the config file at startup if it exists.
func (cw *ConfigWatcher) Start(ctx context.Context) error {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	if cw.watching {
		return nil // Already watching
	}

	// Read config files at startup if they exist
	for _, configPath := range cw.configPaths {
		if _, err := os.Stat(configPath); err == nil {
			// File exists, load it
			if err := cw.loader(configPath); err != nil {
				log.Error().Err(err).Str("path", configPath).Msg("Failed to load config at startup")
			}
		}
	}

	// Add all paths to the watcher
	for _, configPath := range cw.configPaths {
		if err := cw.watcher.Add(configPath); err != nil {
			return fmt.Errorf("failed to add path %s to watcher: %w", configPath, err)
		}
	}

	// Start the watching goroutine
	go cw.watchLoop()

	cw.watching = true
	return nil
}

// watchLoop runs the main watching loop.
func (cw *ConfigWatcher) watchLoop() {
	defer close(cw.done)

	for {
		select {
		case <-cw.ctx.Done():
			return
		case event, ok := <-cw.watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				if err := cw.loader(event.Name); err != nil {
					log.Error().Err(err).Str("path", event.Name).Msg("Error in config loader")
				}
			}
		case err, ok := <-cw.watcher.Errors:
			if !ok {
				return
			}
			log.Error().Err(err).Msg("Error watching config file")
		}
	}
}

// Stop stops the watcher and cleans up resources.
func (cw *ConfigWatcher) Stop() error {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	if !cw.watching {
		return nil // Already stopped
	}

	if cw.cancel != nil {
		cw.cancel()
	}

	if cw.watcher != nil {
		if err := cw.watcher.Close(); err != nil {
			return err
		}
	}

	// Wait for the watch loop to finish
	<-cw.done

	cw.watching = false
	return nil
}

// IsWatching returns true if the watcher is currently active.
func (cw *ConfigWatcher) IsWatching() bool {
	cw.mu.RLock()
	defer cw.mu.RUnlock()
	return cw.watching
}

// WatchConfig starts watching configuration files using the global watcher.
// This function maintains backward compatibility with the existing API.
func WatchConfig(configPaths []string, loader ConfigLoaderFunc) error {
	globalMu.Lock()
	defer globalMu.Unlock()

	// Close existing watcher if any
	if globalWatcher != nil {
		if err := globalWatcher.Stop(); err != nil {
			log.Error().Err(err).Msg("Error stopping existing watcher")
		}
		globalWatcher = nil
	}

	// Create new watcher
	watcher, err := NewConfigWatcher(configPaths, loader)
	if err != nil {
		return err
	}

	// Start watching
	if err := watcher.Start(context.Background()); err != nil {
		watcher.Stop()
		return err
	}

	globalWatcher = watcher
	return nil
}

// StopWatching stops the global configuration watcher.
func StopWatching() error {
	globalMu.Lock()
	defer globalMu.Unlock()

	if globalWatcher == nil {
		return nil
	}

	err := globalWatcher.Stop()
	globalWatcher = nil
	return err
}

// LoadAndValidateConfig loads configuration from a file and validates it.
func LoadAndValidateConfig(path string) (*Config, error) {
	config, err := LoadConfig(path)
	if err != nil {
		return nil, err
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// LoadAndValidateConfigFromString loads configuration from a string and validates it.
func LoadAndValidateConfigFromString(configString string) (*Config, error) {
	config, err := NewConfigFromString(configString)
	if err != nil {
		return nil, err
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}
