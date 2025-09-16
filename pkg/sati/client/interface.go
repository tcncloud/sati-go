package client

import "context"

// ClientInterface defines the interface for the TCN Exile Gate service client.
// This interface abstracts the concrete client implementation and allows for
// easier testing and dependency injection following clean architecture principles.
type ClientInterface interface {
	// Close closes the client connection.
	Close() error

	// Agent Management
	AddAgentCallResponse(ctx context.Context, params AddAgentCallResponseParams) (AddAgentCallResponseResult, error)
	GetAgentByID(ctx context.Context, params GetAgentByIDParams) (GetAgentByIDResult, error)
	GetAgentByPartnerID(ctx context.Context, params GetAgentByPartnerIDParams) (GetAgentByPartnerIDResult, error)
	GetAgentStatus(ctx context.Context, params GetAgentStatusParams) (GetAgentStatusResult, error)
	UpdateAgentStatus(ctx context.Context, params UpdateAgentStatusParams) (UpdateAgentStatusResult, error)
	UpsertAgent(ctx context.Context, params UpsertAgentParams) (UpsertAgentResult, error)

	// Configuration
	GetClientConfiguration(ctx context.Context, params GetClientConfigurationParams) (GetClientConfigurationResult, error)
	GetOrganizationInfo(ctx context.Context, params GetOrganizationInfoParams) (GetOrganizationInfoResult, error)
	RotateCertificate(ctx context.Context, params RotateCertificateParams) (RotateCertificateResult, error)

	// Events and Jobs
	PollEvents(ctx context.Context, params PollEventsParams) (PollEventsResult, error)
	StreamJobs(ctx context.Context, params StreamJobsParams) <-chan StreamJobsResult
	SubmitJobResults(ctx context.Context, params SubmitJobResultsParams) (SubmitJobResultsResult, error)

	// Call Management
	Dial(ctx context.Context, params DialParams) (DialResult, error)
	PutCallOnSimpleHold(ctx context.Context, params PutCallOnSimpleHoldParams) (PutCallOnSimpleHoldResult, error)
	TakeCallOffSimpleHold(ctx context.Context, params TakeCallOffSimpleHoldParams) (TakeCallOffSimpleHoldResult, error)
	Transfer(ctx context.Context, params TransferParams) (TransferResult, error)

	// Recording
	StartCallRecording(ctx context.Context, params StartCallRecordingParams) (StartCallRecordingResult, error)
	StopCallRecording(ctx context.Context, params StopCallRecordingParams) (StopCallRecordingResult, error)
	GetRecordingStatus(ctx context.Context, params GetRecordingStatusParams) (GetRecordingStatusResult, error)

	// Lists and Data
	ListAgents(ctx context.Context, params ListAgentsParams) <-chan ListAgentsResult
	ListHuntGroupPauseCodes(ctx context.Context, params ListHuntGroupPauseCodesParams) (ListHuntGroupPauseCodesResult, error)
	ListNCLRulesetNames(ctx context.Context, params ListNCLRulesetNamesParams) (ListNCLRulesetNamesResult, error)
	ListScrubLists(ctx context.Context, params ListScrubListsParams) (ListScrubListsResult, error)
	ListSkills(ctx context.Context, params ListSkillsParams) (ListSkillsResult, error)
	ListAgentSkills(ctx context.Context, params ListAgentSkillsParams) (ListAgentSkillsResult, error)

	// Scrub List Management
	AddScrubListEntries(ctx context.Context, params AddScrubListEntriesParams) (AddScrubListEntriesResult, error)
	UpdateScrubListEntry(ctx context.Context, params UpdateScrubListEntryParams) (UpdateScrubListEntryResult, error)
	RemoveScrubListEntries(ctx context.Context, params RemoveScrubListEntriesParams) (RemoveScrubListEntriesResult, error)

	// Skills Management
	AssignAgentSkill(ctx context.Context, params AssignAgentSkillParams) (AssignAgentSkillResult, error)
	UnassignAgentSkill(ctx context.Context, params UnassignAgentSkillParams) (UnassignAgentSkillResult, error)

	// Voice Recordings
	SearchVoiceRecordings(ctx context.Context, params SearchVoiceRecordingsParams) <-chan SearchVoiceRecordingsResult
	GetVoiceRecordingDownloadLink(ctx context.Context, params GetVoiceRecordingDownloadLinkParams) (GetVoiceRecordingDownloadLinkResult, error)
	ListSearchableRecordingFields(ctx context.Context, params ListSearchableRecordingFieldsParams) (ListSearchableRecordingFieldsResult, error)

	// Logging
	Log(ctx context.Context, params LogParams) (LogResult, error)
}
