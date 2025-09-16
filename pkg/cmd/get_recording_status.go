package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func GetRecordingStatusCmd(configPath *string) *cobra.Command {
	var partnerAgentID string

	cmd := &cobra.Command{
		Use:   "get-recording-status",
		Short: "Call GateService.GetRecordingStatus",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := saticonfig.LoadConfig(*configPath)
			if err != nil {
				return err
			}

			// Use the new client constructor
			client, err := saticlient.NewClient(cfg)
			if err != nil {
				return err
			}
			defer handleClientClose(client) // Ensure connection is closed

			ctx, cancel := createContext(DefaultTimeout)
			defer cancel()

			// Build the params struct
			params := saticlient.GetRecordingStatusParams{
				PartnerAgentID: partnerAgentID,
			}

			// Call the client method
			resp, err := client.GetRecordingStatus(ctx, params)
			if err != nil {
				return err
			}
			if OutputFormat == OutputFormatJSON {
				data, err := json.MarshalIndent(resp, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
			} else {
				fmt.Printf("Recording Status: %s\n", resp.Status)
				if resp.RecordingSid != "" {
					fmt.Printf("Recording SID: %s\n", resp.RecordingSid)
				}
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&partnerAgentID, "partner-agent-id", "", "Partner Agent ID (required)")
	markFlagRequired(cmd, "partner-agent-id")

	return cmd
}
