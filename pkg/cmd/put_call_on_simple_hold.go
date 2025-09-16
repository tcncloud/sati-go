package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func PutCallOnSimpleHoldCmd(configPath *string) *cobra.Command {
	var partnerAgentID string

	cmd := &cobra.Command{
		Use:   "put-call-on-simple-hold",
		Short: "Call GateService.PutCallOnSimpleHold",
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

			params := saticlient.PutCallOnSimpleHoldParams{
				PartnerAgentID: partnerAgentID,
			}
			resp, err := client.PutCallOnSimpleHold(ctx, params)
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
				fmt.Println("Call put on simple hold successfully")
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&partnerAgentID, "partner-agent-id", "", "Partner Agent ID (required)")
	markFlagRequired(cmd, "partner-agent-id")

	return cmd
}
