package exileconfig

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/tcncloud/sati-go/pkg/ports"
)

// MockDomainService is a mock implementation of ports.DomainService
type MockDomainService struct {
	startExileConfigCalled           bool
	startExileConfigError            error
	clientConfigurationChangedCalled bool
	clientConfigurationChangedError  error
	lastOldConfig                    *ports.GetClientConfigurationResult
	lastNewConfig                    *ports.GetClientConfigurationResult
}

func (m *MockDomainService) StartConfigWatcher(ctx context.Context) error {
	return nil
}

func (m *MockDomainService) StartExileClientConfiguration() error {
	m.startExileConfigCalled = true
	return m.startExileConfigError
}

func (m *MockDomainService) StartPollEvents() error {
	return nil
}

func (m *MockDomainService) StartStreamJobs() error {
	return nil
}

func (m *MockDomainService) StartHostPlugin() error {
	return nil
}

func (m *MockDomainService) StopAllProcesses() error {
	return nil
}

func (m *MockDomainService) IsRunning() bool {
	return false
}

func (m *MockDomainService) ClientConfigurationChanged(oldConfig, newConfig *ports.GetClientConfigurationResult) error {
	m.clientConfigurationChangedCalled = true
	m.lastOldConfig = oldConfig
	m.lastNewConfig = newConfig
	return m.clientConfigurationChangedError
}

// MockConfigChangeHandler is a mock implementation of ports.ConfigChangeHandler
type MockConfigChangeHandler struct {
	onConfigChangedCalled bool
	oldConfig             *ports.GetClientConfigurationResult
	newConfig             *ports.GetClientConfigurationResult
	shouldRestartCalled   bool
	shouldRestartResult   bool
}

func (m *MockConfigChangeHandler) OnConfigChanged(oldConfig, newConfig *ports.GetClientConfigurationResult) {
	m.onConfigChangedCalled = true
	m.oldConfig = oldConfig
	m.newConfig = newConfig
}

func (m *MockConfigChangeHandler) ShouldRestartProcesses(oldConfig, newConfig *ports.GetClientConfigurationResult) bool {
	m.shouldRestartCalled = true
	return m.shouldRestartResult
}

// MockClientInterface is a mock implementation of ports.ClientInterface
type MockClientInterface struct {
	getClientConfigResult ports.GetClientConfigurationResult
	getClientConfigError  error
}

func (m *MockClientInterface) Close() error {
	return nil
}

func (m *MockClientInterface) GetClientConfiguration(ctx context.Context, params ports.GetClientConfigurationParams) (ports.GetClientConfigurationResult, error) {
	return m.getClientConfigResult, m.getClientConfigError
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

func (m *MockClientInterface) PollEvents(ctx context.Context, params ports.PollEventsParams) (ports.PollEventsResult, error) {
	return ports.PollEventsResult{}, nil
}

func (m *MockClientInterface) StreamJobs(ctx context.Context, params ports.StreamJobsParams) <-chan ports.StreamJobsResult {
	return nil
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

// setupTestExileConfig creates a test exile config with mocked dependencies
func setupTestExileConfig() (*ExileConfig, *MockClientInterface, *MockDomainService) {
	logger := zerolog.Nop()
	mockClient := &MockClientInterface{}
	mockDomainService := &MockDomainService{}

	exileConfig := NewExileConfig(
		mockClient,
		mockDomainService,
		&logger,
	)

	return exileConfig, mockClient, mockDomainService
}

func TestExileConfig_NewExileConfig(t *testing.T) {
	logger := zerolog.Nop()
	mockClient := &MockClientInterface{}
	mockDomainService := &MockDomainService{}

	exileConfig := NewExileConfig(
		mockClient,
		mockDomainService,
		&logger,
	)

	if exileConfig == nil {
		t.Fatal("Expected exile config to be created")
	}

	if exileConfig.exileClient != mockClient {
		t.Error("Expected exile client to be set")
	}

	if exileConfig.domainService != mockDomainService {
		t.Error("Expected domain service to be set")
	}

	if exileConfig.log != &logger {
		t.Error("Expected logger to be set")
	}

	if exileConfig.IsRunning() {
		t.Error("Expected exile config to not be running initially")
	}
}

func TestExileConfig_OnExileClientStarted(t *testing.T) {
	exileConfig, _, _ := setupTestExileConfig()
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		err := exileConfig.OnExileClientStarted(ctx)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if !exileConfig.IsRunning() {
			t.Error("Expected exile config to be running")
		}
	})

	t.Run("AlreadyRunning", func(t *testing.T) {
		// Start it once
		_ = exileConfig.OnExileClientStarted(ctx)

		// Try to start again
		err := exileConfig.OnExileClientStarted(ctx)

		if err != nil {
			t.Errorf("Expected no error when already running, got: %v", err)
		}
	})
}

func TestExileConfig_OnExileClientStopped(t *testing.T) {
	exileConfig, _, _ := setupTestExileConfig()
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Start first
		_ = exileConfig.OnExileClientStarted(ctx)

		err := exileConfig.OnExileClientStopped()

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if exileConfig.IsRunning() {
			t.Error("Expected exile config to not be running after stop")
		}
	})

	t.Run("NotRunning", func(t *testing.T) {
		// Try to stop when not running
		err := exileConfig.OnExileClientStopped()

		if err != nil {
			t.Errorf("Expected no error when not running, got: %v", err)
		}
	})
}

func TestExileConfig_IsRunning(t *testing.T) {
	exileConfig, _, _ := setupTestExileConfig()
	ctx := context.Background()

	if exileConfig.IsRunning() {
		t.Error("Expected exile config to not be running initially")
	}

	// Start running
	_ = exileConfig.OnExileClientStarted(ctx)

	if !exileConfig.IsRunning() {
		t.Error("Expected exile config to be running after start")
	}

	// Stop running
	_ = exileConfig.OnExileClientStopped()

	if exileConfig.IsRunning() {
		t.Error("Expected exile config to not be running after stop")
	}
}

func TestExileConfig_hasConfigurationChanged(t *testing.T) {
	exileConfig, _, _ := setupTestExileConfig()

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

	// Test with nil old config
	if !exileConfig.hasConfigurationChanged(nil, config1) {
		t.Error("Expected config change when old config is nil")
	}

	// Test with different configs
	if !exileConfig.hasConfigurationChanged(config1, config2) {
		t.Error("Expected config change when configs are different")
	}

	// Test with same configs
	if exileConfig.hasConfigurationChanged(config1, config3) {
		t.Error("Expected no config change when configs are the same")
	}

	// Test with nil new config
	if !exileConfig.hasConfigurationChanged(config1, nil) {
		t.Error("Expected config change when new config is nil")
	}
}

func TestExileConfig_GetLastConfiguration(t *testing.T) {
	exileConfig, _, _ := setupTestExileConfig()

	// Initially should be nil
	if exileConfig.GetLastConfiguration() != nil {
		t.Error("Expected last configuration to be nil initially")
	}

	// Set a configuration
	config := &ports.GetClientConfigurationResult{
		OrgID:         "test-org",
		OrgName:       "Test Org",
		ConfigName:    "test-config",
		ConfigPayload: "test-payload",
	}

	exileConfig.mu.Lock()
	exileConfig.lastConfig = config
	exileConfig.mu.Unlock()

	// Should return the set configuration
	result := exileConfig.GetLastConfiguration()
	if result == nil {
		t.Error("Expected last configuration to be set")
	}

	if result.OrgID != config.OrgID {
		t.Errorf("Expected OrgID %s, got %s", config.OrgID, result.OrgID)
	}
}

func TestExileConfig_ConfigurationChange(t *testing.T) {
	exileConfig, mockClient, mockDomainService := setupTestExileConfig()
	ctx := context.Background()

	// Set up initial configuration
	initialConfig := ports.GetClientConfigurationResult{
		OrgID:         "org1",
		OrgName:       "Org 1",
		ConfigName:    "config1",
		ConfigPayload: "payload1",
	}

	mockClient.getClientConfigResult = initialConfig
	mockClient.getClientConfigError = nil

	// Test configuration check directly
	err := exileConfig.checkConfiguration(ctx)
	if err != nil {
		t.Fatalf("Failed to check configuration: %v", err)
	}

	// Verify initial configuration was set
	lastConfig := exileConfig.GetLastConfiguration()
	if lastConfig == nil {
		t.Error("Expected last configuration to be set")
		return
	}

	if lastConfig.OrgID != initialConfig.OrgID {
		t.Errorf("Expected OrgID %s, got %s", initialConfig.OrgID, lastConfig.OrgID)
	}

	// Change configuration
	newConfig := ports.GetClientConfigurationResult{
		OrgID:         "org2",
		OrgName:       "Org 2",
		ConfigName:    "config2",
		ConfigPayload: "payload2",
	}

	mockClient.getClientConfigResult = newConfig

	// Test configuration check again
	err = exileConfig.checkConfiguration(ctx)
	if err != nil {
		t.Fatalf("Failed to check configuration: %v", err)
	}

	// Verify domain service was called
	if !mockDomainService.clientConfigurationChangedCalled {
		t.Error("Expected domain service ClientConfigurationChanged to be called")
	}

	// Verify the correct configs were passed
	if mockDomainService.lastOldConfig == nil {
		t.Error("Expected old config to be passed to domain service")
	}

	if mockDomainService.lastNewConfig == nil {
		t.Error("Expected new config to be passed to domain service")
	}

	if mockDomainService.lastOldConfig.OrgID != initialConfig.OrgID {
		t.Errorf("Expected old config OrgID %s, got %s", initialConfig.OrgID, mockDomainService.lastOldConfig.OrgID)
	}

	if mockDomainService.lastNewConfig.OrgID != newConfig.OrgID {
		t.Errorf("Expected new config OrgID %s, got %s", newConfig.OrgID, mockDomainService.lastNewConfig.OrgID)
	}
}

func TestExileConfig_ConfigurationError(t *testing.T) {
	exileConfig, mockClient, _ := setupTestExileConfig()
	ctx := context.Background()

	// Set up error condition
	mockClient.getClientConfigError = errors.New("configuration error")

	// Start running
	err := exileConfig.OnExileClientStarted(ctx)
	if err != nil {
		t.Fatalf("Failed to start exile config: %v", err)
	}

	// Wait for the check to complete
	time.Sleep(2 * time.Second)

	// The error should be logged but not cause the monitoring to stop
	if !exileConfig.IsRunning() {
		t.Error("Expected exile config to still be running after error")
	}

	// Stop running
	_ = exileConfig.OnExileClientStopped()
}

func TestExileConfig_ImplementsConfigWatcher(t *testing.T) {
	exileConfig, _, _ := setupTestExileConfig()

	// Test that ExileConfig implements ports.ConfigWatcher interface
	var _ ports.ConfigWatcher = exileConfig
}
