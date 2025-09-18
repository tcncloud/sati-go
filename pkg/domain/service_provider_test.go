package domain

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
)

func TestDomainServiceProvider_NewDomainServiceProvider(t *testing.T) {
	logger := zerolog.Nop()
	domain := NewDomain(&logger)

	provider := NewDomainServiceProvider(domain)

	if provider == nil {
		t.Fatal("Expected provider to be created")
	}

	if provider.domain != domain {
		t.Error("Expected domain to be set")
	}
}

func TestDomainServiceProvider_GetConfigWatcherStarter(t *testing.T) {
	logger := zerolog.Nop()
	domain := NewDomain(&logger)
	provider := NewDomainServiceProvider(domain)

	starter := provider.GetConfigWatcherStarter()

	if starter == nil {
		t.Fatal("Expected starter function to be returned")
	}

	// Test that the function works
	ctx := context.Background()
	err := starter(ctx)

	// Should not panic and should return nil (no config watcher set)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestDomainServiceProvider_GetExileClientConfigurationStarter(t *testing.T) {
	logger := zerolog.Nop()
	domain := NewDomain(&logger)
	provider := NewDomainServiceProvider(domain)

	starter := provider.GetExileClientConfigurationStarter()

	if starter == nil {
		t.Fatal("Expected starter function to be returned")
	}

	// Test that the function works
	err := starter()

	// Should not panic and should return nil (no client set)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestDomainServiceProvider_GetPollEventsStarter(t *testing.T) {
	logger := zerolog.Nop()
	domain := NewDomain(&logger)
	provider := NewDomainServiceProvider(domain)

	starter := provider.GetPollEventsStarter()

	if starter == nil {
		t.Fatal("Expected starter function to be returned")
	}

	// Test that the function works
	err := starter()

	// Should not panic and should return nil (no client set)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestDomainServiceProvider_GetStreamJobsStarter(t *testing.T) {
	logger := zerolog.Nop()
	domain := NewDomain(&logger)
	provider := NewDomainServiceProvider(domain)

	starter := provider.GetStreamJobsStarter()

	if starter == nil {
		t.Fatal("Expected starter function to be returned")
	}

	// Test that the function works
	err := starter()

	// Should not panic and should return nil (no client set)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestDomainServiceProvider_GetHostPluginStarter(t *testing.T) {
	logger := zerolog.Nop()
	domain := NewDomain(&logger)
	provider := NewDomainServiceProvider(domain)

	starter := provider.GetHostPluginStarter()

	if starter == nil {
		t.Fatal("Expected starter function to be returned")
	}

	// Test that the function works
	err := starter()

	// Should not panic and should return nil
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestDomainServiceProvider_GetStopAllProcessesStopper(t *testing.T) {
	logger := zerolog.Nop()
	domain := NewDomain(&logger)
	provider := NewDomainServiceProvider(domain)

	stopper := provider.GetStopAllProcessesStopper()

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

func TestDomainServiceProvider_GetIsRunningChecker(t *testing.T) {
	logger := zerolog.Nop()
	domain := NewDomain(&logger)
	provider := NewDomainServiceProvider(domain)

	checker := provider.GetIsRunningChecker()

	if checker == nil {
		t.Fatal("Expected checker function to be returned")
	}

	// Test that the function works
	running := checker()

	// Should not panic and should return false initially
	if running {
		t.Error("Expected domain to not be running initially")
	}
}
