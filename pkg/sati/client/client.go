package client

import (
	"context"

	gatev2 "github.com/tcncloud/sati-go/internal/genproto/tcnapi/exile/gate/v2"
	"github.com/tcncloud/sati-go/pkg/sati"
	"google.golang.org/grpc"
)

// Client provides methods for interacting with the GateService API.
type Client struct {
	conn *grpc.ClientConn
	gate gatev2.GateServiceClient
}

// NewClient creates a new Sati API client.
// It takes the configuration and sets up the gRPC connection and client stub.
func NewClient(cfg *sati.Config) (*Client, error) {
	conn, err := sati.SetupClient(cfg) // Reuse existing setup logic
	if err != nil {
		return nil, err
	}
	return &Client{
		conn: conn,
		gate: gatev2.NewGateServiceClient(conn),
	}, nil
}

// Close terminates the underlying gRPC connection.
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// --- GateService Methods ---

func (c *Client) AddAgentCallResponse(ctx context.Context, req *gatev2.AddAgentCallResponseRequest) (*gatev2.AddAgentCallResponseResponse, error) {
	return c.gate.AddAgentCallResponse(ctx, req)
}

func (c *Client) AddScrubListEntries(ctx context.Context, req *gatev2.AddScrubListEntriesRequest) (*gatev2.AddScrubListEntriesResponse, error) {
	return c.gate.AddScrubListEntries(ctx, req)
}

func (c *Client) Dial(ctx context.Context, req *gatev2.DialRequest) (*gatev2.DialResponse, error) {
	return c.gate.Dial(ctx, req)
}

func (c *Client) GetAgentById(ctx context.Context, req *gatev2.GetAgentByIdRequest) (*gatev2.GetAgentByIdResponse, error) {
	return c.gate.GetAgentById(ctx, req)
}

func (c *Client) GetAgentByPartnerId(ctx context.Context, req *gatev2.GetAgentByPartnerIdRequest) (*gatev2.GetAgentByPartnerIdResponse, error) {
	return c.gate.GetAgentByPartnerId(ctx, req)
}

func (c *Client) GetAgentStatus(ctx context.Context, req *gatev2.GetAgentStatusRequest) (*gatev2.GetAgentStatusResponse, error) {
	return c.gate.GetAgentStatus(ctx, req)
}

func (c *Client) GetClientConfiguration(ctx context.Context, req *gatev2.GetClientConfigurationRequest) (*gatev2.GetClientConfigurationResponse, error) {
	return c.gate.GetClientConfiguration(ctx, req)
}

func (c *Client) GetOrganizationInfo(ctx context.Context, req *gatev2.GetOrganizationInfoRequest) (*gatev2.GetOrganizationInfoResponse, error) {
	return c.gate.GetOrganizationInfo(ctx, req)
}

func (c *Client) GetRecordingStatus(ctx context.Context, req *gatev2.GetRecordingStatusRequest) (*gatev2.GetRecordingStatusResponse, error) {
	return c.gate.GetRecordingStatus(ctx, req)
}

// ListAgents returns a stream of agents.
func (c *Client) ListAgents(ctx context.Context, req *gatev2.ListAgentsRequest) (gatev2.GateService_ListAgentsClient, error) {
	return c.gate.ListAgents(ctx, req)
}

func (c *Client) ListHuntGroupPauseCodes(ctx context.Context, req *gatev2.ListHuntGroupPauseCodesRequest) (*gatev2.ListHuntGroupPauseCodesResponse, error) {
	return c.gate.ListHuntGroupPauseCodes(ctx, req)
}

func (c *Client) ListScrubLists(ctx context.Context, req *gatev2.ListScrubListsRequest) (*gatev2.ListScrubListsResponse, error) {
	return c.gate.ListScrubLists(ctx, req)
}

func (c *Client) Log(ctx context.Context, req *gatev2.LogRequest) (*gatev2.LogResponse, error) {
	return c.gate.Log(ctx, req)
}

func (c *Client) PollEvents(ctx context.Context, req *gatev2.PollEventsRequest) (*gatev2.PollEventsResponse, error) {
	return c.gate.PollEvents(ctx, req)
}

func (c *Client) PutCallOnSimpleHold(ctx context.Context, req *gatev2.PutCallOnSimpleHoldRequest) (*gatev2.PutCallOnSimpleHoldResponse, error) {
	return c.gate.PutCallOnSimpleHold(ctx, req)
}

func (c *Client) RemoveScrubListEntries(ctx context.Context, req *gatev2.RemoveScrubListEntriesRequest) (*gatev2.RemoveScrubListEntriesResponse, error) {
	return c.gate.RemoveScrubListEntries(ctx, req)
}

func (c *Client) RotateCertificate(ctx context.Context, req *gatev2.RotateCertificateRequest) (*gatev2.RotateCertificateResponse, error) {
	return c.gate.RotateCertificate(ctx, req)
}

func (c *Client) StartCallRecording(ctx context.Context, req *gatev2.StartCallRecordingRequest) (*gatev2.StartCallRecordingResponse, error) {
	return c.gate.StartCallRecording(ctx, req)
}

func (c *Client) StopCallRecording(ctx context.Context, req *gatev2.StopCallRecordingRequest) (*gatev2.StopCallRecordingResponse, error) {
	return c.gate.StopCallRecording(ctx, req)
}

// StreamJobs returns a stream of job messages.
func (c *Client) StreamJobs(ctx context.Context, req *gatev2.StreamJobsRequest) (gatev2.GateService_StreamJobsClient, error) {
	return c.gate.StreamJobs(ctx, req)
}

func (c *Client) SubmitJobResults(ctx context.Context, req *gatev2.SubmitJobResultsRequest) (*gatev2.SubmitJobResultsResponse, error) {
	return c.gate.SubmitJobResults(ctx, req)
}

func (c *Client) TakeCallOffSimpleHold(ctx context.Context, req *gatev2.TakeCallOffSimpleHoldRequest) (*gatev2.TakeCallOffSimpleHoldResponse, error) {
	return c.gate.TakeCallOffSimpleHold(ctx, req)
}

func (c *Client) UpdateAgentStatus(ctx context.Context, req *gatev2.UpdateAgentStatusRequest) (*gatev2.UpdateAgentStatusResponse, error) {
	return c.gate.UpdateAgentStatus(ctx, req)
}

func (c *Client) UpdateScrubListEntry(ctx context.Context, req *gatev2.UpdateScrubListEntryRequest) (*gatev2.UpdateScrubListEntryResponse, error) {
	return c.gate.UpdateScrubListEntry(ctx, req)
}

func (c *Client) UpsertAgent(ctx context.Context, req *gatev2.UpsertAgentRequest) (*gatev2.UpsertAgentResponse, error) {
	return c.gate.UpsertAgent(ctx, req)
}
