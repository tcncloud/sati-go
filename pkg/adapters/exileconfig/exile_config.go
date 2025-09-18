package exileconfig

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/tcncloud/sati-go/pkg/ports"
)

// ExileConfig is a service that monitors the exile client configuration and notifies the domain service when changes occur.
// It provides handlers that the exile client can call when it starts/stops, and runs a main loop to check for configuration changes.
type ExileConfig struct {
	exileClient   ports.ClientInterface
	domainService ports.DomainService
	log           *zerolog.Logger

	// State management
	mu         sync.RWMutex
	isRunning  bool
	lastConfig *ports.GetClientConfigurationResult
	stopChan   chan struct{}
	ticker     *time.Ticker
}

// NewExileConfig creates a new ExileConfig instance.
func NewExileConfig(
	exileClient ports.ClientInterface,
	domainService ports.DomainService,
	log *zerolog.Logger,
) *ExileConfig {
	return &ExileConfig{
		exileClient:   exileClient,
		domainService: domainService,
		log:           log,
		stopChan:      make(chan struct{}),
	}
}

// OnExileClientStarted is the handler that should be called by the exile client when it starts.
// This starts the configuration monitoring loop.
func (ec *ExileConfig) OnExileClientStarted(ctx context.Context) error {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	if ec.isRunning {
		ec.log.Warn().Msg("Exile config monitoring already running")
		return nil
	}

	ec.log.Info().Msg("Starting exile config monitoring")

	// Start the configuration monitoring goroutine
	go ec.monitorConfiguration(ctx)

	ec.isRunning = true
	ec.log.Info().Msg("Exile config monitoring started successfully")

	return nil
}

// OnExileClientStopped is the handler that should be called by the exile client when it stops.
// This stops the configuration monitoring loop.
func (ec *ExileConfig) OnExileClientStopped() error {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	if !ec.isRunning {
		ec.log.Warn().Msg("Exile config monitoring not running")
		return nil
	}

	ec.log.Info().Msg("Stopping exile config monitoring")

	// Signal stop
	close(ec.stopChan)

	// Stop ticker if running
	if ec.ticker != nil {
		ec.ticker.Stop()
		ec.ticker = nil
	}

	ec.isRunning = false
	ec.log.Info().Msg("Exile config monitoring stopped successfully")

	return nil
}

// IsRunning returns true if the configuration monitoring is currently active.
func (ec *ExileConfig) IsRunning() bool {
	ec.mu.RLock()
	defer ec.mu.RUnlock()

	return ec.isRunning
}

// monitorConfiguration runs the main loop that checks for configuration changes every 1 minute.
func (ec *ExileConfig) monitorConfiguration(ctx context.Context) {
	// Check configuration every 1 minute
	ec.ticker = time.NewTicker(1 * time.Minute)
	defer ec.ticker.Stop()

	ec.log.Info().Msg("Configuration monitoring loop started")

	for {
		select {
		case <-ctx.Done():
			ec.log.Info().Msg("Context cancelled, stopping config monitoring")
			return
		case <-ec.stopChan:
			ec.log.Info().Msg("Stop signal received, stopping config monitoring")
			return
		case <-ec.ticker.C:
			if err := ec.checkConfiguration(ctx); err != nil {
				ec.log.Error().Err(err).Msg("Failed to check configuration")
			}
		}
	}
}

// checkConfiguration checks the current configuration and handles changes.
func (ec *ExileConfig) checkConfiguration(ctx context.Context) error {
	// Create a timeout context for the configuration check
	checkCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Get current configuration
	params := ports.GetClientConfigurationParams{}
	result, err := ec.exileClient.GetClientConfiguration(checkCtx, params)
	if err != nil {
		return err
	}

	ec.mu.Lock()
	configChanged := ec.lastConfig == nil || ec.hasConfigurationChanged(ec.lastConfig, &result)
	oldConfig := ec.lastConfig
	ec.lastConfig = &result
	ec.mu.Unlock()

	if configChanged {
		ec.log.Info().
			Str("org_id", result.OrgID).
			Str("org_name", result.OrgName).
			Str("config_name", result.ConfigName).
			Msg("Configuration changed, notifying domain service")

		// Notify the domain service about the configuration change
		if err := ec.domainService.ClientConfigurationChanged(oldConfig, &result); err != nil {
			ec.log.Error().Err(err).Msg("Failed to notify domain service of configuration change")
			return err
		}

		ec.log.Info().Msg("Domain service notified of configuration change")
	}

	return nil
}

// hasConfigurationChanged checks if the configuration has changed.
func (ec *ExileConfig) hasConfigurationChanged(old, new *ports.GetClientConfigurationResult) bool {
	if old == nil || new == nil {
		return old != new
	}

	return old.OrgID != new.OrgID ||
		old.OrgName != new.OrgName ||
		old.ConfigName != new.ConfigName ||
		old.ConfigPayload != new.ConfigPayload
}

// GetLastConfiguration returns the last known configuration.
func (ec *ExileConfig) GetLastConfiguration() *ports.GetClientConfigurationResult {
	ec.mu.RLock()
	defer ec.mu.RUnlock()

	return ec.lastConfig
}

// ConfigChangeHandler interface implementation

// OnConfigChanged is called when the configuration changes.
func (ec *ExileConfig) OnConfigChanged(oldConfig, newConfig *ports.GetClientConfigurationResult) {
	ec.log.Info().
		Str("old_org_id", oldConfig.OrgID).
		Str("new_org_id", newConfig.OrgID).
		Str("old_config_name", oldConfig.ConfigName).
		Str("new_config_name", newConfig.ConfigName).
		Msg("Configuration change detected")

	// Log the change details
	ec.log.Debug().
		Interface("old_config", oldConfig).
		Interface("new_config", newConfig).
		Msg("Configuration change details")
}

// ShouldRestartProcesses returns true if processes should be restarted.
func (ec *ExileConfig) ShouldRestartProcesses(oldConfig, newConfig *ports.GetClientConfigurationResult) bool {
	if oldConfig == nil || newConfig == nil {
		return oldConfig != newConfig
	}

	// Restart processes if any critical configuration changed
	return oldConfig.OrgID != newConfig.OrgID ||
		oldConfig.OrgName != newConfig.OrgName ||
		oldConfig.ConfigName != newConfig.ConfigName ||
		oldConfig.ConfigPayload != newConfig.ConfigPayload
}

// ConfigWatcher interface implementation

// Start begins watching for configuration changes.
func (ec *ExileConfig) Start(ctx context.Context) error {
	return ec.OnExileClientStarted(ctx)
}

// Stop stops the configuration watcher.
func (ec *ExileConfig) Stop() error {
	return ec.OnExileClientStopped()
}

// IsWatching returns true if the watcher is currently active.
func (ec *ExileConfig) IsWatching() bool {
	return ec.IsRunning()
}

// Ensure ExileConfig implements required interfaces.
var _ ports.ConfigWatcher = (*ExileConfig)(nil)
var _ ports.ConfigChangeHandler = (*ExileConfig)(nil)
