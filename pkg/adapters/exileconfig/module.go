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

package exileconfig

import (
	"context"

	"github.com/tcncloud/sati-go/pkg/ports"
	"go.uber.org/fx"
)

// Module provides the exile config adapter module for dependency injection.
// It includes the ExileConfig service that watches for configuration changes
// and notifies the domain service when changes occur.
//
// Usage example:
//
//	app := fx.New(
//	  exileconfig.Module,
//	  fx.Invoke(func(watcher ports.ConfigWatcher) {
//	    ctx := context.Background()
//	    watcher.Start(ctx)
//	  }),
//	)
var Module = fx.Module("exileconfig",
	// Provide the ExileConfig service
	fx.Provide(NewExileConfig),

	// Provide the ConfigWatcher interface
	fx.Provide(func(exileConfig *ExileConfig) ports.ConfigWatcher {
		return exileConfig
	}),

	// Provide service methods as injectable functions
	fx.Provide(func(ec *ExileConfig) func(context.Context) error {
		return ec.Start
	}),

	fx.Provide(func(ec *ExileConfig) func() error {
		return ec.Stop
	}),

	fx.Provide(func(ec *ExileConfig) func() bool {
		return ec.IsWatching
	}),

	fx.Provide(func(ec *ExileConfig) func() *ports.GetClientConfigurationResult {
		return ec.GetLastConfiguration
	}),
)

// Ensure ExileConfig implements required interfaces.
var _ ports.ConfigWatcher = (*ExileConfig)(nil)
var _ ports.ConfigChangeHandler = (*ExileConfig)(nil)

// ExileConfigServiceProvider provides clean access to exile config service methods.
type ExileConfigServiceProvider struct {
	exileConfig *ExileConfig
}

// NewExileConfigServiceProvider creates a new exile config service provider.
func NewExileConfigServiceProvider(exileConfig *ExileConfig) *ExileConfigServiceProvider {
	return &ExileConfigServiceProvider{
		exileConfig: exileConfig,
	}
}

// GetConfigWatcherStarter returns the config watcher starter function.
func (esp *ExileConfigServiceProvider) GetConfigWatcherStarter() func(context.Context) error {
	return esp.exileConfig.Start
}

// GetConfigWatcherStopper returns the config watcher stopper function.
func (esp *ExileConfigServiceProvider) GetConfigWatcherStopper() func() error {
	return esp.exileConfig.Stop
}

// GetIsWatchingChecker returns the is watching checker function.
func (esp *ExileConfigServiceProvider) GetIsWatchingChecker() func() bool {
	return esp.exileConfig.IsWatching
}

// GetLastConfigurationGetter returns the last configuration getter function.
func (esp *ExileConfigServiceProvider) GetLastConfigurationGetter() func() *ports.GetClientConfigurationResult {
	return esp.exileConfig.GetLastConfiguration
}
