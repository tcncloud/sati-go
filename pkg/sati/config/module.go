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

package config

import (
	"github.com/tcncloud/sati-go/pkg/interfaces"
	"go.uber.org/fx"
)

// Module provides the configuration module for dependency injection.
// It includes the ConfigWatcher implementation and related services.
//
// Usage example:
//
//	app := fx.New(
//	  config.Module,
//	  fx.Invoke(func(watcherFactory func([]string, ConfigLoaderFunc) (interfaces.ConfigWatcher, error)) {
//	    watcher, err := watcherFactory([]string{"/path/to/config.cfg"}, loaderFunc)
//	    if err != nil {
//	      log.Fatal(err)
//	    }
//	    watcher.Start(context.Background())
//	  }),
//	)
var Module = fx.Module("config",
	// Provide the ConfigLoader as an implementation of interfaces.ConfigLoader
	fx.Provide(func() interfaces.ConfigLoader {
		return &configLoader{}
	}),

	// Provide utility functions for configuration management
	fx.Provide(LoadConfig),
	fx.Provide(NewConfigFromString),
	fx.Provide(LoadAndValidateConfig),
	fx.Provide(LoadAndValidateConfigFromString),

	// Provide a factory for creating ConfigWatcher instances
	fx.Provide(func() func([]string, ConfigLoaderFunc) (interfaces.ConfigWatcher, error) {
		return func(configPaths []string, loader ConfigLoaderFunc) (interfaces.ConfigWatcher, error) {
			return NewConfigWatcher(configPaths, loader)
		}
	}),
)

// configLoader implements the interfaces.ConfigLoader interface
type configLoader struct{}

// LoadConfig loads configuration from a file path.
func (cl *configLoader) LoadConfig(path string) (interface{}, error) {
	return LoadConfig(path)
}

// LoadConfigFromString loads configuration from a base64 encoded string.
func (cl *configLoader) LoadConfigFromString(configString string) (interface{}, error) {
	return NewConfigFromString(configString)
}
