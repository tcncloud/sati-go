package hostplugin

import (
	"go.uber.org/fx"

	"github.com/rs/zerolog"
	"github.com/tcncloud/sati-go/pkg/ports"
)

// Module provides the hostplugin module for dependency injection.
// It includes the HostPluginProcess implementation for hosting plugins and managing their lifecycle.
//
// Usage example:
//
//	app := fx.New(
//	  hostplugin.Module,
//	  fx.Invoke(func(process ports.HostPluginProcess) {
//	    process.Run(context.Background())
//	  }),
//	)
var Module = fx.Module("hostplugin",
	// Provide the concrete HostPluginProcess implementation
	fx.Provide(NewHostPluginProcess),

	// Provide the interface implementation
	fx.Provide(func(process *HostPluginProcess) ports.HostPluginProcess {
		return process
	}),
)

// NewHostPluginProcessWithLogger creates a new HostPluginProcess with a specific logger.
// This is useful for testing or when you need a specific logger instance.
func NewHostPluginProcessWithLogger(log *zerolog.Logger) *HostPluginProcess {
	return NewHostPluginProcess(log)
}

