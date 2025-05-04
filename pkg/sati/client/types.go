package client

import (
	gatev2pb "buf.build/gen/go/tcn/exileapi/protocolbuffers/go/tcnapi/exile/gate/v2"
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

// Add definitions for other methods (PollEvents, StreamJobs, etc.) following this pattern.
