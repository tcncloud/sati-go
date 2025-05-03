package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/url"
	"strings"

	gatev2grpc "buf.build/gen/go/tcn/exileapi/grpc/go/tcnapi/exile/gate/v2/gatev2grpc"
	gatev2pb "buf.build/gen/go/tcn/exileapi/protocolbuffers/go/tcnapi/exile/gate/v2"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client provides methods for interacting with the GateService API.
type Client struct {
	conn *grpc.ClientConn
	gate gatev2grpc.GateServiceClient
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
		gate: gatev2grpc.NewGateServiceClient(conn),
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

func (c *Client) AddAgentCallResponse(ctx context.Context, req *gatev2pb.AddAgentCallResponseRequest) (*gatev2pb.AddAgentCallResponseResponse, error) {
	return c.gate.AddAgentCallResponse(ctx, req)
}

func (c *Client) AddScrubListEntries(ctx context.Context, req *gatev2pb.AddScrubListEntriesRequest) (*gatev2pb.AddScrubListEntriesResponse, error) {
	return c.gate.AddScrubListEntries(ctx, req)
}

func (c *Client) Dial(ctx context.Context, req *gatev2pb.DialRequest) (*gatev2pb.DialResponse, error) {
	return c.gate.Dial(ctx, req)
}

func (c *Client) GetAgentById(ctx context.Context, req *gatev2pb.GetAgentByIdRequest) (*gatev2pb.GetAgentByIdResponse, error) {
	return c.gate.GetAgentById(ctx, req)
}

func (c *Client) GetAgentByPartnerId(ctx context.Context, req *gatev2pb.GetAgentByPartnerIdRequest) (*gatev2pb.GetAgentByPartnerIdResponse, error) {
	return c.gate.GetAgentByPartnerId(ctx, req)
}

func (c *Client) GetAgentStatus(ctx context.Context, req *gatev2pb.GetAgentStatusRequest) (*gatev2pb.GetAgentStatusResponse, error) {
	return c.gate.GetAgentStatus(ctx, req)
}

func (c *Client) GetClientConfiguration(ctx context.Context, req *gatev2pb.GetClientConfigurationRequest) (*gatev2pb.GetClientConfigurationResponse, error) {
	return c.gate.GetClientConfiguration(ctx, req)
}

func (c *Client) GetOrganizationInfo(ctx context.Context, req *gatev2pb.GetOrganizationInfoRequest) (*gatev2pb.GetOrganizationInfoResponse, error) {
	return c.gate.GetOrganizationInfo(ctx, req)
}

func (c *Client) GetRecordingStatus(ctx context.Context, req *gatev2pb.GetRecordingStatusRequest) (*gatev2pb.GetRecordingStatusResponse, error) {
	return c.gate.GetRecordingStatus(ctx, req)
}

// ListAgents returns a stream of agents.
func (c *Client) ListAgents(ctx context.Context, req *gatev2pb.ListAgentsRequest) (gatev2grpc.GateService_ListAgentsClient, error) {
	return c.gate.ListAgents(ctx, req)
}

func (c *Client) ListHuntGroupPauseCodes(ctx context.Context, req *gatev2pb.ListHuntGroupPauseCodesRequest) (*gatev2pb.ListHuntGroupPauseCodesResponse, error) {
	return c.gate.ListHuntGroupPauseCodes(ctx, req)
}

func (c *Client) ListScrubLists(ctx context.Context, req *gatev2pb.ListScrubListsRequest) (*gatev2pb.ListScrubListsResponse, error) {
	return c.gate.ListScrubLists(ctx, req)
}

func (c *Client) Log(ctx context.Context, req *gatev2pb.LogRequest) (*gatev2pb.LogResponse, error) {
	return c.gate.Log(ctx, req)
}

func (c *Client) PollEvents(ctx context.Context, req *gatev2pb.PollEventsRequest) (*gatev2pb.PollEventsResponse, error) {
	return c.gate.PollEvents(ctx, req)
}

func (c *Client) PutCallOnSimpleHold(ctx context.Context, req *gatev2pb.PutCallOnSimpleHoldRequest) (*gatev2pb.PutCallOnSimpleHoldResponse, error) {
	return c.gate.PutCallOnSimpleHold(ctx, req)
}

func (c *Client) RemoveScrubListEntries(ctx context.Context, req *gatev2pb.RemoveScrubListEntriesRequest) (*gatev2pb.RemoveScrubListEntriesResponse, error) {
	return c.gate.RemoveScrubListEntries(ctx, req)
}

func (c *Client) RotateCertificate(ctx context.Context, req *gatev2pb.RotateCertificateRequest) (*gatev2pb.RotateCertificateResponse, error) {
	return c.gate.RotateCertificate(ctx, req)
}

func (c *Client) StartCallRecording(ctx context.Context, req *gatev2pb.StartCallRecordingRequest) (*gatev2pb.StartCallRecordingResponse, error) {
	return c.gate.StartCallRecording(ctx, req)
}

func (c *Client) StopCallRecording(ctx context.Context, req *gatev2pb.StopCallRecordingRequest) (*gatev2pb.StopCallRecordingResponse, error) {
	return c.gate.StopCallRecording(ctx, req)
}

// StreamJobs returns a stream of job messages.
func (c *Client) StreamJobs(ctx context.Context, req *gatev2pb.StreamJobsRequest) (gatev2grpc.GateService_StreamJobsClient, error) {
	return c.gate.StreamJobs(ctx, req)
}

func (c *Client) SubmitJobResults(ctx context.Context, req *gatev2pb.SubmitJobResultsRequest) (*gatev2pb.SubmitJobResultsResponse, error) {
	return c.gate.SubmitJobResults(ctx, req)
}

func (c *Client) TakeCallOffSimpleHold(ctx context.Context, req *gatev2pb.TakeCallOffSimpleHoldRequest) (*gatev2pb.TakeCallOffSimpleHoldResponse, error) {
	return c.gate.TakeCallOffSimpleHold(ctx, req)
}

func (c *Client) UpdateAgentStatus(ctx context.Context, req *gatev2pb.UpdateAgentStatusRequest) (*gatev2pb.UpdateAgentStatusResponse, error) {
	return c.gate.UpdateAgentStatus(ctx, req)
}

func (c *Client) UpdateScrubListEntry(ctx context.Context, req *gatev2pb.UpdateScrubListEntryRequest) (*gatev2pb.UpdateScrubListEntryResponse, error) {
	return c.gate.UpdateScrubListEntry(ctx, req)
}

func (c *Client) UpsertAgent(ctx context.Context, req *gatev2pb.UpsertAgentRequest) (*gatev2pb.UpsertAgentResponse, error) {
	return c.gate.UpsertAgent(ctx, req)
}

// --- Internal Helper Functions ---

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
