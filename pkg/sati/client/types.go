package client

import (
	gatev2pb "github.com/tcncloud/sati-go/internal/genproto/tcnapi/exile/gate/v2"
)

// --- Common/Shared Types ---

// Agent represents agent information.
type Agent struct {
	UserID         string
	OrgID          string
	FirstName      string
	LastName       string
	Username       string
	PartnerAgentID string
	// Add other relevant fields from corev2.Agent if needed
}

// --- Dial ---

type DialParams struct {
	PartnerAgentID string // Required
	PhoneNumber    string // Required
	CallerID       *string
	PoolID         *string
	RecordID       *string
}

type DialResult struct {
	CallSid string
}

// --- AddScrubListEntries ---

type ScrubListEntryInput struct {
	Content string // Required
	Notes   *string
}

type AddScrubListEntriesParams struct {
	ScrubListID string                // Required
	Entries     []ScrubListEntryInput // Required, non-empty
	CountryCode *string
}

// AddScrubListEntriesResult represents the (currently empty) response.
// If the proto response adds fields later, add them here.
type AddScrubListEntriesResult struct{}

// --- ListAgents ---

// ListAgentsParams represents the (currently empty) request parameters.
// If the proto request adds fields later, add them here.
type ListAgentsParams struct{}

// ListAgentsResult contains an agent received from the stream.
type ListAgentsResult struct {
	Agent *Agent
	Error error // Used to propagate errors on the channel
}

// --- UpdateAgentStatus ---

type UpdateAgentStatusParams struct {
	PartnerAgentID string              // Required
	NewState       gatev2pb.AgentState // Required (Using proto enum directly is often okay)
	Reason         *string
}

// UpdateAgentStatusResult represents the (currently empty) response.
type UpdateAgentStatusResult struct{}

// --- GetAgentById ---
type GetAgentByIdParams struct {
	UserID string // Required
}

type GetAgentByIdResult struct {
	Agent *Agent
}

// --- GetClientConfiguration ---
type GetClientConfigurationParams struct{}

type GetClientConfigurationResult struct {
	OrgID         string
	OrgName       string
	ConfigName    string
	ConfigPayload string
}

// --- ListNCLRulesetNames ---
type ListNCLRulesetNamesParams struct{}

type ListNCLRulesetNamesResult struct {
	RulesetNames []string
}

// --- ListSkills ---
type ListSkillsParams struct{}

type ListSkillsResult struct {
	Skills []Skill
}

type Skill struct {
	ID          string
	Name        string
	Description string
}

// --- ListAgentSkills ---
type ListAgentSkillsParams struct {
	PartnerAgentID string // Required
}

type ListAgentSkillsResult struct {
	Skills []Skill
}

// --- AssignAgentSkill ---
type AssignAgentSkillParams struct {
	PartnerAgentID string // Required
	SkillID        string // Required
}

type AssignAgentSkillResult struct{}

// --- UnassignAgentSkill ---
type UnassignAgentSkillParams struct {
	PartnerAgentID string // Required
	SkillID        string // Required
}

type UnassignAgentSkillResult struct{}

// --- SearchVoiceRecordings ---
type SearchVoiceRecordingsParams struct {
	StartDate    *string
	EndDate      *string
	AgentID      *string
	CallSid      *string
	RecordingSid *string
	SearchFields []string
	SearchQuery  *string
	PageSize     *int32
	PageToken    *string
}

type SearchVoiceRecordingsResult struct {
	Recording *VoiceRecording
	NextToken string
	Error     error
}

type VoiceRecording struct {
	RecordingSid string
	CallSid      string
	AgentID      string
	StartTime    string
	EndTime      string
	Duration     int32
	FileSize     int64
	Status       string
}

// --- GetVoiceRecordingDownloadLink ---
type GetVoiceRecordingDownloadLinkParams struct {
	RecordingSid string // Required
}

type GetVoiceRecordingDownloadLinkResult struct {
	DownloadURL string
	ExpiresAt   string
}

// --- ListSearchableRecordingFields ---
type ListSearchableRecordingFieldsParams struct{}

type ListSearchableRecordingFieldsResult struct {
	Fields []SearchableField
}

type SearchableField struct {
	Name        string
	DisplayName string
	Type        string
}

// --- Transfer ---
type TransferParams struct {
	CallSid                 string // Required
	ReceivingPartnerAgentID *string
	Outbound                *TransferOutbound
	Queue                   *TransferQueue
}

type TransferOutbound struct {
	PhoneNumber string
	CallerID    *string
	PoolID      *string
	RecordID    *string
}

type TransferQueue struct {
	QueueID string
}

type TransferResult struct{}

// Add definitions for other methods (PollEvents, StreamJobs, etc.) following this pattern.
