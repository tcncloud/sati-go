package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/url"
	"strings"

	gatev2pb "github.com/tcncloud/sati-go/internal/genproto/tcnapi/exile/gate/v2" // Keep for internal mapping
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/wrapperspb" // Needed for optional fields
)

// Client provides methods for interacting with the GateService API.
type Client struct {
	conn *grpc.ClientConn
	gate gatev2pb.GateServiceClient
}

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

// AddAgentCallResponse remains unchanged for now as its types weren't defined in types.go
func (c *Client) AddAgentCallResponse(ctx context.Context, req *gatev2pb.AddAgentCallResponseRequest) (*gatev2pb.AddAgentCallResponseResponse, error) {
	return c.gate.AddAgentCallResponse(ctx, req)
}

// AddScrubListEntries adds entries to a scrub list.
func (c *Client) AddScrubListEntries(ctx context.Context, params AddScrubListEntriesParams) (AddScrubListEntriesResult, error) {
	if params.ScrubListID == "" || len(params.Entries) == 0 {
		return AddScrubListEntriesResult{}, fmt.Errorf("ScrubListID and at least one Entry are required")
	}

	pbEntries := make([]*gatev2pb.AddScrubListEntriesRequest_Entry, 0, len(params.Entries))
	for _, e := range params.Entries {
		if e.Content == "" {
			return AddScrubListEntriesResult{}, fmt.Errorf("entry content cannot be empty")
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
		return AddScrubListEntriesResult{}, err
	}

	return AddScrubListEntriesResult{}, nil // Return empty struct for success
}

// Dial initiates an outbound call.
func (c *Client) Dial(ctx context.Context, params DialParams) (DialResult, error) {
	if params.PartnerAgentID == "" || params.PhoneNumber == "" {
		return DialResult{}, fmt.Errorf("PartnerAgentID and PhoneNumber are required")
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
		return DialResult{}, err
	}
	if resp == nil {
		return DialResult{}, fmt.Errorf("received nil response from gRPC Dial")
	}

	result := DialResult{
		CallSid: resp.GetCallSid(),
	}
	return result, nil
}

// GetAgentById retrieves agent details by User ID.
func (c *Client) GetAgentById(ctx context.Context, params GetAgentByIdParams) (GetAgentByIdResult, error) {
	if params.UserID == "" {
		return GetAgentByIdResult{}, fmt.Errorf("UserID is required")
	}
	req := &gatev2pb.GetAgentByIdRequest{UserId: params.UserID}

	resp, err := c.gate.GetAgentById(ctx, req)
	if err != nil {
		return GetAgentByIdResult{}, err
	}
	if resp == nil || resp.Agent == nil {
		// Consider returning a specific "not found" error type here
		return GetAgentByIdResult{}, fmt.Errorf("agent not found or nil response received")
	}

	result := GetAgentByIdResult{
		Agent: mapProtoAgentToAgent(resp.Agent), // Use mapping helper
	}
	return result, nil
}

// GetAgentByPartnerId remains unchanged for now
func (c *Client) GetAgentByPartnerId(ctx context.Context, req *gatev2pb.GetAgentByPartnerIdRequest) (*gatev2pb.GetAgentByPartnerIdResponse, error) {
	return c.gate.GetAgentByPartnerId(ctx, req)
}

// GetAgentStatus remains unchanged for now
func (c *Client) GetAgentStatus(ctx context.Context, req *gatev2pb.GetAgentStatusRequest) (*gatev2pb.GetAgentStatusResponse, error) {
	return c.gate.GetAgentStatus(ctx, req)
}

// GetClientConfiguration retrieves client configuration details.
func (c *Client) GetClientConfiguration(ctx context.Context, params GetClientConfigurationParams) (GetClientConfigurationResult, error) {
	req := &gatev2pb.GetClientConfigurationRequest{}
	resp, err := c.gate.GetClientConfiguration(ctx, req)
	if err != nil {
		return GetClientConfigurationResult{}, err
	}
	if resp == nil {
		return GetClientConfigurationResult{}, fmt.Errorf("received nil response from gRPC GetClientConfiguration")
	}

	result := GetClientConfigurationResult{
		OrgID:         resp.GetOrgId(),
		OrgName:       resp.GetOrgName(),
		ConfigName:    resp.GetConfigName(),
		ConfigPayload: resp.GetConfigPayload(),
	}
	return result, nil
}

// GetOrganizationInfo remains unchanged for now
func (c *Client) GetOrganizationInfo(ctx context.Context, req *gatev2pb.GetOrganizationInfoRequest) (*gatev2pb.GetOrganizationInfoResponse, error) {
	return c.gate.GetOrganizationInfo(ctx, req)
}

// GetRecordingStatus remains unchanged for now
func (c *Client) GetRecordingStatus(ctx context.Context, req *gatev2pb.GetRecordingStatusRequest) (*gatev2pb.GetRecordingStatusResponse, error) {
	return c.gate.GetRecordingStatus(ctx, req)
}

// ListAgents returns a channel that emits agent details.
func (c *Client) ListAgents(ctx context.Context, params ListAgentsParams) <-chan ListAgentsResult {
	resultsChan := make(chan ListAgentsResult)
	req := &gatev2pb.ListAgentsRequest{}

	go func() {
		defer close(resultsChan)
		stream, err := c.gate.ListAgents(ctx, req)
		if err != nil {
			resultsChan <- ListAgentsResult{Error: fmt.Errorf("failed to start ListAgents stream: %w", err)}
			return
		}

		for {
			resp, err := stream.Recv()
			if err != nil {
				if !IsStreamEnd(err) { // Don't send EOF as error
					resultsChan <- ListAgentsResult{Error: fmt.Errorf("error receiving from ListAgents stream: %w", err)}
				}
				return // End goroutine on EOF or error
			}
			if resp == nil || resp.Agent == nil {
				// Handle nil response/agent, maybe log or send specific error
				resultsChan <- ListAgentsResult{Error: fmt.Errorf("received nil agent in ListAgents stream")}
				continue
			}
			resultsChan <- ListAgentsResult{Agent: mapProtoAgentToAgent(resp.Agent)}
		}
	}()

	return resultsChan
}

// ListHuntGroupPauseCodes remains unchanged for now
func (c *Client) ListHuntGroupPauseCodes(ctx context.Context, req *gatev2pb.ListHuntGroupPauseCodesRequest) (*gatev2pb.ListHuntGroupPauseCodesResponse, error) {
	return c.gate.ListHuntGroupPauseCodes(ctx, req)
}

// ListScrubLists remains unchanged for now
func (c *Client) ListScrubLists(ctx context.Context, req *gatev2pb.ListScrubListsRequest) (*gatev2pb.ListScrubListsResponse, error) {
	return c.gate.ListScrubLists(ctx, req)
}

// Log remains unchanged for now
func (c *Client) Log(ctx context.Context, req *gatev2pb.LogRequest) (*gatev2pb.LogResponse, error) {
	return c.gate.Log(ctx, req)
}

// PollEvents remains unchanged for now
func (c *Client) PollEvents(ctx context.Context, req *gatev2pb.PollEventsRequest) (*gatev2pb.PollEventsResponse, error) {
	return c.gate.PollEvents(ctx, req)
}

// PutCallOnSimpleHold remains unchanged for now
func (c *Client) PutCallOnSimpleHold(ctx context.Context, req *gatev2pb.PutCallOnSimpleHoldRequest) (*gatev2pb.PutCallOnSimpleHoldResponse, error) {
	return c.gate.PutCallOnSimpleHold(ctx, req)
}

// RemoveScrubListEntries remains unchanged for now
func (c *Client) RemoveScrubListEntries(ctx context.Context, req *gatev2pb.RemoveScrubListEntriesRequest) (*gatev2pb.RemoveScrubListEntriesResponse, error) {
	return c.gate.RemoveScrubListEntries(ctx, req)
}

// RotateCertificate remains unchanged for now
func (c *Client) RotateCertificate(ctx context.Context, req *gatev2pb.RotateCertificateRequest) (*gatev2pb.RotateCertificateResponse, error) {
	return c.gate.RotateCertificate(ctx, req)
}

// StartCallRecording remains unchanged for now
func (c *Client) StartCallRecording(ctx context.Context, req *gatev2pb.StartCallRecordingRequest) (*gatev2pb.StartCallRecordingResponse, error) {
	return c.gate.StartCallRecording(ctx, req)
}

// StopCallRecording remains unchanged for now
func (c *Client) StopCallRecording(ctx context.Context, req *gatev2pb.StopCallRecordingRequest) (*gatev2pb.StopCallRecordingResponse, error) {
	return c.gate.StopCallRecording(ctx, req)
}

// StreamJobs remains unchanged for now
func (c *Client) StreamJobs(ctx context.Context, req *gatev2pb.StreamJobsRequest) (gatev2pb.GateService_StreamJobsClient, error) {
	return c.gate.StreamJobs(ctx, req)
}

// SubmitJobResults remains unchanged for now
func (c *Client) SubmitJobResults(ctx context.Context, req *gatev2pb.SubmitJobResultsRequest) (*gatev2pb.SubmitJobResultsResponse, error) {
	return c.gate.SubmitJobResults(ctx, req)
}

// TakeCallOffSimpleHold remains unchanged for now
func (c *Client) TakeCallOffSimpleHold(ctx context.Context, req *gatev2pb.TakeCallOffSimpleHoldRequest) (*gatev2pb.TakeCallOffSimpleHoldResponse, error) {
	return c.gate.TakeCallOffSimpleHold(ctx, req)
}

// UpdateAgentStatus updates the state of an agent.
func (c *Client) UpdateAgentStatus(ctx context.Context, params UpdateAgentStatusParams) (UpdateAgentStatusResult, error) {
	if params.PartnerAgentID == "" {
		return UpdateAgentStatusResult{}, fmt.Errorf("PartnerAgentID is required")
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
		return UpdateAgentStatusResult{}, err
	}
	return UpdateAgentStatusResult{}, nil
}

// UpdateScrubListEntry remains unchanged for now
func (c *Client) UpdateScrubListEntry(ctx context.Context, req *gatev2pb.UpdateScrubListEntryRequest) (*gatev2pb.UpdateScrubListEntryResponse, error) {
	return c.gate.UpdateScrubListEntry(ctx, req)
}

// UpsertAgent remains unchanged for now
func (c *Client) UpsertAgent(ctx context.Context, req *gatev2pb.UpsertAgentRequest) (*gatev2pb.UpsertAgentResponse, error) {
	return c.gate.UpsertAgent(ctx, req)
}

// --- Internal Helper Functions ---

// mapProtoAgentToAgent converts a proto Agent to our custom client Agent type.
// TODO: Fix corev2 import path before using this.
func mapProtoAgentToAgent(pbAgent *gatev2pb.Agent) *Agent {
	if pbAgent == nil {
		return nil
	}
	return &Agent{
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
		return nil, fmt.Errorf("failed to append CA cert")
	}
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
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
	return err == io.EOF
}
