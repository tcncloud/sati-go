package domain

import (
	"context"
	"testing"
	"time"

	"github.com/tcncloud/sati-go/pkg/ports"
)

func TestExileClientConfigurationProcess_hasConfigChanged(t *testing.T) {
	domain, _, _ := setupTestDomain()
	process := &ExileClientConfigurationProcess{
		domain: domain,
	}

	config1 := &ports.GetClientConfigurationResult{
		OrgID:         "org1",
		OrgName:       "Org 1",
		ConfigName:    "config1",
		ConfigPayload: "payload1",
	}

	config2 := &ports.GetClientConfigurationResult{
		OrgID:         "org2",
		OrgName:       "Org 2",
		ConfigName:    "config2",
		ConfigPayload: "payload2",
	}

	config3 := &ports.GetClientConfigurationResult{
		OrgID:         "org1",
		OrgName:       "Org 1",
		ConfigName:    "config1",
		ConfigPayload: "payload1",
	}

	// Test with different configs
	if !process.hasConfigChanged(config1, config2) {
		t.Error("Expected config change when configs are different")
	}

	// Test with same configs
	if process.hasConfigChanged(config1, config3) {
		t.Error("Expected no config change when configs are the same")
	}
}

func TestExileClientConfigurationProcess_checkConfiguration(t *testing.T) {
	domain, _, mockClient := setupTestDomain()
	process := &ExileClientConfigurationProcess{
		domain: domain,
	}

	expectedConfig := ports.GetClientConfigurationResult{
		OrgID:         "test-org",
		OrgName:       "Test Org",
		ConfigName:    "test-config",
		ConfigPayload: "test-payload",
	}

	mockClient.getClientConfigResult = expectedConfig
	mockClient.getClientConfigError = nil

	err := process.checkConfiguration()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if process.lastConfig == nil {
		t.Error("Expected last config to be set")
	}

	if process.lastConfig.OrgID != expectedConfig.OrgID {
		t.Errorf("Expected OrgID %s, got %s", expectedConfig.OrgID, process.lastConfig.OrgID)
	}
}

func TestExileClientConfigurationProcess_run(t *testing.T) {
	domain, _, mockClient := setupTestDomain()
	process := &ExileClientConfigurationProcess{
		domain: domain,
		ticker: time.NewTicker(10 * time.Millisecond), // Fast ticker for testing
	}

	expectedConfig := ports.GetClientConfigurationResult{
		OrgID:         "test-org",
		OrgName:       "Test Org",
		ConfigName:    "test-config",
		ConfigPayload: "test-payload",
	}

	mockClient.getClientConfigResult = expectedConfig
	mockClient.getClientConfigError = nil

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Run the process
	go process.run(ctx)

	// Wait for context to be done
	<-ctx.Done()

	// Verify that checkConfiguration was called
	if process.lastConfig == nil {
		t.Error("Expected last config to be set")
	}
}

func TestPollEventsProcess_pollEvents(t *testing.T) {
	domain, _, mockClient := setupTestDomain()
	process := &PollEventsProcess{
		domain: domain,
	}

	expectedEvents := []ports.Event{
		{Type: "test"},
		{Type: "test"},
	}

	expectedResult := ports.PollEventsResult{
		Events: expectedEvents,
	}

	mockClient.pollEventsResult = expectedResult
	mockClient.pollEventsError = nil

	err := process.pollEvents()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestPollEventsProcess_run(t *testing.T) {
	domain, _, mockClient := setupTestDomain()
	process := &PollEventsProcess{
		domain: domain,
	}

	expectedEvents := []ports.Event{
		{Type: "test"},
	}

	expectedResult := ports.PollEventsResult{
		Events: expectedEvents,
	}

	mockClient.pollEventsResult = expectedResult
	mockClient.pollEventsError = nil

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Run the process
	go process.run(ctx)

	// Wait for context to be done
	<-ctx.Done()
}

func TestStreamJobsProcess_streamJobs(t *testing.T) {
	domain, _, mockClient := setupTestDomain()
	process := &StreamJobsProcess{
		domain: domain,
	}

	// Create a test channel
	resultsChan := make(chan ports.StreamJobsResult, 1)
	resultsChan <- ports.StreamJobsResult{
		Job: &ports.Job{
			JobID: "job1",
			Type:  "test",
		},
	}
	close(resultsChan)

	mockClient.streamJobsChan = resultsChan

	err := process.streamJobs()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestStreamJobsProcess_run(t *testing.T) {
	domain, _, mockClient := setupTestDomain()
	process := &StreamJobsProcess{
		domain: domain,
	}

	// Create a test channel
	resultsChan := make(chan ports.StreamJobsResult, 1)
	resultsChan <- ports.StreamJobsResult{
		Job: &ports.Job{
			JobID: "job1",
			Type:  "test",
		},
	}
	close(resultsChan)

	mockClient.streamJobsChan = resultsChan

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Run the process
	go process.run(ctx)

	// Wait for context to be done
	<-ctx.Done()
}

// HostPluginProcess tests are now in the adapters package
// since HostPluginProcess is now an interface

func TestProcessStopMethods(t *testing.T) {
	domain, _, _ := setupTestDomain()

	// Test ExileClientConfigurationProcess stop
	exileProcess := &ExileClientConfigurationProcess{
		domain: domain,
		cancel: func() {}, // Mock cancel function
	}
	exileProcess.stop() // Should not panic

	// Test PollEventsProcess stop
	pollProcess := &PollEventsProcess{
		domain: domain,
		cancel: func() {}, // Mock cancel function
	}
	pollProcess.stop() // Should not panic

	// Test StreamJobsProcess stop
	streamProcess := &StreamJobsProcess{
		domain: domain,
		cancel: func() {}, // Mock cancel function
	}
	streamProcess.stop() // Should not panic

	// HostPluginProcess stop test removed since it's now an interface
}
