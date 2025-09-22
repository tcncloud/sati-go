package hostplugin

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/tcncloud/sati-go/pkg/ports"
)

// HostPluginProcess implements the ports.HostPluginProcess interface.
type HostPluginProcess struct {
	log    *zerolog.Logger
	cancel context.CancelFunc
}

// NewHostPluginProcess creates a new HostPluginProcess instance.
func NewHostPluginProcess(log *zerolog.Logger) *HostPluginProcess {
	return &HostPluginProcess{
		log: log,
	}
}

// Run starts the host plugin process.
func (p *HostPluginProcess) Run(ctx context.Context) {
	// This is where the plugin hosting logic would go
	// For now, it's a placeholder that handles events and jobs
	p.log.Info().Msg("Host plugin process running")

	<-ctx.Done()
}

// Stop stops the host plugin process.
func (p *HostPluginProcess) Stop() {
	if p.cancel != nil {
		p.cancel()
	}
}

// DispatchEvents dispatches events to the plugin.
func (p *HostPluginProcess) DispatchEvents(events []ports.Event) {
	// Dispatch events to the plugin
	p.log.Debug().Int("count", len(events)).Msg("Dispatching events to plugin")
	// Plugin dispatch logic would go here
}

// DispatchJob dispatches a job to the plugin.
func (p *HostPluginProcess) DispatchJob(job *ports.Job) {
	// Dispatch job to the plugin
	p.log.Debug().Str("job_id", job.JobID).Str("type", job.Type).Msg("Dispatching job to plugin")
	// Plugin dispatch logic would go here
}

// Ensure HostPluginProcess implements ports.HostPluginProcess interface.
var _ ports.HostPluginProcess = (*HostPluginProcess)(nil)

