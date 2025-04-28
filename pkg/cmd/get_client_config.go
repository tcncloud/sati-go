package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	gatev2 "git.tcncloud.net/experiments/sati-client/internal/genproto/tcnapi/exile/gate/v2"
	"git.tcncloud.net/experiments/sati-client/pkg/sati"
	"github.com/spf13/cobra"
)

func GetClientConfigCmd(configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "get-client-config",
		Short: "Call GateService.GetClientConfiguration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := sati.LoadConfig(*configPath)
			if err != nil {
				return err
			}
			conn, err := sati.SetupClient(cfg)
			if err != nil {
				return err
			}
			defer conn.Close()
			client := gatev2.NewGateServiceClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			resp, err := client.GetClientConfiguration(ctx, &gatev2.GetClientConfigurationRequest{})
			if err != nil {
				return err
			}
			if OutputFormat == "json" {
				data, err := json.MarshalIndent(resp, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
			} else {
				fmt.Printf("OrgID: %s\nOrgName: %s\nConfigName: %s\nConfigPayload: %s\n", resp.GetOrgId(), resp.GetOrgName(), resp.GetConfigName(), resp.GetConfigPayload())
			}
			return nil
		},
	}
}
