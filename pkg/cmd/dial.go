package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
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

			// Build the custom Params struct
			params := saticlient.DialParams{
				PartnerAgentID: partnerAgentId,
				PhoneNumber:    phoneNumber,
			}
			if callerId != "" {
				params.CallerID = &callerId
			}
			if poolId != "" {
				params.PoolID = &poolId
			}
			if recordId != "" {
				params.RecordID = &recordId
			}

			// Call the client method with custom Params
			resp, err := client.Dial(ctx, params)
			if err != nil {
				return err
			}
			// Use the custom Result struct field
			fmt.Printf("Call SID: %s\n", resp.CallSid)
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
