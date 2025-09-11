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

func GetAgentByIDCmd(configPath *string) *cobra.Command {
	var userID string

	cmd := &cobra.Command{
		Use:   "get-agent-by-id",
		Short: "Call GateService.GetAgentById",
		RunE: func(cmd *cobra.Command, args []string) error {
			if userID == "" {
				return ErrUserIDRequired
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

			// Build the custom Params struct
			params := saticlient.GetAgentByIDParams{UserID: userID}

			// Call the client method with custom Params
			resp, err := client.GetAgentByID(ctx, params)
			if err != nil {
				return err
			}
			// Use the custom Result struct
			if OutputFormat == OutputFormatJSON {
				data, err := json.MarshalIndent(resp.Agent, "", "  ") // Marshal the Agent struct directly
				if err != nil {
					return err
				}
				fmt.Println(string(data))
			} else {
				if resp.Agent != nil {
					fmt.Printf("UserID: %s, OrgID: %s, FirstName: %s, LastName: %s, Username: %s, PartnerAgentID: %s\n",
						resp.Agent.UserID, resp.Agent.OrgID, resp.Agent.FirstName, resp.Agent.LastName, resp.Agent.Username, resp.Agent.PartnerAgentID)
				} else {
					fmt.Println("Agent not found.")
				}
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&userID, "user-id", "", "User ID (required)")
	markFlagRequired(cmd, "user-id")

	return cmd
}
