package domain

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/tcncloud/sati-go/pkg/ports"
)

// Constants for timeouts and intervals.
const (
	// ConfigCheckInterval is the interval for checking client configuration changes.
	ConfigCheckInterval = 60 * time.Second
	// DefaultTimeout is the default timeout for API calls.
	DefaultTimeout = 30 * time.Second
	// RetryDelay is the delay between retries for failed operations.
	RetryDelay = 5 * time.Second
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

// ExileClientConfigurationProcess manages the client configuration fetching and process coordination.
type ExileClientConfigurationProcess struct {
	domain     *Domain
	cancel     context.CancelFunc
	ticker     *time.Ticker
	lastConfig *ports.GetClientConfigurationResult
	mu         sync.RWMutex
}

// PollEventsProcess manages event polling from the exile client.
type PollEventsProcess struct {
	domain *Domain
	cancel context.CancelFunc
}

// StreamJobsProcess manages job streaming from the exile client.
type StreamJobsProcess struct {
	domain *Domain
	cancel context.CancelFunc
}

// HostPluginProcess manages the plugin hosting and event/job dispatching.
type HostPluginProcess struct {
	domain *Domain
	cancel context.CancelFunc
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

// ExileClientConfigurationProcess methods

func (p *ExileClientConfigurationProcess) run(ctx context.Context) {
	defer p.ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-p.ticker.C:
			if err := p.checkConfiguration(); err != nil {
				p.domain.log.Error().Err(err).Msg("Failed to check client configuration")
			}
		}
	}
}

func (p *ExileClientConfigurationProcess) checkConfiguration() error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	params := ports.GetClientConfigurationParams{}

	result, err := p.domain.client.GetClientConfiguration(ctx, params)
	if err != nil {
		return err
	}

	p.mu.Lock()
	configChanged := p.lastConfig == nil || p.hasConfigChanged(p.lastConfig, &result)
	p.lastConfig = &result
	p.mu.Unlock()

	if configChanged {
		p.domain.log.Info().Msg("Client configuration changed, restarting processes")
		p.restartProcesses()
	}

	return nil
}

func (p *ExileClientConfigurationProcess) hasConfigChanged(old, new *ports.GetClientConfigurationResult) bool {
	return old.OrgID != new.OrgID ||
		old.OrgName != new.OrgName ||
		old.ConfigName != new.ConfigName ||
		old.ConfigPayload != new.ConfigPayload
}

func (p *ExileClientConfigurationProcess) restartProcesses() {
	// Stop existing processes
	p.domain.mu.Lock()

	if p.domain.pollEventsProcess != nil {
		p.domain.pollEventsProcess.stop()
		p.domain.pollEventsProcess = nil
	}

	if p.domain.streamJobsProcess != nil {
		p.domain.streamJobsProcess.stop()
		p.domain.streamJobsProcess = nil
	}

	if p.domain.hostPluginProcess != nil {
		p.domain.hostPluginProcess.stop()
		p.domain.hostPluginProcess = nil
	}

	p.domain.mu.Unlock()

	// Restart processes
	if err := p.domain.StartPollEvents(); err != nil {
		p.domain.log.Error().Err(err).Msg("Failed to restart poll events process")
	}

	if err := p.domain.StartStreamJobs(); err != nil {
		p.domain.log.Error().Err(err).Msg("Failed to restart stream jobs process")
	}

	if err := p.domain.StartHostPlugin(); err != nil {
		p.domain.log.Error().Err(err).Msg("Failed to restart host plugin process")
	}
}

func (p *ExileClientConfigurationProcess) stop() {
	if p.cancel != nil {
		p.cancel()
	}
}

// PollEventsProcess methods

func (p *PollEventsProcess) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := p.pollEvents(); err != nil {
				p.domain.log.Error().Err(err).Msg("Failed to poll events")
				// Wait before retrying
				time.Sleep(RetryDelay)
			}
		}
	}
}

func (p *PollEventsProcess) pollEvents() error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	params := ports.PollEventsParams{}

	result, err := p.domain.client.PollEvents(ctx, params)
	if err != nil {
		return err
	}

	// Dispatch events to host plugin process
	if p.domain.hostPluginProcess != nil {
		p.domain.hostPluginProcess.dispatchEvents(result.Events)
	}

	return nil
}

func (p *PollEventsProcess) stop() {
	if p.cancel != nil {
		p.cancel()
	}
}

// StreamJobsProcess methods

func (p *StreamJobsProcess) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := p.streamJobs(); err != nil {
				p.domain.log.Error().Err(err).Msg("Failed to stream jobs")
				// Wait before retrying
				time.Sleep(RetryDelay)
			}
		}
	}
}

func (p *StreamJobsProcess) streamJobs() error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	params := ports.StreamJobsParams{}
	resultsChan := p.domain.client.StreamJobs(ctx, params)

	for result := range resultsChan {
		if result.Error != nil {
			return result.Error
		}

		// Dispatch job to host plugin process
		if p.domain.hostPluginProcess != nil {
			p.domain.hostPluginProcess.dispatchJob(result.Job)
		}
	}

	return nil
}

func (p *StreamJobsProcess) stop() {
	if p.cancel != nil {
		p.cancel()
	}
}

// HostPluginProcess methods

func (p *HostPluginProcess) run(ctx context.Context) {
	// This is where the plugin hosting logic would go
	// For now, it's a placeholder that handles events and jobs
	p.domain.log.Info().Msg("Host plugin process running")

	<-ctx.Done()
}

func (p *HostPluginProcess) dispatchEvents(events []ports.Event) {
	// Dispatch events to the plugin
	p.domain.log.Debug().Int("count", len(events)).Msg("Dispatching events to plugin")
	// Plugin dispatch logic would go here
}

func (p *HostPluginProcess) dispatchJob(job *ports.Job) {
	// Dispatch job to the plugin
	p.domain.log.Debug().Str("job_id", job.JobID).Str("type", job.Type).Msg("Dispatching job to plugin")
	// Plugin dispatch logic would go here
}

func (p *HostPluginProcess) stop() {
	if p.cancel != nil {
		p.cancel()
	}
}
