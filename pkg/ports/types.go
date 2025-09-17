package ports

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

// --- GetAgentByID ---.
type GetAgentByIDParams struct {
	UserID string // Required
}

type GetAgentByIDResult struct {
	Agent *Agent
}

// --- GetClientConfiguration ---.
type GetClientConfigurationParams struct{}

type GetClientConfigurationResult struct {
	OrgID         string
	OrgName       string
	ConfigName    string
	ConfigPayload string
}

// --- ListNCLRulesetNames ---.
type ListNCLRulesetNamesParams struct{}

type ListNCLRulesetNamesResult struct {
	RulesetNames []string
}

// --- ListSkills ---.
type ListSkillsParams struct{}

type ListSkillsResult struct {
	Skills []Skill
}

type Skill struct {
	ID          string
	Name        string
	Description string
}

// --- ListAgentSkills ---.
type ListAgentSkillsParams struct {
	PartnerAgentID string // Required
}

type ListAgentSkillsResult struct {
	Skills []Skill
}

// --- AssignAgentSkill ---.
type AssignAgentSkillParams struct {
	PartnerAgentID string // Required
	SkillID        string // Required
}

type AssignAgentSkillResult struct{}

// --- UnassignAgentSkill ---.
type UnassignAgentSkillParams struct {
	PartnerAgentID string // Required
	SkillID        string // Required
}

type UnassignAgentSkillResult struct{}

// --- SearchVoiceRecordings ---.
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

// --- GetVoiceRecordingDownloadLink ---.
type GetVoiceRecordingDownloadLinkParams struct {
	RecordingSid string // Required
}

type GetVoiceRecordingDownloadLinkResult struct {
	DownloadURL string
	ExpiresAt   string
}

// --- ListSearchableRecordingFields ---.
type ListSearchableRecordingFieldsParams struct{}

type ListSearchableRecordingFieldsResult struct {
	Fields []SearchableField
}

type SearchableField struct {
	Name        string
	DisplayName string
	Type        string
}

// --- Transfer ---.
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

// --- PollEvents ---
type PollEventsParams struct{}

type PollEventsResult struct {
	Events []Event
}

type Event struct {
	Type          string
	Telephony     *ExileTelephonyResult
	AgentCall     *ExileAgentCall
	AgentResponse *ExileAgentResponse
}

type ExileTelephonyResult struct {
	CallSid        int64
	CallType       string
	CreateTime     string
	UpdateTime     string
	Status         string
	Result         string
	CallerID       string
	PhoneNumber    string
	StartTime      string
	EndTime        string
	DeliveryLength int64
	LinkbackLength int64
	PoolID         string
	RecordID       string
	ClientSid      int64
	OrgID          string
	InternalKey    string
}

type ExileAgentCall struct {
	AgentCallSid             int64
	CallSid                  int64
	CallType                 string
	TalkDuration             int64
	CallWaitDuration         int64
	WrapUpDuration           int64
	PauseDuration            int64
	TransferDuration         int64
	ManualDuration           int64
	PreviewDuration          int64
	HoldDuration             int64
	AgentWaitDuration        int64
	SuspendedDuration        int64
	ExternalTransferDuration int64
	CreateTime               string
	UpdateTime               string
	OrgID                    string
	UserID                   string
	InternalKey              string
	PartnerAgentID           string
}

type ExileAgentResponse struct {
	AgentCallResponseSid int64
	CallSid              int64
	CallType             string
	ResponseKey          string
	ResponseValue        string
	CreateTime           string
	UpdateTime           string
	ClientSid            int64
	OrgID                string
	AgentSid             int64
	UserID               string
	InternalKey          string
	PartnerAgentID       string
}

// --- StreamJobs ---
type StreamJobsParams struct{}

type StreamJobsResult struct {
	Job   *Job
	Error error
}

type Job struct {
	JobID string
	Type  string
	Data  map[string]interface{}
}

// --- SubmitJobResults ---
type SubmitJobResultsParams struct {
	JobID             string
	EndOfTransmission bool
	Results           map[string]interface{}
}

type SubmitJobResultsResult struct{}

// --- GetAgentStatus ---
type GetAgentStatusParams struct {
	PartnerAgentID string
}

type GetAgentStatusResult struct {
	AgentStatus *AgentStatus
}

type AgentStatus struct {
	PartnerAgentID string
	State          string
	Reason         string
	LastUpdate     string
}

// --- ListHuntGroupPauseCodes ---
type ListHuntGroupPauseCodesParams struct {
	PartnerAgentID string
}

type ListHuntGroupPauseCodesResult struct {
	PauseCodes []PauseCode
}

type PauseCode struct {
	Code        string
	Description string
	Duration    int64
}

// --- PutCallOnSimpleHold ---
type PutCallOnSimpleHoldParams struct {
	PartnerAgentID string
}

type PutCallOnSimpleHoldResult struct{}

// --- TakeCallOffSimpleHold ---
type TakeCallOffSimpleHoldParams struct {
	PartnerAgentID string
}

type TakeCallOffSimpleHoldResult struct{}

// --- StartCallRecording ---
type StartCallRecordingParams struct {
	PartnerAgentID string
}

type StartCallRecordingResult struct{}

// --- StopCallRecording ---
type StopCallRecordingParams struct {
	PartnerAgentID string
}

type StopCallRecordingResult struct{}

// --- GetRecordingStatus ---
type GetRecordingStatusParams struct {
	PartnerAgentID string
}

type GetRecordingStatusResult struct {
	IsRecording bool
}

// --- ListScrubLists ---
type ListScrubListsParams struct{}

type ListScrubListsResult struct {
	ScrubLists []ScrubList
}

type ScrubList struct {
	ID          string
	Name        string
	Description string
	CountryCode string
}

// --- UpdateScrubListEntry ---
type UpdateScrubListEntryParams struct {
	ScrubListID string
	EntryID     string
	Content     string
	Notes       *string
}

type UpdateScrubListEntryResult struct{}

// --- RemoveScrubListEntries ---
type RemoveScrubListEntriesParams struct {
	ScrubListID string
	EntryIDs    []string
}

type RemoveScrubListEntriesResult struct{}

// --- UpsertAgent ---
type UpsertAgentParams struct {
	UserID         string
	OrgID          string
	FirstName      string
	LastName       string
	Username       string
	PartnerAgentID string
}

type UpsertAgentResult struct{}

// --- GetAgentByPartnerID ---
type GetAgentByPartnerIDParams struct {
	PartnerAgentID string
}

type GetAgentByPartnerIDResult struct {
	Agent *Agent
}

// --- AddAgentCallResponse ---
type AddAgentCallResponseParams struct {
	PartnerAgentID string
	CallSid        int64
	ResponseKey    string
	ResponseValue  string
	AgentSid       int64
}

type AddAgentCallResponseResult struct{}

// --- Log ---
type LogParams struct {
	Level   string
	Message string
	Context map[string]interface{}
}

type LogResult struct{}

// --- RotateCertificate ---
type RotateCertificateParams struct{}

type RotateCertificateResult struct {
	Certificate   string
	PrivateKey    string
	CACertificate string
}

// --- GetOrganizationInfo ---
type GetOrganizationInfoParams struct{}

type GetOrganizationInfoResult struct {
	OrgID   string
	OrgName string
}
