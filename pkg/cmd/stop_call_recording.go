package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func StopCallRecordingCmd(configPath *string) *cobra.Command {
	var partnerAgentID string

	cmd := &cobra.Command{
		Use:   "stop-call-recording",
		Short: "Call GateService.StopCallRecording",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := saticonfig.LoadConfig(*configPath)
			if err != nil {
				return err
			}

			client, err := saticlient.NewClient(cfg)
			if err != nil {
				return err
			}
			defer handleClientClose(client)

			ctx, cancel := createContext(DefaultTimeout)
			defer cancel()

			params := saticlient.StopCallRecordingParams{
				PartnerAgentID: partnerAgentID,
			}
			resp, err := client.StopCallRecording(ctx, params)
			if err != nil {
				return err
			}
			fmt.Printf("%+v\n", resp)

			return nil
		},
	}
	cmd.Flags().StringVar(&partnerAgentID, "partner-agent-id", "", "Partner Agent ID (required)")
	markFlagRequired(cmd, "partner-agent-id")

	return cmd
}
