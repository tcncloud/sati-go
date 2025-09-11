package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func GetClientConfigCmd(configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "get-client-config",
		Short: "Call GateService.GetClientConfiguration",
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

			// Build the custom Params struct (empty in this case)
			params := saticlient.GetClientConfigurationParams{}

			// Call the client method with custom Params
			resp, err := client.GetClientConfiguration(ctx, params)
			if err != nil {
				return err
			}
			// Use the custom Result struct
			if OutputFormat == OutputFormatJSON {
				data, err := json.MarshalIndent(resp, "", "  ") // Marshal the Result struct
				if err != nil {
					return err
				}
				fmt.Println(string(data))
			} else {
				fmt.Printf("OrgID: %s\nOrgName: %s\nConfigName: %s\nConfigPayload: %s\n",
					resp.OrgID, resp.OrgName, resp.ConfigName, resp.ConfigPayload) // Use direct fields
			}

			return nil
		},
	}
}
