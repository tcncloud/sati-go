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
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func ListAgentsCmd(configPath *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-agents",
		Short: "Call GateService.ListAgents",
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

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			// Build the custom Params struct (empty)
			params := saticlient.ListAgentsParams{}

			// Call the client stream method - returns a channel
			resultsChan := client.ListAgents(ctx, params)

			var agents []*saticlient.Agent // Store results from channel
			for result := range resultsChan {
				if result.Error != nil {
					// Handle potential errors during streaming
					// Depending on desired behavior, could log, return, or collect errors
					return fmt.Errorf("error streaming agents: %w", result.Error)
				}
				if result.Agent != nil {
					agents = append(agents, result.Agent)
				}
			}

			// Process collected agents
			if OutputFormat == "json" {
				data, err := json.MarshalIndent(agents, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
			} else {
				for _, agent := range agents {
					fmt.Printf("UserID: %s, OrgID: %s, FirstName: %s, LastName: %s, Username: %s, PartnerAgentID: %s\n",
						agent.UserID, agent.OrgID, agent.FirstName, agent.LastName, agent.Username, agent.PartnerAgentID)
				}
			}
			return nil
		},
	}
	return cmd
}
