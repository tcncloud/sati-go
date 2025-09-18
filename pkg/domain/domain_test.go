package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/tcncloud/sati-go/pkg/ports"
)

// MockConfigWatcher is a mock implementation of ports.ConfigWatcher
type MockConfigWatcher struct {
	startError  error
	stopError   error
	watching    bool
	startCalled bool
	stopCalled  bool
}

func (m *MockConfigWatcher) Start(ctx context.Context) error {
	m.startCalled = true
	return m.startError
}

func (m *MockConfigWatcher) Stop() error {
	m.stopCalled = true
	return m.stopError
}

func (m *MockConfigWatcher) IsWatching() bool {
	return m.watching
}

// MockClientInterface is a mock implementation of ports.ClientInterface
type MockClientInterface struct {
	getClientConfigResult ports.GetClientConfigurationResult
	getClientConfigError  error
	pollEventsResult      ports.PollEventsResult
	pollEventsError       error
	streamJobsChan        <-chan ports.StreamJobsResult
	closeError            error
}

func (m *MockClientInterface) Close() error {
	return m.closeError
}

func (m *MockClientInterface) GetClientConfiguration(ctx context.Context, params ports.GetClientConfigurationParams) (ports.GetClientConfigurationResult, error) {
	return m.getClientConfigResult, m.getClientConfigError
}

func (m *MockClientInterface) PollEvents(ctx context.Context, params ports.PollEventsParams) (ports.PollEventsResult, error) {
	return m.pollEventsResult, m.pollEventsError
}

func (m *MockClientInterface) StreamJobs(ctx context.Context, params ports.StreamJobsParams) <-chan ports.StreamJobsResult {
	return m.streamJobsChan
}

// Implement all other required methods with default implementations
func (m *MockClientInterface) AddAgentCallResponse(ctx context.Context, params ports.AddAgentCallResponseParams) (ports.AddAgentCallResponseResult, error) {
	return ports.AddAgentCallResponseResult{}, nil
}

func (m *MockClientInterface) GetAgentByID(ctx context.Context, params ports.GetAgentByIDParams) (ports.GetAgentByIDResult, error) {
	return ports.GetAgentByIDResult{}, nil
}

func (m *MockClientInterface) GetAgentByPartnerID(ctx context.Context, params ports.GetAgentByPartnerIDParams) (ports.GetAgentByPartnerIDResult, error) {
	return ports.GetAgentByPartnerIDResult{}, nil
}

func (m *MockClientInterface) GetAgentStatus(ctx context.Context, params ports.GetAgentStatusParams) (ports.GetAgentStatusResult, error) {
	return ports.GetAgentStatusResult{}, nil
}

func (m *MockClientInterface) UpdateAgentStatus(ctx context.Context, params ports.UpdateAgentStatusParams) (ports.UpdateAgentStatusResult, error) {
	return ports.UpdateAgentStatusResult{}, nil
}

func (m *MockClientInterface) UpsertAgent(ctx context.Context, params ports.UpsertAgentParams) (ports.UpsertAgentResult, error) {
	return ports.UpsertAgentResult{}, nil
}

func (m *MockClientInterface) GetOrganizationInfo(ctx context.Context, params ports.GetOrganizationInfoParams) (ports.GetOrganizationInfoResult, error) {
	return ports.GetOrganizationInfoResult{}, nil
}

func (m *MockClientInterface) RotateCertificate(ctx context.Context, params ports.RotateCertificateParams) (ports.RotateCertificateResult, error) {
	return ports.RotateCertificateResult{}, nil
}

func (m *MockClientInterface) SubmitJobResults(ctx context.Context, params ports.SubmitJobResultsParams) (ports.SubmitJobResultsResult, error) {
	return ports.SubmitJobResultsResult{}, nil
}

func (m *MockClientInterface) Dial(ctx context.Context, params ports.DialParams) (ports.DialResult, error) {
	return ports.DialResult{}, nil
}

func (m *MockClientInterface) PutCallOnSimpleHold(ctx context.Context, params ports.PutCallOnSimpleHoldParams) (ports.PutCallOnSimpleHoldResult, error) {
	return ports.PutCallOnSimpleHoldResult{}, nil
}

func (m *MockClientInterface) TakeCallOffSimpleHold(ctx context.Context, params ports.TakeCallOffSimpleHoldParams) (ports.TakeCallOffSimpleHoldResult, error) {
	return ports.TakeCallOffSimpleHoldResult{}, nil
}

func (m *MockClientInterface) Transfer(ctx context.Context, params ports.TransferParams) (ports.TransferResult, error) {
	return ports.TransferResult{}, nil
}

func (m *MockClientInterface) StartCallRecording(ctx context.Context, params ports.StartCallRecordingParams) (ports.StartCallRecordingResult, error) {
	return ports.StartCallRecordingResult{}, nil
}

func (m *MockClientInterface) StopCallRecording(ctx context.Context, params ports.StopCallRecordingParams) (ports.StopCallRecordingResult, error) {
	return ports.StopCallRecordingResult{}, nil
}

func (m *MockClientInterface) GetRecordingStatus(ctx context.Context, params ports.GetRecordingStatusParams) (ports.GetRecordingStatusResult, error) {
	return ports.GetRecordingStatusResult{}, nil
}

func (m *MockClientInterface) ListAgents(ctx context.Context, params ports.ListAgentsParams) <-chan ports.ListAgentsResult {
	return nil
}

func (m *MockClientInterface) ListHuntGroupPauseCodes(ctx context.Context, params ports.ListHuntGroupPauseCodesParams) (ports.ListHuntGroupPauseCodesResult, error) {
	return ports.ListHuntGroupPauseCodesResult{}, nil
}

func (m *MockClientInterface) ListNCLRulesetNames(ctx context.Context, params ports.ListNCLRulesetNamesParams) (ports.ListNCLRulesetNamesResult, error) {
	return ports.ListNCLRulesetNamesResult{}, nil
}

func (m *MockClientInterface) ListScrubLists(ctx context.Context, params ports.ListScrubListsParams) (ports.ListScrubListsResult, error) {
	return ports.ListScrubListsResult{}, nil
}

func (m *MockClientInterface) ListSkills(ctx context.Context, params ports.ListSkillsParams) (ports.ListSkillsResult, error) {
	return ports.ListSkillsResult{}, nil
}

func (m *MockClientInterface) ListAgentSkills(ctx context.Context, params ports.ListAgentSkillsParams) (ports.ListAgentSkillsResult, error) {
	return ports.ListAgentSkillsResult{}, nil
}

func (m *MockClientInterface) AddScrubListEntries(ctx context.Context, params ports.AddScrubListEntriesParams) (ports.AddScrubListEntriesResult, error) {
	return ports.AddScrubListEntriesResult{}, nil
}

func (m *MockClientInterface) RemoveScrubListEntries(ctx context.Context, params ports.RemoveScrubListEntriesParams) (ports.RemoveScrubListEntriesResult, error) {
	return ports.RemoveScrubListEntriesResult{}, nil
}

func (m *MockClientInterface) UpdateScrubListEntry(ctx context.Context, params ports.UpdateScrubListEntryParams) (ports.UpdateScrubListEntryResult, error) {
	return ports.UpdateScrubListEntryResult{}, nil
}

func (m *MockClientInterface) AssignAgentSkill(ctx context.Context, params ports.AssignAgentSkillParams) (ports.AssignAgentSkillResult, error) {
	return ports.AssignAgentSkillResult{}, nil
}

func (m *MockClientInterface) UnassignAgentSkill(ctx context.Context, params ports.UnassignAgentSkillParams) (ports.UnassignAgentSkillResult, error) {
	return ports.UnassignAgentSkillResult{}, nil
}

func (m *MockClientInterface) GetVoiceRecordingDownloadLink(ctx context.Context, params ports.GetVoiceRecordingDownloadLinkParams) (ports.GetVoiceRecordingDownloadLinkResult, error) {
	return ports.GetVoiceRecordingDownloadLinkResult{}, nil
}

func (m *MockClientInterface) ListSearchableRecordingFields(ctx context.Context, params ports.ListSearchableRecordingFieldsParams) (ports.ListSearchableRecordingFieldsResult, error) {
	return ports.ListSearchableRecordingFieldsResult{}, nil
}

func (m *MockClientInterface) SearchVoiceRecordings(ctx context.Context, params ports.SearchVoiceRecordingsParams) <-chan ports.SearchVoiceRecordingsResult {
	return nil
}

func (m *MockClientInterface) Log(ctx context.Context, params ports.LogParams) (ports.LogResult, error) {
	return ports.LogResult{}, nil
}

// setupTestDomain creates a test domain with mocked dependencies
func setupTestDomain() (*Domain, *MockConfigWatcher, *MockClientInterface) {
	logger := zerolog.Nop()
	mockWatcher := &MockConfigWatcher{}
	mockClient := &MockClientInterface{}

	domain := NewDomain(&logger)
	domain.SetConfigWatcher(mockWatcher)
	domain.SetClient(mockClient)

	return domain, mockWatcher, mockClient
}

func TestDomain_NewDomain(t *testing.T) {
	logger := zerolog.Nop()
	domain := NewDomain(&logger)

	if domain == nil {
		t.Fatal("Expected domain to be created")
	}

	if domain.log != &logger {
		t.Error("Expected logger to be set")
	}

	if domain.shutdownChan == nil {
		t.Error("Expected shutdown channel to be created")
	}

	if domain.IsRunning() {
		t.Error("Expected domain to not be running initially")
	}
}

func TestDomain_SetConfigWatcher(t *testing.T) {
	domain, mockWatcher, _ := setupTestDomain()

	if domain.configWatcher != mockWatcher {
		t.Error("Expected config watcher to be set")
	}
}

func TestDomain_SetClient(t *testing.T) {
	domain, _, _ := setupTestDomain()

	// Test that the client was set (we can't directly compare interfaces)
	if domain.client == nil {
		t.Error("Expected client to be set")
	}
}

func TestDomain_StartConfigWatcher(t *testing.T) {
	domain, mockWatcher, _ := setupTestDomain()
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockWatcher.startError = nil
		mockWatcher.startCalled = false

		err := domain.StartConfigWatcher(ctx)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if !mockWatcher.startCalled {
			t.Error("Expected config watcher start to be called")
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockWatcher.startError = errors.New("start failed")
		mockWatcher.startCalled = false

		err := domain.StartConfigWatcher(ctx)

		if err == nil {
			t.Error("Expected error but got none")
		}

		if !mockWatcher.startCalled {
			t.Error("Expected config watcher start to be called")
		}
	})

	t.Run("NoConfigWatcher", func(t *testing.T) {
		domainNoWatcher := NewDomain(domain.log)
		domainNoWatcher.SetClient(domain.client)

		err := domainNoWatcher.StartConfigWatcher(ctx)

		if err != nil {
			t.Errorf("Expected no error when no config watcher, got: %v", err)
		}
	})
}

func TestDomain_StartExileClientConfiguration(t *testing.T) {
	domain, _, mockClient := setupTestDomain()

	t.Run("Success", func(t *testing.T) {
		mockClient.getClientConfigError = nil

		err := domain.StartExileClientConfiguration()

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if domain.exileConfigProcess == nil {
			t.Error("Expected exile config process to be created")
		}
	})

	t.Run("NoClient", func(t *testing.T) {
		domainNoClient := NewDomain(domain.log)
		domainNoClient.SetConfigWatcher(domain.configWatcher)

		err := domainNoClient.StartExileClientConfiguration()

		if err != nil {
			t.Errorf("Expected no error when no client, got: %v", err)
		}
	})

	t.Run("AlreadyRunning", func(t *testing.T) {
		// Start it once
		_ = domain.StartExileClientConfiguration()

		// Try to start again
		err := domain.StartExileClientConfiguration()

		if err != nil {
			t.Errorf("Expected no error when already running, got: %v", err)
		}
	})
}

func TestDomain_StartPollEvents(t *testing.T) {
	domain, _, mockClient := setupTestDomain()

	t.Run("Success", func(t *testing.T) {
		mockClient.pollEventsError = nil

		err := domain.StartPollEvents()

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if domain.pollEventsProcess == nil {
			t.Error("Expected poll events process to be created")
		}
	})

	t.Run("NoClient", func(t *testing.T) {
		domainNoClient := NewDomain(domain.log)
		domainNoClient.SetConfigWatcher(domain.configWatcher)

		err := domainNoClient.StartPollEvents()

		if err != nil {
			t.Errorf("Expected no error when no client, got: %v", err)
		}
	})

	t.Run("AlreadyRunning", func(t *testing.T) {
		// Start it once
		_ = domain.StartPollEvents()

		// Try to start again
		err := domain.StartPollEvents()

		if err != nil {
			t.Errorf("Expected no error when already running, got: %v", err)
		}
	})
}

func TestDomain_StartStreamJobs(t *testing.T) {
	domain, _, _ := setupTestDomain()

	t.Run("Success", func(t *testing.T) {
		err := domain.StartStreamJobs()

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if domain.streamJobsProcess == nil {
			t.Error("Expected stream jobs process to be created")
		}
	})

	t.Run("NoClient", func(t *testing.T) {
		domainNoClient := NewDomain(domain.log)
		domainNoClient.SetConfigWatcher(domain.configWatcher)

		err := domainNoClient.StartStreamJobs()

		if err != nil {
			t.Errorf("Expected no error when no client, got: %v", err)
		}
	})

	t.Run("AlreadyRunning", func(t *testing.T) {
		// Start it once
		_ = domain.StartStreamJobs()

		// Try to start again
		err := domain.StartStreamJobs()

		if err != nil {
			t.Errorf("Expected no error when already running, got: %v", err)
		}
	})
}

func TestDomain_StartHostPlugin(t *testing.T) {
	domain, _, _ := setupTestDomain()

	t.Run("Success", func(t *testing.T) {
		err := domain.StartHostPlugin()

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if domain.hostPluginProcess == nil {
			t.Error("Expected host plugin process to be created")
		}
	})

	t.Run("AlreadyRunning", func(t *testing.T) {
		// Start it once
		_ = domain.StartHostPlugin()

		// Try to start again
		err := domain.StartHostPlugin()

		if err != nil {
			t.Errorf("Expected no error when already running, got: %v", err)
		}
	})
}

func TestDomain_StopAllProcesses(t *testing.T) {
	domain, mockWatcher, _ := setupTestDomain()

	// Start some processes
	_ = domain.StartHostPlugin()
	_ = domain.StartPollEvents()
	_ = domain.StartStreamJobs()

	t.Run("Success", func(t *testing.T) {
		mockWatcher.stopError = nil
		mockWatcher.stopCalled = false

		err := domain.StopAllProcesses()

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if !mockWatcher.stopCalled {
			t.Error("Expected config watcher stop to be called")
		}

		if domain.hostPluginProcess != nil {
			t.Error("Expected host plugin process to be stopped")
		}

		if domain.pollEventsProcess != nil {
			t.Error("Expected poll events process to be stopped")
		}

		if domain.streamJobsProcess != nil {
			t.Error("Expected stream jobs process to be stopped")
		}

		if domain.IsRunning() {
			t.Error("Expected domain to not be running")
		}
	})

	t.Run("ConfigWatcherStopError", func(t *testing.T) {
		// Reset domain
		domain, mockWatcher, _ = setupTestDomain()
		mockWatcher.stopError = errors.New("stop failed")
		mockWatcher.stopCalled = false

		err := domain.StopAllProcesses()

		if err == nil {
			t.Error("Expected error but got none")
		}

		if !mockWatcher.stopCalled {
			t.Error("Expected config watcher stop to be called")
		}
	})
}

func TestDomain_IsRunning(t *testing.T) {
	domain, _, _ := setupTestDomain()

	if domain.IsRunning() {
		t.Error("Expected domain to not be running initially")
	}

	// Start a process
	_ = domain.StartHostPlugin()

	// Domain should still not be running (isRunning flag is only set in StopAllProcesses)
	if domain.IsRunning() {
		t.Error("Expected domain to not be running")
	}
}

func TestDomain_ImplementsDomainService(t *testing.T) {
	domain, _, _ := setupTestDomain()

	// Test that Domain implements DomainService interface
	var _ ports.DomainService = domain
}
