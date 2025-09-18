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

	// Provide the domain service interface
	fx.Provide(func(domain *Domain) ports.DomainService {
		return domain
	}),

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

// Ensure Domain implements DomainService interface.
var _ ports.DomainService = (*Domain)(nil)
