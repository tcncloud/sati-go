package domain

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/tcncloud/sati-go/pkg/ports"
)

// Domain is the main domain object for the application.
// It implements the following services:
// - ConfigWatcherHandler that will be called by the config watcher service when the
// config file is present and/or updated. When the config is changed this handler needs to restart the ExileClientConfiguration process.
//
// - ExileClientConfiguration process, this will be started by the ConfigWatcherHandler. This process will be responsible for fetching
// on an interval (every 60 seconds to 5 minutes) the client configuration using GetClientConfiguration method from the exile client
// verifying if the client configuration has changed and if so, it will restart the PollEventsProcess, StreamJobsProcess and HostPluginProcess.
//
// - PollEventsProcess - this will be started by the ExileClientConfiguration process. This process will be responsible for polling
// events from the exile client using the PollEvents method from the exile client and dispatching them to the HostPluginProcess.
//
// - StreamJobsProcess - this will be started by the ExileClientConfiguration process. This process will be responsible for streaming
// jobs from the exile client using the StreamJobs method from the exile client and dispatching them to the HostPluginProcess.
//
// - HostPluginProcess - this will be started by the ExileClientConfiguration process. This process will be responsible for hosting the plugin
// and dispatching the events and jobs to the plugin.
type Domain struct {
	log           *zerolog.Logger
	configWatcher ports.ConfigWatcher
	client        ports.ClientInterface

	// Process state
	mu                 sync.RWMutex
	exileConfigProcess *ExileClientConfigurationProcess
	pollEventsProcess  *PollEventsProcess
	streamJobsProcess  *StreamJobsProcess
	hostPluginProcess  *HostPluginProcess
	isRunning          bool
	shutdownChan       chan struct{}
}

// NewDomain creates a new Domain instance.
func NewDomain(log *zerolog.Logger) *Domain {
	return &Domain{
		log:          log,
		shutdownChan: make(chan struct{}),
	}
}

// SetConfigWatcher sets the configuration watcher for the domain.
func (d *Domain) SetConfigWatcher(watcher ports.ConfigWatcher) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.configWatcher = watcher
}

// SetClient sets the client interface for the domain.
func (d *Domain) SetClient(client ports.ClientInterface) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.client = client
}

// StartConfigWatcher starts the configuration watcher.
func (d *Domain) StartConfigWatcher(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.configWatcher == nil {
		d.log.Warn().Msg("Config watcher not configured, skipping startup")

		return nil
	}

	if err := d.configWatcher.Start(ctx); err != nil {
		d.log.Error().Err(err).Msg("Failed to start config watcher")

		return err
	}

	d.log.Info().Msg("Config watcher started successfully")

	return nil
}

// StartExileClientConfiguration starts the exile client configuration process.
func (d *Domain) StartExileClientConfiguration() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.client == nil {
		d.log.Warn().Msg("Client not configured, cannot start exile client configuration process")

		return nil
	}

	if d.exileConfigProcess != nil {
		d.log.Info().Msg("Exile client configuration process already running")

		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	process := &ExileClientConfigurationProcess{
		domain: d,
		cancel: cancel,
		ticker: time.NewTicker(ConfigCheckInterval),
	}

	d.exileConfigProcess = process
	go process.run(ctx)

	d.log.Info().Msg("Exile client configuration process started")

	return nil
}

// StartPollEvents starts the poll events process.
func (d *Domain) StartPollEvents() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.client == nil {
		d.log.Warn().Msg("Client not configured, cannot start poll events process")

		return nil
	}

	if d.pollEventsProcess != nil {
		d.log.Info().Msg("Poll events process already running")

		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	process := &PollEventsProcess{
		domain: d,
		cancel: cancel,
	}

	d.pollEventsProcess = process
	go process.run(ctx)

	d.log.Info().Msg("Poll events process started")

	return nil
}

// StartStreamJobs starts the stream jobs process.
func (d *Domain) StartStreamJobs() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.client == nil {
		d.log.Warn().Msg("Client not configured, cannot start stream jobs process")

		return nil
	}

	if d.streamJobsProcess != nil {
		d.log.Info().Msg("Stream jobs process already running")

		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	process := &StreamJobsProcess{
		domain: d,
		cancel: cancel,
	}

	d.streamJobsProcess = process
	go process.run(ctx)

	d.log.Info().Msg("Stream jobs process started")

	return nil
}

// StartHostPlugin starts the host plugin process.
func (d *Domain) StartHostPlugin() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.hostPluginProcess != nil {
		d.log.Info().Msg("Host plugin process already running")

		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	process := &HostPluginProcess{
		domain: d,
		cancel: cancel,
	}

	d.hostPluginProcess = process
	go process.run(ctx)

	d.log.Info().Msg("Host plugin process started")

	return nil
}

// StopAllProcesses stops all running processes.
func (d *Domain) StopAllProcesses() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.log.Info().Msg("Stopping all domain processes")

	// Stop processes in reverse order
	if d.hostPluginProcess != nil {
		d.hostPluginProcess.stop()
		d.hostPluginProcess = nil
	}

	if d.streamJobsProcess != nil {
		d.streamJobsProcess.stop()
		d.streamJobsProcess = nil
	}

	if d.pollEventsProcess != nil {
		d.pollEventsProcess.stop()
		d.pollEventsProcess = nil
	}

	if d.exileConfigProcess != nil {
		d.exileConfigProcess.stop()
		d.exileConfigProcess = nil
	}

	if d.configWatcher != nil {
		if err := d.configWatcher.Stop(); err != nil {
			d.log.Error().Err(err).Msg("Failed to stop config watcher")

			return err
		}
	}

	d.isRunning = false
	close(d.shutdownChan)

	d.log.Info().Msg("All domain processes stopped")

	return nil
}

// IsRunning returns true if the domain is currently running.
func (d *Domain) IsRunning() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.isRunning
}

// Ensure Domain implements DomainService interface.
var _ ports.DomainService = (*Domain)(nil)
