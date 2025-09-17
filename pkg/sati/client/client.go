package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"

	gatev2pb "github.com/tcncloud/sati-go/internal/genproto/tcnapi/exile/gate/v2" // Keep for internal mapping
	"github.com/tcncloud/sati-go/pkg/ports"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/wrapperspb" // Needed for optional fields
)

// Common error constants for client operations.
var (
	ErrScrubListIDRequired     = errors.New("ScrubListID and at least one Entry are required")
	ErrEntryContentEmpty       = errors.New("entry content cannot be empty")
	ErrDialParamsRequired      = errors.New("PartnerAgentID and PhoneNumber are required")
	ErrDialResponseNil         = errors.New("received nil response from gRPC Dial")
	ErrUserIDRequired          = errors.New("UserID is required")
	ErrAgentNotFound           = errors.New("agent not found or nil response received")
	ErrClientConfigResponseNil = errors.New("received nil response from gRPC GetClientConfiguration")
	ErrListAgentsStreamNil     = errors.New("received nil agent in ListAgents stream")
	ErrPartnerAgentIDRequired  = errors.New("PartnerAgentID is required")
	ErrCAAppendFailed          = errors.New("failed to append CA cert")
)

// Client provides methods for interacting with the GateService API.
type Client struct {
	conn *grpc.ClientConn
	gate gatev2pb.GateServiceClient
}

// Ensure Client implements the ClientInterface interface
var _ ports.ClientInterface = (*Client)(nil)

// NewClient creates a new Sati API client.
// It takes the configuration and sets up the gRPC connection and client stub.
func NewClient(cfg *saticonfig.Config) (*Client, error) {
	conn, err := setupConnection(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
		gate: gatev2pb.NewGateServiceClient(conn),
	}, nil
}

// Close terminates the underlying gRPC connection.
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

// --- GateService Methods (Refactored) ---

// AddAgentCallResponse adds a response to an agent call.
func (c *Client) AddAgentCallResponse(ctx context.Context, params ports.AddAgentCallResponseParams) (ports.AddAgentCallResponseResult, error) {
	req := &gatev2pb.AddAgentCallResponseRequest{
		PartnerAgentId:   params.PartnerAgentID,
		CallSid:          strconv.FormatInt(params.CallSid, 10),
		CallType:         gatev2pb.CallType_CALL_TYPE_INBOUND, // Default to inbound
		CurrentSessionId: params.AgentSid,
		Key:              params.ResponseKey,
		Value:            params.ResponseValue,
	}

	_, err := c.gate.AddAgentCallResponse(ctx, req)
	if err != nil {
		return ports.AddAgentCallResponseResult{}, err
	}

	return ports.AddAgentCallResponseResult{}, nil
}

// AddScrubListEntries adds entries to a scrub list.
func (c *Client) AddScrubListEntries(ctx context.Context, params ports.AddScrubListEntriesParams) (ports.AddScrubListEntriesResult, error) {
	if params.ScrubListID == "" || len(params.Entries) == 0 {
		return ports.AddScrubListEntriesResult{}, ErrScrubListIDRequired
	}

	pbEntries := make([]*gatev2pb.AddScrubListEntriesRequest_Entry, 0, len(params.Entries))
	for _, e := range params.Entries {
		if e.Content == "" {
			return ports.AddScrubListEntriesResult{}, ErrEntryContentEmpty
		}

		pbEntry := &gatev2pb.AddScrubListEntriesRequest_Entry{
			Content: e.Content,
		}
		if e.Notes != nil {
			pbEntry.Notes = wrapperspb.String(*e.Notes)
		}

		pbEntries = append(pbEntries, pbEntry)
	}

	req := &gatev2pb.AddScrubListEntriesRequest{
		ScrubListId: params.ScrubListID,
		Entries:     pbEntries,
	}
	if params.CountryCode != nil {
		req.CountryCode = *params.CountryCode // Assuming proto field is string, not wrapper
	}

	_, err := c.gate.AddScrubListEntries(ctx, req)
	if err != nil {
		return ports.AddScrubListEntriesResult{}, err
	}

	return ports.AddScrubListEntriesResult{}, nil // Return empty struct for success
}

// Dial initiates an outbound call.
func (c *Client) Dial(ctx context.Context, params ports.DialParams) (ports.DialResult, error) {
	if params.PartnerAgentID == "" || params.PhoneNumber == "" {
		return ports.DialResult{}, ErrDialParamsRequired
	}

	req := &gatev2pb.DialRequest{
		PartnerAgentId: params.PartnerAgentID,
		PhoneNumber:    params.PhoneNumber,
	}
	if params.CallerID != nil {
		req.CallerId = wrapperspb.String(*params.CallerID)
	}

	if params.PoolID != nil {
		req.PoolId = wrapperspb.String(*params.PoolID)
	}

	if params.RecordID != nil {
		req.RecordId = wrapperspb.String(*params.RecordID)
	}

	resp, err := c.gate.Dial(ctx, req)
	if err != nil {
		return ports.DialResult{}, err
	}

	if resp == nil {
		return ports.DialResult{}, ErrDialResponseNil
	}

	result := ports.DialResult{
		CallSid: resp.GetCallSid(),
	}

	return result, nil
}

// GetAgentById retrieves agent details by User ID.
func (c *Client) GetAgentByID(ctx context.Context, params ports.GetAgentByIDParams) (ports.GetAgentByIDResult, error) {
	if params.UserID == "" {
		return ports.GetAgentByIDResult{}, ErrUserIDRequired
	}

	req := &gatev2pb.GetAgentByIdRequest{UserId: params.UserID}

	resp, err := c.gate.GetAgentById(ctx, req)
	if err != nil {
		return ports.GetAgentByIDResult{}, err
	}

	if resp == nil || resp.GetAgent() == nil {
		// Consider returning a specific "not found" error type here
		return ports.GetAgentByIDResult{}, ErrAgentNotFound
	}

	result := ports.GetAgentByIDResult{
		Agent: mapProtoAgentToAgent(resp.GetAgent()), // Use mapping helper
	}

	return result, nil
}

// GetAgentByPartnerID retrieves agent details by Partner Agent ID.
func (c *Client) GetAgentByPartnerID(ctx context.Context, params ports.GetAgentByPartnerIDParams) (ports.GetAgentByPartnerIDResult, error) {
	if params.PartnerAgentID == "" {
		return ports.GetAgentByPartnerIDResult{}, ErrPartnerAgentIDRequired
	}

	req := &gatev2pb.GetAgentByPartnerIdRequest{
		PartnerAgentId: params.PartnerAgentID,
	}

	resp, err := c.gate.GetAgentByPartnerId(ctx, req)
	if err != nil {
		return ports.GetAgentByPartnerIDResult{}, err
	}

	if resp == nil || resp.GetAgent() == nil {
		return ports.GetAgentByPartnerIDResult{}, ErrAgentNotFound
	}

	result := ports.GetAgentByPartnerIDResult{
		Agent: mapProtoAgentToAgent(resp.GetAgent()),
	}

	return result, nil
}

// GetAgentStatus retrieves the current state of an agent.
func (c *Client) GetAgentStatus(ctx context.Context, params ports.GetAgentStatusParams) (ports.GetAgentStatusResult, error) {
	if params.PartnerAgentID == "" {
		return ports.GetAgentStatusResult{}, ErrPartnerAgentIDRequired
	}

	req := &gatev2pb.GetAgentStatusRequest{
		PartnerAgentId: params.PartnerAgentID,
	}

	resp, err := c.gate.GetAgentStatus(ctx, req)
	if err != nil {
		return ports.GetAgentStatusResult{}, err
	}

	if resp == nil {
		return ports.GetAgentStatusResult{}, errors.New("received nil response from gRPC GetAgentStatus")
	}

	result := ports.GetAgentStatusResult{
		AgentStatus: &ports.AgentStatus{
			PartnerAgentID: resp.GetPartnerAgentId(),
			State:          resp.GetAgentState().String(),
			Reason:         "",                              // Not available in this response
			LastUpdate:     time.Now().Format(time.RFC3339), // Not available in this response
		},
	}

	return result, nil
}

// GetClientConfiguration retrieves client configuration details.
func (c *Client) GetClientConfiguration(ctx context.Context, params ports.GetClientConfigurationParams) (ports.GetClientConfigurationResult, error) {
	req := &gatev2pb.GetClientConfigurationRequest{}

	resp, err := c.gate.GetClientConfiguration(ctx, req)
	if err != nil {
		return ports.GetClientConfigurationResult{}, err
	}

	if resp == nil {
		return ports.GetClientConfigurationResult{}, ErrClientConfigResponseNil
	}

	result := ports.GetClientConfigurationResult{
		OrgID:         resp.GetOrgId(),
		OrgName:       resp.GetOrgName(),
		ConfigName:    resp.GetConfigName(),
		ConfigPayload: resp.GetConfigPayload(),
	}

	return result, nil
}

// GetOrganizationInfo retrieves organization details.
func (c *Client) GetOrganizationInfo(ctx context.Context, params ports.GetOrganizationInfoParams) (ports.GetOrganizationInfoResult, error) {
	req := &gatev2pb.GetOrganizationInfoRequest{}

	resp, err := c.gate.GetOrganizationInfo(ctx, req)
	if err != nil {
		return ports.GetOrganizationInfoResult{}, err
	}

	if resp == nil {
		return ports.GetOrganizationInfoResult{}, errors.New("received nil response from gRPC GetOrganizationInfo")
	}

	result := ports.GetOrganizationInfoResult{
		OrgID:   resp.GetOrgId(),
		OrgName: resp.GetOrgName(),
	}

	return result, nil
}

// GetRecordingStatus checks if a call is currently being recorded.
func (c *Client) GetRecordingStatus(ctx context.Context, params ports.GetRecordingStatusParams) (ports.GetRecordingStatusResult, error) {
	req := &gatev2pb.GetRecordingStatusRequest{
		PartnerAgentId: params.PartnerAgentID,
	}

	resp, err := c.gate.GetRecordingStatus(ctx, req)
	if err != nil {
		return ports.GetRecordingStatusResult{}, err
	}

	if resp == nil {
		return ports.GetRecordingStatusResult{}, errors.New("received nil response from gRPC GetRecordingStatus")
	}

	result := ports.GetRecordingStatusResult{
		IsRecording: resp.GetIsRecording(),
	}

	return result, nil
}

// ListAgents returns a channel that emits agent details.
func (c *Client) ListAgents(ctx context.Context, params ports.ListAgentsParams) <-chan ports.ListAgentsResult {
	resultsChan := make(chan ports.ListAgentsResult)
	req := &gatev2pb.ListAgentsRequest{}

	go func() {
		defer close(resultsChan)

		stream, err := c.gate.ListAgents(ctx, req)
		if err != nil {
			resultsChan <- ports.ListAgentsResult{Error: fmt.Errorf("failed to start ListAgents stream: %w", err)}

			return
		}

		for {
			resp, err := stream.Recv()
			if err != nil {
				if !IsStreamEnd(err) { // Don't send EOF as error
					resultsChan <- ports.ListAgentsResult{Error: fmt.Errorf("error receiving from ListAgents stream: %w", err)}
				}

				return // End goroutine on EOF or error
			}

			if resp == nil || resp.GetAgent() == nil {
				// Handle nil response/agent, maybe log or send specific error
				resultsChan <- ports.ListAgentsResult{Error: ErrListAgentsStreamNil}

				continue
			}

			resultsChan <- ports.ListAgentsResult{Agent: mapProtoAgentToAgent(resp.GetAgent())}
		}
	}()

	return resultsChan
}

// ListHuntGroupPauseCodes lists the pause codes defined for a hunt group.
func (c *Client) ListHuntGroupPauseCodes(ctx context.Context, params ports.ListHuntGroupPauseCodesParams) (ports.ListHuntGroupPauseCodesResult, error) {
	req := &gatev2pb.ListHuntGroupPauseCodesRequest{
		PartnerAgentId: params.PartnerAgentID,
	}

	resp, err := c.gate.ListHuntGroupPauseCodes(ctx, req)
	if err != nil {
		return ports.ListHuntGroupPauseCodesResult{}, err
	}

	if resp == nil {
		return ports.ListHuntGroupPauseCodesResult{}, errors.New("received nil response from gRPC ListHuntGroupPauseCodes")
	}

	pauseCodes := make([]ports.PauseCode, 0, len(resp.GetPauseCodes()))
	for _, pc := range resp.GetPauseCodes() {
		pauseCodes = append(pauseCodes, ports.PauseCode{
			Code:        pc,
			Description: "", // Not available in this response
			Duration:    0,  // Not available in this response
		})
	}

	result := ports.ListHuntGroupPauseCodesResult{
		PauseCodes: pauseCodes,
	}

	return result, nil
}

// ListScrubLists lists all available scrub lists.
func (c *Client) ListScrubLists(ctx context.Context, params ports.ListScrubListsParams) (ports.ListScrubListsResult, error) {
	req := &gatev2pb.ListScrubListsRequest{}

	resp, err := c.gate.ListScrubLists(ctx, req)
	if err != nil {
		return ports.ListScrubListsResult{}, err
	}

	if resp == nil {
		return ports.ListScrubListsResult{}, errors.New("received nil response from gRPC ListScrubLists")
	}

	scrubLists := make([]ports.ScrubList, 0, len(resp.GetScrubLists()))
	for _, sl := range resp.GetScrubLists() {
		scrubLists = append(scrubLists, ports.ScrubList{
			ID:          sl.GetScrubListId(),
			Name:        "", // Not available in this response
			Description: "", // Not available in this response
			CountryCode: "", // Not available in this response
		})
	}

	result := ports.ListScrubListsResult{
		ScrubLists: scrubLists,
	}

	return result, nil
}

// Log sends a log message.
func (c *Client) Log(ctx context.Context, params ports.LogParams) (ports.LogResult, error) {
	// Create a structured log message
	logMessage := fmt.Sprintf("Level: %s, Message: %s", params.Level, params.Message)
	req := &gatev2pb.LogRequest{
		Payload: logMessage,
	}

	_, err := c.gate.Log(ctx, req)
	if err != nil {
		return ports.LogResult{}, err
	}

	return ports.LogResult{}, nil
}

// PollEvents polls for events from the Operator platform.
func (c *Client) PollEvents(ctx context.Context, params ports.PollEventsParams) (ports.PollEventsResult, error) {
	req := &gatev2pb.PollEventsRequest{}

	resp, err := c.gate.PollEvents(ctx, req)
	if err != nil {
		return ports.PollEventsResult{}, err
	}

	if resp == nil {
		return ports.PollEventsResult{}, errors.New("received nil response from gRPC PollEvents")
	}

	events := make([]ports.Event, 0, len(resp.GetEvents()))
	for _, event := range resp.GetEvents() {
		e := ports.Event{}

		if telephony := event.GetTelephonyResult(); telephony != nil {
			e.Telephony = &ports.ExileTelephonyResult{
				CallSid:        telephony.GetCallSid(),
				CallType:       telephony.GetCallType(),
				CreateTime:     telephony.GetCreateTime().AsTime().Format(time.RFC3339),
				UpdateTime:     telephony.GetUpdateTime().AsTime().Format(time.RFC3339),
				Status:         telephony.GetStatus().String(),
				Result:         telephony.GetResult().String(),
				CallerID:       telephony.GetCallerId(),
				PhoneNumber:    telephony.GetPhoneNumber(),
				StartTime:      telephony.GetStartTime().AsTime().Format(time.RFC3339),
				EndTime:        telephony.GetEndTime().AsTime().Format(time.RFC3339),
				DeliveryLength: telephony.GetDeliveryLength(),
				LinkbackLength: telephony.GetLinkbackLength(),
				PoolID:         telephony.GetPoolId(),
				RecordID:       telephony.GetRecordId(),
				ClientSid:      telephony.GetClientSid(),
				OrgID:          telephony.GetOrgId(),
				InternalKey:    telephony.GetInternalKey(),
			}
		}

		if agentCall := event.GetAgentCall(); agentCall != nil {
			e.AgentCall = &ports.ExileAgentCall{
				AgentCallSid:             agentCall.GetAgentCallSid(),
				CallSid:                  agentCall.GetCallSid(),
				CallType:                 agentCall.GetCallType(),
				TalkDuration:             agentCall.GetTalkDuration(),
				CallWaitDuration:         agentCall.GetCallWaitDuration(),
				WrapUpDuration:           agentCall.GetWrapUpDuration(),
				PauseDuration:            agentCall.GetPauseDuration(),
				TransferDuration:         agentCall.GetTransferDuration(),
				ManualDuration:           agentCall.GetManualDuration(),
				PreviewDuration:          agentCall.GetPreviewDuration(),
				HoldDuration:             agentCall.GetHoldDuration(),
				AgentWaitDuration:        agentCall.GetAgentWaitDuration(),
				SuspendedDuration:        agentCall.GetSuspendedDuration(),
				ExternalTransferDuration: agentCall.GetExternalTransferDuration(),
				CreateTime:               agentCall.GetCreateTime().AsTime().Format(time.RFC3339),
				UpdateTime:               agentCall.GetUpdateTime().AsTime().Format(time.RFC3339),
				OrgID:                    agentCall.GetOrgId(),
				UserID:                   agentCall.GetUserId(),
				InternalKey:              agentCall.GetInternalKey(),
				PartnerAgentID:           agentCall.GetPartnerAgentId(),
			}
		}

		if agentResponse := event.GetAgentResponse(); agentResponse != nil {
			e.AgentResponse = &ports.ExileAgentResponse{
				AgentCallResponseSid: agentResponse.GetAgentCallResponseSid(),
				CallSid:              agentResponse.GetCallSid(),
				CallType:             agentResponse.GetCallType(),
				ResponseKey:          agentResponse.GetResponseKey(),
				ResponseValue:        agentResponse.GetResponseValue(),
				CreateTime:           agentResponse.GetCreateTime().AsTime().Format(time.RFC3339),
				UpdateTime:           agentResponse.GetUpdateTime().AsTime().Format(time.RFC3339),
				ClientSid:            agentResponse.GetClientSid(),
				OrgID:                agentResponse.GetOrgId(),
				AgentSid:             agentResponse.GetAgentSid(),
				UserID:               agentResponse.GetUserId(),
				InternalKey:          agentResponse.GetInternalKey(),
				PartnerAgentID:       agentResponse.GetPartnerAgentId(),
			}
		}

		events = append(events, e)
	}

	result := ports.PollEventsResult{
		Events: events,
	}

	return result, nil
}

// PutCallOnSimpleHold puts a call on simple hold.
func (c *Client) PutCallOnSimpleHold(ctx context.Context, params ports.PutCallOnSimpleHoldParams) (ports.PutCallOnSimpleHoldResult, error) {
	req := &gatev2pb.PutCallOnSimpleHoldRequest{
		PartnerAgentId: params.PartnerAgentID,
	}

	_, err := c.gate.PutCallOnSimpleHold(ctx, req)
	if err != nil {
		return ports.PutCallOnSimpleHoldResult{}, err
	}

	return ports.PutCallOnSimpleHoldResult{}, nil
}

// RemoveScrubListEntries removes entries from a scrub list.
func (c *Client) RemoveScrubListEntries(ctx context.Context, params ports.RemoveScrubListEntriesParams) (ports.RemoveScrubListEntriesResult, error) {
	req := &gatev2pb.RemoveScrubListEntriesRequest{
		ScrubListId: params.ScrubListID,
		Entries:     params.EntryIDs,
	}

	_, err := c.gate.RemoveScrubListEntries(ctx, req)
	if err != nil {
		return ports.RemoveScrubListEntriesResult{}, err
	}

	return ports.RemoveScrubListEntriesResult{}, nil
}

// RotateCertificate rotates the client certificate.
func (c *Client) RotateCertificate(ctx context.Context, params ports.RotateCertificateParams) (ports.RotateCertificateResult, error) {
	req := &gatev2pb.RotateCertificateRequest{}

	resp, err := c.gate.RotateCertificate(ctx, req)
	if err != nil {
		return ports.RotateCertificateResult{}, err
	}

	if resp == nil {
		return ports.RotateCertificateResult{}, errors.New("received nil response from gRPC RotateCertificate")
	}

	result := ports.RotateCertificateResult{
		Certificate:   resp.GetEncodedCertificate(),
		PrivateKey:    "", // Not available in this response
		CACertificate: "", // Not available in this response
	}

	return result, nil
}

// StartCallRecording starts recording a call.
func (c *Client) StartCallRecording(ctx context.Context, params ports.StartCallRecordingParams) (ports.StartCallRecordingResult, error) {
	req := &gatev2pb.StartCallRecordingRequest{
		PartnerAgentId: params.PartnerAgentID,
	}

	_, err := c.gate.StartCallRecording(ctx, req)
	if err != nil {
		return ports.StartCallRecordingResult{}, err
	}

	return ports.StartCallRecordingResult{}, nil
}

// StopCallRecording stops recording a call.
func (c *Client) StopCallRecording(ctx context.Context, params ports.StopCallRecordingParams) (ports.StopCallRecordingResult, error) {
	req := &gatev2pb.StopCallRecordingRequest{
		PartnerAgentId: params.PartnerAgentID,
	}

	_, err := c.gate.StopCallRecording(ctx, req)
	if err != nil {
		return ports.StopCallRecordingResult{}, err
	}

	return ports.StopCallRecordingResult{}, nil
}

// StreamJobs returns a channel that emits jobs from the Operator platform.
func (c *Client) StreamJobs(ctx context.Context, params ports.StreamJobsParams) <-chan ports.StreamJobsResult {
	resultChan := make(chan ports.StreamJobsResult, 1)

	go func() {
		defer close(resultChan)

		req := &gatev2pb.StreamJobsRequest{}

		stream, err := c.gate.StreamJobs(ctx, req)
		if err != nil {
			resultChan <- ports.StreamJobsResult{Error: err}
			return
		}

		for {
			resp, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}

			if err != nil {
				resultChan <- ports.StreamJobsResult{Error: err}
				return
			}

			// Process the job in the response
			jobData := make(map[string]interface{})
			// Convert job data to map - this would need to be implemented based on the actual job structure
			jobData["job_id"] = resp.GetJobId()
			jobData["type"] = "" // Type not available in this response

			resultChan <- ports.StreamJobsResult{
				Job: &ports.Job{
					JobID: resp.GetJobId(),
					Type:  "", // Type not available in this response
					Data:  jobData,
				},
			}
		}
	}()

	return resultChan
}

// SubmitJobResults submits results for jobs received via StreamJobs.
func (c *Client) SubmitJobResults(ctx context.Context, params ports.SubmitJobResultsParams) (ports.SubmitJobResultsResult, error) {
	req := &gatev2pb.SubmitJobResultsRequest{
		JobId:             params.JobID,
		EndOfTransmission: params.EndOfTransmission,
	}

	// Convert results map to protobuf format
	// This would need to be implemented based on the actual job result structure
	// For now, we'll leave it empty

	_, err := c.gate.SubmitJobResults(ctx, req)
	if err != nil {
		return ports.SubmitJobResultsResult{}, err
	}

	return ports.SubmitJobResultsResult{}, nil
}

// TakeCallOffSimpleHold takes a call off simple hold.
func (c *Client) TakeCallOffSimpleHold(ctx context.Context, params ports.TakeCallOffSimpleHoldParams) (ports.TakeCallOffSimpleHoldResult, error) {
	req := &gatev2pb.TakeCallOffSimpleHoldRequest{
		PartnerAgentId: params.PartnerAgentID,
	}

	_, err := c.gate.TakeCallOffSimpleHold(ctx, req)
	if err != nil {
		return ports.TakeCallOffSimpleHoldResult{}, err
	}

	return ports.TakeCallOffSimpleHoldResult{}, nil
}

// UpdateAgentStatus updates the state of an agent.
func (c *Client) UpdateAgentStatus(ctx context.Context, params ports.UpdateAgentStatusParams) (ports.UpdateAgentStatusResult, error) {
	if params.PartnerAgentID == "" {
		return ports.UpdateAgentStatusResult{}, ErrPartnerAgentIDRequired
	}
	// Note: We assume AgentState enum values are stable and okay to expose directly.
	req := &gatev2pb.UpdateAgentStatusRequest{
		PartnerAgentId: params.PartnerAgentID,
		NewState:       params.NewState,
	}
	if params.Reason != nil {
		// Assuming the proto field is just a string, not a wrapper
		req.Reason = *params.Reason
	}

	_, err := c.gate.UpdateAgentStatus(ctx, req)
	if err != nil {
		return ports.UpdateAgentStatusResult{}, err
	}

	return ports.UpdateAgentStatusResult{}, nil
}

// UpdateScrubListEntry updates an existing scrub list entry.
func (c *Client) UpdateScrubListEntry(ctx context.Context, params ports.UpdateScrubListEntryParams) (ports.UpdateScrubListEntryResult, error) {
	req := &gatev2pb.UpdateScrubListEntryRequest{
		ScrubListId: params.ScrubListID,
		Content:     params.Content,
	}

	if params.Notes != nil {
		req.Notes = wrapperspb.String(*params.Notes)
	}

	_, err := c.gate.UpdateScrubListEntry(ctx, req)
	if err != nil {
		return ports.UpdateScrubListEntryResult{}, err
	}

	return ports.UpdateScrubListEntryResult{}, nil
}

// UpsertAgent creates or updates agent information.
func (c *Client) UpsertAgent(ctx context.Context, params ports.UpsertAgentParams) (ports.UpsertAgentResult, error) {
	req := &gatev2pb.UpsertAgentRequest{
		Username:       params.Username,
		PartnerAgentId: params.PartnerAgentID,
		FirstName:      params.FirstName,
		LastName:       params.LastName,
		Password:       "", // Password not available in params
	}

	_, err := c.gate.UpsertAgent(ctx, req)
	if err != nil {
		return ports.UpsertAgentResult{}, err
	}

	return ports.UpsertAgentResult{}, nil
}

// --- New Methods for Missing Commands ---

// ListNCLRulesetNames calls the ListNCLRulesetNames RPC.
func (c *Client) ListNCLRulesetNames(ctx context.Context, params ports.ListNCLRulesetNamesParams) (ports.ListNCLRulesetNamesResult, error) {
	req := &gatev2pb.ListNCLRulesetNamesRequest{}

	resp, err := c.gate.ListNCLRulesetNames(ctx, req)
	if err != nil {
		return ports.ListNCLRulesetNamesResult{}, err
	}

	return ports.ListNCLRulesetNamesResult{
		RulesetNames: resp.GetRulesetNames(),
	}, nil
}

// ListSkills calls the ListSkills RPC.
func (c *Client) ListSkills(ctx context.Context, params ports.ListSkillsParams) (ports.ListSkillsResult, error) {
	req := &gatev2pb.ListSkillsRequest{}

	resp, err := c.gate.ListSkills(ctx, req)
	if err != nil {
		return ports.ListSkillsResult{}, err
	}

	skills := make([]ports.Skill, 0, len(resp.GetSkills()))
	for _, skill := range resp.GetSkills() {
		skills = append(skills, ports.Skill{
			ID:          skill.GetSkillId(),
			Name:        skill.GetName(),
			Description: skill.GetDescription(),
		})
	}

	return ports.ListSkillsResult{
		Skills: skills,
	}, nil
}

// ListAgentSkills calls the ListAgentSkills RPC.
func (c *Client) ListAgentSkills(ctx context.Context, params ports.ListAgentSkillsParams) (ports.ListAgentSkillsResult, error) {
	req := &gatev2pb.ListAgentSkillsRequest{
		PartnerAgentId: params.PartnerAgentID,
	}

	resp, err := c.gate.ListAgentSkills(ctx, req)
	if err != nil {
		return ports.ListAgentSkillsResult{}, err
	}

	skills := make([]ports.Skill, 0, len(resp.GetSkills()))
	for _, skill := range resp.GetSkills() {
		skills = append(skills, ports.Skill{
			ID:          skill.GetSkillId(),
			Name:        skill.GetName(),
			Description: skill.GetDescription(),
		})
	}

	return ports.ListAgentSkillsResult{
		Skills: skills,
	}, nil
}

// AssignAgentSkill calls the AssignAgentSkill RPC.
func (c *Client) AssignAgentSkill(ctx context.Context, params ports.AssignAgentSkillParams) (ports.AssignAgentSkillResult, error) {
	req := &gatev2pb.AssignAgentSkillRequest{
		PartnerAgentId: params.PartnerAgentID,
		SkillId:        params.SkillID,
	}

	_, err := c.gate.AssignAgentSkill(ctx, req)
	if err != nil {
		return ports.AssignAgentSkillResult{}, err
	}

	return ports.AssignAgentSkillResult{}, nil
}

// UnassignAgentSkill calls the UnassignAgentSkill RPC.
func (c *Client) UnassignAgentSkill(ctx context.Context, params ports.UnassignAgentSkillParams) (ports.UnassignAgentSkillResult, error) {
	req := &gatev2pb.UnassignAgentSkillRequest{
		PartnerAgentId: params.PartnerAgentID,
		SkillId:        params.SkillID,
	}

	_, err := c.gate.UnassignAgentSkill(ctx, req)
	if err != nil {
		return ports.UnassignAgentSkillResult{}, err
	}

	return ports.UnassignAgentSkillResult{}, nil
}

// SearchVoiceRecordings calls the SearchVoiceRecordings RPC (streaming).
func (c *Client) SearchVoiceRecordings(ctx context.Context, params ports.SearchVoiceRecordingsParams) <-chan ports.SearchVoiceRecordingsResult {
	resultChan := make(chan ports.SearchVoiceRecordingsResult, 1)

	go func() {
		defer close(resultChan)

		req := &gatev2pb.SearchVoiceRecordingsRequest{}

		// Build search options from parameters
		var searchOptions []*gatev2pb.SearchOption

		if params.StartDate != nil {
			searchOptions = append(searchOptions, &gatev2pb.SearchOption{
				Field:    "start_time",
				Operator: gatev2pb.Operator_EQUAL, // Use EQUAL for now, date comparison might need different approach
				Value:    *params.StartDate,
			})
		}

		if params.EndDate != nil {
			searchOptions = append(searchOptions, &gatev2pb.SearchOption{
				Field:    "start_time",
				Operator: gatev2pb.Operator_EQUAL, // Use EQUAL for now, date comparison might need different approach
				Value:    *params.EndDate,
			})
		}

		if params.AgentID != nil {
			searchOptions = append(searchOptions, &gatev2pb.SearchOption{
				Field:    "partner_agent_ids",
				Operator: gatev2pb.Operator_CONTAINS,
				Value:    *params.AgentID,
			})
		}

		if params.CallSid != nil {
			searchOptions = append(searchOptions, &gatev2pb.SearchOption{
				Field:    "call_sid",
				Operator: gatev2pb.Operator_EQUAL,
				Value:    *params.CallSid,
			})
		}

		req.SearchOptions = searchOptions

		stream, err := c.gate.SearchVoiceRecordings(ctx, req)
		if err != nil {
			resultChan <- ports.SearchVoiceRecordingsResult{Error: err}

			return
		}

		for {
			resp, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}

			if err != nil {
				resultChan <- ports.SearchVoiceRecordingsResult{Error: err}

				return
			}

			// Process each recording in the response
			for _, recording := range resp.GetRecordings() {
				voiceRecording := &ports.VoiceRecording{
					RecordingSid: recording.GetName(), // Use name as recording ID
					CallSid:      strconv.FormatInt(recording.GetCallSid(), 10),
					AgentID:      strings.Join(recording.GetPartnerAgentIds(), ","),
					StartTime:    recording.GetStartTime().AsTime().Format(time.RFC3339),
					EndTime:      recording.GetStartTime().AsTime().Add(recording.GetDuration().AsDuration()).Format(time.RFC3339),
					Duration:     int32(recording.GetDuration().AsDuration().Seconds()),
					FileSize:     0,           // Not available in this response
					Status:       "completed", // Assume completed for now
				}

				resultChan <- ports.SearchVoiceRecordingsResult{
					Recording: voiceRecording,
					NextToken: "", // Not available in this response
				}
			}
		}
	}()

	return resultChan
}

// GetVoiceRecordingDownloadLink calls the GetVoiceRecordingDownloadLink RPC.
func (c *Client) GetVoiceRecordingDownloadLink(ctx context.Context, params ports.GetVoiceRecordingDownloadLinkParams) (ports.GetVoiceRecordingDownloadLinkResult, error) {
	req := &gatev2pb.GetVoiceRecordingDownloadLinkRequest{
		RecordingId: params.RecordingSid,
	}

	resp, err := c.gate.GetVoiceRecordingDownloadLink(ctx, req)
	if err != nil {
		return ports.GetVoiceRecordingDownloadLinkResult{}, err
	}

	return ports.GetVoiceRecordingDownloadLinkResult{
		DownloadURL: resp.GetDownloadLink(),
		ExpiresAt:   "", // Not available in this response
	}, nil
}

// ListSearchableRecordingFields calls the ListSearchableRecordingFields RPC.
func (c *Client) ListSearchableRecordingFields(ctx context.Context, params ports.ListSearchableRecordingFieldsParams) (ports.ListSearchableRecordingFieldsResult, error) {
	req := &gatev2pb.ListSearchableRecordingFieldsRequest{}

	resp, err := c.gate.ListSearchableRecordingFields(ctx, req)
	if err != nil {
		return ports.ListSearchableRecordingFieldsResult{}, err
	}

	fields := make([]ports.SearchableField, 0, len(resp.GetFields()))
	for _, fieldName := range resp.GetFields() {
		fields = append(fields, ports.SearchableField{
			Name:        fieldName,
			DisplayName: fieldName, // Use field name as display name
			Type:        "string",  // Default type
		})
	}

	return ports.ListSearchableRecordingFieldsResult{
		Fields: fields,
	}, nil
}

// Transfer calls the Transfer RPC.
func (c *Client) Transfer(ctx context.Context, params ports.TransferParams) (ports.TransferResult, error) {
	req := &gatev2pb.TransferRequest{
		PartnerAgentId: params.CallSid, // Use CallSid as PartnerAgentId for now
	}

	// Set destination based on transfer type
	switch {
	case params.ReceivingPartnerAgentID != nil:
		req.Destination = &gatev2pb.TransferRequest_ReceivingPartnerAgentId{
			ReceivingPartnerAgentId: &gatev2pb.TransferRequest_Agent{
				PartnerAgentId: *params.ReceivingPartnerAgentID,
			},
		}
	case params.Outbound != nil:
		outbound := &gatev2pb.TransferRequest_Outbound{
			Destination: params.Outbound.PhoneNumber,
		}
		if params.Outbound.CallerID != nil {
			outbound.CallerId = *params.Outbound.CallerID
		}

		req.Destination = &gatev2pb.TransferRequest_Outbound_{
			Outbound: outbound,
		}
	case params.Queue != nil:
		req.Destination = &gatev2pb.TransferRequest_Queue_{
			Queue: &gatev2pb.TransferRequest_Queue{},
		}
	}

	_, err := c.gate.Transfer(ctx, req)
	if err != nil {
		return ports.TransferResult{}, err
	}

	return ports.TransferResult{}, nil
}

// --- Internal Helper Functions ---

// mapProtoAgentToAgent converts a proto Agent to our custom client Agent type.
func mapProtoAgentToAgent(pbAgent *gatev2pb.Agent) *ports.Agent {
	if pbAgent == nil {
		return nil
	}

	return &ports.Agent{
		UserID:         pbAgent.GetUserId(), // Use Getters for safety
		OrgID:          pbAgent.GetOrgId(),
		FirstName:      pbAgent.GetFirstName(),
		LastName:       pbAgent.GetLastName(),
		Username:       pbAgent.GetUsername(),
		PartnerAgentID: pbAgent.GetPartnerAgentId(),
	}
}

// setupConnection configures and establishes the gRPC connection.
func setupConnection(cfg *saticonfig.Config) (*grpc.ClientConn, error) {
	cert, err := tls.X509KeyPair([]byte(cfg.Certificate), []byte(cfg.PrivateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to load client cert: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM([]byte(cfg.CACertificate)); !ok {
		return nil, ErrCAAppendFailed
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS12, // Set minimum TLS version to 1.2
	})

	endpoint := parseAPIEndpoint(cfg.APIEndpoint)

	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to API: %w", err)
	}

	return conn, nil
}

// parseAPIEndpoint ensures the endpoint is in a format grpc.NewClient understands.
func parseAPIEndpoint(raw string) string {
	if len(raw) == 0 {
		return raw
	}

	if u, err := url.Parse(raw); err == nil && u.Host != "" {
		host := u.Host
		if u.Scheme == "https" && !strings.Contains(host, ":") {
			host += ":443"
		}

		return host
	}

	if strings.HasPrefix(raw, "http://") {
		host := strings.TrimPrefix(raw, "http://")

		return host
	}

	if strings.HasPrefix(raw, "https://") {
		host := strings.TrimPrefix(raw, "https://")
		if !strings.Contains(host, ":") {
			host += ":443"
		}

		return host
	}

	return raw
}

// --- Utility Functions ---

// IsStreamEnd returns true if the error indicates the end of a gRPC stream.
func IsStreamEnd(err error) bool {
	return errors.Is(err, io.EOF)
}
