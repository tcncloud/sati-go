// Package ports defines the interfaces for the application's ports.
package ports

import (
	"context"

	"github.com/tcncloud/sati-go/pkg/sati/client"
)

// ClientInterface defines the interface for the TCN Exile Gate service client.
// This interface abstracts the concrete client implementation and allows for
// easier testing and dependency injection following clean architecture principles.
type ClientInterface interface {
	// Close closes the client connection.
	Close() error

	// Agent Management
	AddAgentCallResponse(ctx context.Context, params client.AddAgentCallResponseParams) (client.AddAgentCallResponseResult, error)
	GetAgentByID(ctx context.Context, params client.GetAgentByIDParams) (client.GetAgentByIDResult, error)
	GetAgentByPartnerID(ctx context.Context, params client.GetAgentByPartnerIDParams) (client.GetAgentByPartnerIDResult, error)
	GetAgentStatus(ctx context.Context, params client.GetAgentStatusParams) (client.GetAgentStatusResult, error)
	UpdateAgentStatus(ctx context.Context, params client.UpdateAgentStatusParams) (client.UpdateAgentStatusResult, error)
	UpsertAgent(ctx context.Context, params client.UpsertAgentParams) (client.UpsertAgentResult, error)

	// Configuration
	GetClientConfiguration(ctx context.Context, params client.GetClientConfigurationParams) (client.GetClientConfigurationResult, error)
	GetOrganizationInfo(ctx context.Context, params client.GetOrganizationInfoParams) (client.GetOrganizationInfoResult, error)
	RotateCertificate(ctx context.Context, params client.RotateCertificateParams) (client.RotateCertificateResult, error)

	// Events and Jobs
	PollEvents(ctx context.Context, params client.PollEventsParams) (client.PollEventsResult, error)
	StreamJobs(ctx context.Context, params client.StreamJobsParams) <-chan client.StreamJobsResult
	SubmitJobResults(ctx context.Context, params client.SubmitJobResultsParams) (client.SubmitJobResultsResult, error)

	// Call Management
	Dial(ctx context.Context, params client.DialParams) (client.DialResult, error)
	PutCallOnSimpleHold(ctx context.Context, params client.PutCallOnSimpleHoldParams) (client.PutCallOnSimpleHoldResult, error)
	TakeCallOffSimpleHold(ctx context.Context, params client.TakeCallOffSimpleHoldParams) (client.TakeCallOffSimpleHoldResult, error)
	Transfer(ctx context.Context, params client.TransferParams) (client.TransferResult, error)

	// Recording
	StartCallRecording(ctx context.Context, params client.StartCallRecordingParams) (client.StartCallRecordingResult, error)
	StopCallRecording(ctx context.Context, params client.StopCallRecordingParams) (client.StopCallRecordingResult, error)
	GetRecordingStatus(ctx context.Context, params client.GetRecordingStatusParams) (client.GetRecordingStatusResult, error)

	// Lists and Data
	ListAgents(ctx context.Context, params client.ListAgentsParams) <-chan client.ListAgentsResult
	ListHuntGroupPauseCodes(ctx context.Context, params client.ListHuntGroupPauseCodesParams) (client.ListHuntGroupPauseCodesResult, error)
	ListNCLRulesetNames(ctx context.Context, params client.ListNCLRulesetNamesParams) (client.ListNCLRulesetNamesResult, error)
	ListScrubLists(ctx context.Context, params client.ListScrubListsParams) (client.ListScrubListsResult, error)
	ListSkills(ctx context.Context, params client.ListSkillsParams) (client.ListSkillsResult, error)
	ListAgentSkills(ctx context.Context, params client.ListAgentSkillsParams) (client.ListAgentSkillsResult, error)

	// Scrub List Management
	AddScrubListEntries(ctx context.Context, params client.AddScrubListEntriesParams) (client.AddScrubListEntriesResult, error)
	UpdateScrubListEntry(ctx context.Context, params client.UpdateScrubListEntryParams) (client.UpdateScrubListEntryResult, error)
	RemoveScrubListEntries(ctx context.Context, params client.RemoveScrubListEntriesParams) (client.RemoveScrubListEntriesResult, error)

	// Skills Management
	AssignAgentSkill(ctx context.Context, params client.AssignAgentSkillParams) (client.AssignAgentSkillResult, error)
	UnassignAgentSkill(ctx context.Context, params client.UnassignAgentSkillParams) (client.UnassignAgentSkillResult, error)

	// Voice Recordings
	SearchVoiceRecordings(ctx context.Context, params client.SearchVoiceRecordingsParams) <-chan client.SearchVoiceRecordingsResult
	GetVoiceRecordingDownloadLink(ctx context.Context, params client.GetVoiceRecordingDownloadLinkParams) (client.GetVoiceRecordingDownloadLinkResult, error)
	ListSearchableRecordingFields(ctx context.Context, params client.ListSearchableRecordingFieldsParams) (client.ListSearchableRecordingFieldsResult, error)

	// Logging
	Log(ctx context.Context, params client.LogParams) (client.LogResult, error)
}
