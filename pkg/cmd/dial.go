package cmd

import (
	"context"
	"fmt"
	"time"

	gatev2 "buf.build/gen/go/tcn/exileapi/protocolbuffers/go/tcnapi/exile/gate/v2"
	"github.com/spf13/cobra"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func DialCmd(configPath *string) *cobra.Command {
	var partnerAgentId, phoneNumber, callerId, poolId, recordId string
	cmd := &cobra.Command{
		Use:   "dial",
		Short: "Call GateService.Dial",
		RunE: func(cmd *cobra.Command, args []string) error {
			if partnerAgentId == "" || phoneNumber == "" {
				return fmt.Errorf("--partner-agent-id and --phone-number are required")
			}
			cfg, err := saticonfig.LoadConfig(*configPath)
			if err != nil {
				return err
			}

			// Use the new client constructor
			client, err := saticlient.NewClient(cfg)
			if err != nil {
				return err
			}
			defer client.Close() // Ensure connection is closed

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Build the request struct
			request := &gatev2.DialRequest{
				PartnerAgentId: partnerAgentId,
				PhoneNumber:    phoneNumber,
			}
			if callerId != "" {
				request.CallerId = &wrapperspb.StringValue{Value: callerId}
			}
			if poolId != "" {
				request.PoolId = &wrapperspb.StringValue{Value: poolId}
			}
			if recordId != "" {
				request.RecordId = &wrapperspb.StringValue{Value: recordId}
			}

			// Call the client method
			resp, err := client.Dial(ctx, request)
			if err != nil {
				return err
			}
			fmt.Printf("%+v\n", resp)
			return nil
		},
	}
	cmd.Flags().StringVar(&partnerAgentId, "partner-agent-id", "", "Partner Agent ID (required)")
	cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Phone Number (required)")
	cmd.Flags().StringVar(&callerId, "caller-id", "", "Caller ID (optional)")
	cmd.Flags().StringVar(&poolId, "pool-id", "", "Pool ID (optional)")
	cmd.Flags().StringVar(&recordId, "record-id", "", "Record ID (optional)")
	cmd.MarkFlagRequired("partner-agent-id")
	cmd.MarkFlagRequired("phone-number")
	return cmd
}
