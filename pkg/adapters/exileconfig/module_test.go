package exileconfig

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/tcncloud/sati-go/pkg/ports"
	"go.uber.org/fx"
)

func TestModule(t *testing.T) {
	// Create a test app with the module
	app := fx.New(
		Module,
		fx.Provide(func() *zerolog.Logger {
			logger := zerolog.Nop()
			return &logger
		}),
		fx.Provide(func() ports.ClientInterface {
			return &MockClientInterface{}
		}),
		fx.Provide(func() ports.DomainService {
			return &MockDomainService{}
		}),
		fx.Invoke(func(watcher ports.ConfigWatcher) {
			// Test that the watcher is provided
			if watcher == nil {
				t.Error("Expected ConfigWatcher to be provided")
			}
		}),
	)

	// Test that the app starts without errors
	err := app.Err()
	if err != nil {
		t.Fatalf("Module failed to initialize: %v", err)
	}

	// Start the app
	ctx := context.Background()
	err = app.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start app: %v", err)
	}

	// Stop the app
	err = app.Stop(ctx)
	if err != nil {
		t.Fatalf("Failed to stop app: %v", err)
	}
}

func TestExileConfigServiceProvider_NewExileConfigServiceProvider(t *testing.T) {
	logger := zerolog.Nop()
	mockClient := &MockClientInterface{}
	mockDomainService := &MockDomainService{}

	exileConfig := NewExileConfig(
		mockClient,
		mockDomainService,
		&logger,
	)

	provider := NewExileConfigServiceProvider(exileConfig)

	if provider == nil {
		t.Fatal("Expected provider to be created")
	}

	if provider.exileConfig != exileConfig {
		t.Error("Expected exile config to be set")
	}
}

func TestExileConfigServiceProvider_GetConfigWatcherStarter(t *testing.T) {
	logger := zerolog.Nop()
	mockClient := &MockClientInterface{}
	mockDomainService := &MockDomainService{}

	exileConfig := NewExileConfig(
		mockClient,
		mockDomainService,
		&logger,
	)

	provider := NewExileConfigServiceProvider(exileConfig)

	starter := provider.GetConfigWatcherStarter()

	if starter == nil {
		t.Fatal("Expected starter function to be returned")
	}

	// Test that the function works
	ctx := context.Background()
	err := starter(ctx)

	// Should not panic and should return nil
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestExileConfigServiceProvider_GetConfigWatcherStopper(t *testing.T) {
	logger := zerolog.Nop()
	mockClient := &MockClientInterface{}
	mockDomainService := &MockDomainService{}

	exileConfig := NewExileConfig(
		mockClient,
		mockDomainService,
		&logger,
	)

	provider := NewExileConfigServiceProvider(exileConfig)

	stopper := provider.GetConfigWatcherStopper()

	if stopper == nil {
		t.Fatal("Expected stopper function to be returned")
	}

	// Test that the function works
	err := stopper()

	// Should not panic and should return nil
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestExileConfigServiceProvider_GetIsWatchingChecker(t *testing.T) {
	logger := zerolog.Nop()
	mockClient := &MockClientInterface{}
	mockDomainService := &MockDomainService{}

	exileConfig := NewExileConfig(
		mockClient,
		mockDomainService,
		&logger,
	)

	provider := NewExileConfigServiceProvider(exileConfig)

	checker := provider.GetIsWatchingChecker()

	if checker == nil {
		t.Fatal("Expected checker function to be returned")
	}

	// Test that the function works
	watching := checker()

	// Should not panic and should return false initially
	if watching {
		t.Error("Expected exile config to not be watching initially")
	}
}

func TestExileConfigServiceProvider_GetLastConfigurationGetter(t *testing.T) {
	logger := zerolog.Nop()
	mockClient := &MockClientInterface{}
	mockDomainService := &MockDomainService{}

	exileConfig := NewExileConfig(
		mockClient,
		mockDomainService,
		&logger,
	)

	provider := NewExileConfigServiceProvider(exileConfig)

	getter := provider.GetLastConfigurationGetter()

	if getter == nil {
		t.Fatal("Expected getter function to be returned")
	}

	// Test that the function works
	config := getter()

	// Should not panic and should return nil initially
	if config != nil {
		t.Error("Expected last configuration to be nil initially")
	}
}

func TestExileConfig_ImplementsConfigChangeHandler(t *testing.T) {
	logger := zerolog.Nop()
	mockClient := &MockClientInterface{}
	mockDomainService := &MockDomainService{}

	exileConfig := NewExileConfig(
		mockClient,
		mockDomainService,
		&logger,
	)

	// Test that ExileConfig implements ports.ConfigChangeHandler interface
	var _ ports.ConfigChangeHandler = exileConfig
}
