package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tcncloud/sati-go/pkg/ports"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func DialCmd(configPath *string) *cobra.Command {
	var partnerAgentID, phoneNumber, callerID, poolID, recordID string

	cmd := &cobra.Command{
		Use:   "dial",
		Short: "Call GateService.Dial",
		RunE: func(cmd *cobra.Command, args []string) error {
			if partnerAgentID == "" || phoneNumber == "" {
				return ErrRequiredFieldsMissing
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
			defer handleClientClose(client)

			ctx, cancel := createContext(DefaultTimeout)
			defer cancel()

			// Build the custom Params struct
			params := ports.DialParams{
				PartnerAgentID: partnerAgentID,
				PhoneNumber:    phoneNumber,
			}
			if callerID != "" {
				params.CallerID = &callerID
			}
			if poolID != "" {
				params.PoolID = &poolID
			}
			if recordID != "" {
				params.RecordID = &recordID
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
	cmd.Flags().StringVar(&partnerAgentID, "partner-agent-id", "", "Partner Agent ID (required)")
	cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Phone Number (required)")
	cmd.Flags().StringVar(&callerID, "caller-id", "", "Caller ID (optional)")
	cmd.Flags().StringVar(&poolID, "pool-id", "", "Pool ID (optional)")
	cmd.Flags().StringVar(&recordID, "record-id", "", "Record ID (optional)")
	markFlagRequired(cmd, "partner-agent-id")
	markFlagRequired(cmd, "phone-number")

	return cmd
}
