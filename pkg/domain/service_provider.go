package domain

import "context"

// DomainServiceProvider provides clean access to domain service methods.
type DomainServiceProvider struct {
	domain *Domain
}

// NewDomainServiceProvider creates a new domain service provider.
func NewDomainServiceProvider(domain *Domain) *DomainServiceProvider {
	return &DomainServiceProvider{
		domain: domain,
	}
}

// GetConfigWatcherStarter returns the config watcher starter function.
func (dsp *DomainServiceProvider) GetConfigWatcherStarter() func(context.Context) error {
	return dsp.domain.StartConfigWatcher
}

// GetExileClientConfigurationStarter returns the exile client configuration starter function.
func (dsp *DomainServiceProvider) GetExileClientConfigurationStarter() func() error {
	return dsp.domain.StartExileClientConfiguration
}

// GetPollEventsStarter returns the poll events starter function.
func (dsp *DomainServiceProvider) GetPollEventsStarter() func() error {
	return dsp.domain.StartPollEvents
}

// GetStreamJobsStarter returns the stream jobs starter function.
func (dsp *DomainServiceProvider) GetStreamJobsStarter() func() error {
	return dsp.domain.StartStreamJobs
}

// GetHostPluginStarter returns the host plugin starter function.
func (dsp *DomainServiceProvider) GetHostPluginStarter() func() error {
	return dsp.domain.StartHostPlugin
}

// GetStopAllProcessesStopper returns the stop all processes function.
func (dsp *DomainServiceProvider) GetStopAllProcessesStopper() func() error {
	return dsp.domain.StopAllProcesses
}

// GetIsRunningChecker returns the is running checker function.
func (dsp *DomainServiceProvider) GetIsRunningChecker() func() bool {
	return dsp.domain.IsRunning
}
