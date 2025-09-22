package ports

import "context"

// DomainService defines the interface for domain services.
type DomainService interface {
	// StartConfigWatcher starts the configuration watcher.
	StartConfigWatcher(ctx context.Context) error

	// StartExileClientConfiguration starts the exile client configuration process.
	StartExileClientConfiguration() error

	// StartPollEvents starts the poll events process.
	StartPollEvents() error

	// StartStreamJobs starts the stream jobs process.
	StartStreamJobs() error

	// StartHostPlugin starts the host plugin process.
	StartHostPlugin() error

	// StopAllProcesses stops all running processes.
	StopAllProcesses() error

	// IsRunning returns true if the domain is currently running.
	IsRunning() bool

	// ClientConfigurationChanged is called when the client configuration changes.
	ClientConfigurationChanged(oldConfig, newConfig *GetClientConfigurationResult) error
}

// ProcessManager defines the interface for managing domain processes.
type ProcessManager interface {
	// StartProcess starts a specific process by name.
	StartProcess(name string) error

	// StopProcess stops a specific process by name.
	StopProcess(name string) error

	// StopAllProcesses stops all running processes.
	StopAllProcesses() error

	// IsProcessRunning returns true if the specified process is running.
	IsProcessRunning(name string) bool

	// GetProcessNames returns a list of all available process names.
	GetProcessNames() []string
}

// ConfigChangeHandler defines the interface for handling configuration changes.
type ConfigChangeHandler interface {
	// OnConfigChanged is called when the configuration changes.
	OnConfigChanged(oldConfig, newConfig *GetClientConfigurationResult)

	// ShouldRestartProcesses returns true if processes should be restarted.
	ShouldRestartProcesses(oldConfig, newConfig *GetClientConfigurationResult) bool
}

// EventDispatcher defines the interface for dispatching events to plugins.
type EventDispatcher interface {
	// DispatchEvents dispatches events to the plugin.
	DispatchEvents(events []Event)

	// DispatchJob dispatches a job to the plugin.
	DispatchJob(job *Job)
}

// HostPluginProcess defines the interface for hosting plugins and managing their lifecycle.
type HostPluginProcess interface {
	// Run starts the host plugin process.
	Run(ctx context.Context)

	// Stop stops the host plugin process.
	Stop()

	// DispatchEvents dispatches events to the plugin.
	DispatchEvents(events []Event)

	// DispatchJob dispatches a job to the plugin.
	DispatchJob(job *Job)
}
