package domain

import (
	"context"
	"sync"
	"time"

	"github.com/tcncloud/sati-go/pkg/ports"
)

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

// HostPluginProcess is now an interface defined in ports package
// The concrete implementation should be in the adapters package

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
		p.domain.hostPluginProcess.Stop()
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
		p.domain.hostPluginProcess.DispatchEvents(result.Events)
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
			p.domain.hostPluginProcess.DispatchJob(result.Job)
		}
	}

	return nil
}

func (p *StreamJobsProcess) stop() {
	if p.cancel != nil {
		p.cancel()
	}
}

// HostPluginProcess methods are now in the concrete implementation
// in the adapters package
