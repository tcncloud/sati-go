package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
			defer client.Close() // Ensure connection is closed

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Build the custom Params struct (empty in this case)
			params := saticlient.GetClientConfigurationParams{}

			// Call the client method with custom Params
			resp, err := client.GetClientConfiguration(ctx, params)
			if err != nil {
				return err
			}
			// Use the custom Result struct
			if OutputFormat == "json" {
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
