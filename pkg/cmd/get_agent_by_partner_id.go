// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//
// Copyright 2024 TCN Inc

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func GetAgentByPartnerIDCmd(configPath *string) *cobra.Command {
	var partnerAgentID string

	cmd := &cobra.Command{
		Use:   "get-agent-by-partner-id",
		Short: "Call GateService.GetAgentByPartnerId",
		RunE: func(cmd *cobra.Command, args []string) error {
			if partnerAgentID == "" {
				return ErrPartnerAgentIDRequired
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
			defer handleClientClose(client) // Ensure connection is closed

			ctx, cancel := createContext(DefaultTimeout)
			defer cancel()

			// Build the params struct
			params := saticlient.GetAgentByPartnerIDParams{
				PartnerAgentID: partnerAgentID,
			}

			// Call the client method
			resp, err := client.GetAgentByPartnerID(ctx, params)
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
				if resp.Agent != nil {
					fmt.Printf("Agent found:\n")
					fmt.Printf("  UserID: %s\n", resp.Agent.UserID)
					fmt.Printf("  OrgID: %s\n", resp.Agent.OrgID)
					fmt.Printf("  Username: %s\n", resp.Agent.Username)
					fmt.Printf("  PartnerAgentID: %s\n", resp.Agent.PartnerAgentID)
					fmt.Printf("  FirstName: %s\n", resp.Agent.FirstName)
					fmt.Printf("  LastName: %s\n", resp.Agent.LastName)
				} else {
					fmt.Println("Agent not found")
				}
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&partnerAgentID, "partner-agent-id", "", "Partner Agent ID (required)")
	markFlagRequired(cmd, "partner-agent-id")

	return cmd
}
