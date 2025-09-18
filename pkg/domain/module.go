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

package domain

import (
	"context"

	"github.com/tcncloud/sati-go/pkg/ports"
	"go.uber.org/fx"
)

// Module provides the domain module for dependency injection.
// It includes the main Domain service and related business logic.
//
// Usage example:
//
//	app := fx.New(
//	  domain.Module,
//	  fx.Invoke(func(d *Domain) {
//	    // Use domain service
//	    ctx := context.Background()
//	    d.StartConfigWatcher(ctx)
//	  }),
//	)
var Module = fx.Module("domain",
	// Provide the main Domain service
	fx.Provide(NewDomain),

	// Provide the domain service provider
	fx.Provide(NewDomainServiceProvider),

	// Provide domain service methods as injectable functions
	fx.Provide(func(d *Domain) func(context.Context) error {
		return d.StartConfigWatcher
	}),

	// Provide a function to set the config watcher
	fx.Provide(func(d *Domain) func(ports.ConfigWatcher) {
		return d.SetConfigWatcher
	}),

	// Provide a function to set the client
	fx.Provide(func(d *Domain) func(ports.ClientInterface) {
		return d.SetClient
	}),

	// Provide domain service methods
	fx.Provide(func(d *Domain) func() error {
		return d.StartExileClientConfiguration
	}),

	fx.Provide(func(d *Domain) func() error {
		return d.StartPollEvents
	}),

	fx.Provide(func(d *Domain) func() error {
		return d.StartStreamJobs
	}),

	fx.Provide(func(d *Domain) func() error {
		return d.StartHostPlugin
	}),

	fx.Provide(func(d *Domain) func() error {
		return d.StopAllProcesses
	}),

	fx.Provide(func(d *Domain) func() bool {
		return d.IsRunning
	}),
)

// DomainService defines the interface for domain services.
type DomainService interface {
	StartConfigWatcher(ctx context.Context) error
	StartExileClientConfiguration() error
	StartPollEvents() error
	StartStreamJobs() error
	StartHostPlugin() error
	StopAllProcesses() error
	IsRunning() bool
}

// Ensure Domain implements DomainService interface.
var _ DomainService = (*Domain)(nil)

// DomainServiceProvider provides a way to get domain services.
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
