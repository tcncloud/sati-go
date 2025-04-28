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

	gatev2 "git.tcncloud.net/experiments/sati-client/internal/genproto/tcnapi/exile/gate/v2"
	"git.tcncloud.net/experiments/sati-client/pkg/sati"
	"github.com/spf13/cobra"
)

func ListAgentsCmd(configPath *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-agents",
		Short: "Call GateService.ListAgents",
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
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			stream, err := client.ListAgents(ctx, &gatev2.ListAgentsRequest{})
			if err != nil {
				return err
			}

			var agents []*gatev2.Agent
			for {
				resp, err := stream.Recv()
				if err != nil {
					if sati.IsStreamEnd(err) {
						break
					}
					return err
				}
				agents = append(agents, resp.Agent)
			}

			if OutputFormat == "json" {
				data, err := json.MarshalIndent(agents, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
			} else {
				for _, agent := range agents {
					fmt.Printf("UserID: %s, OrgID: %s, FirstName: %s, LastName: %s, Username: %s, PartnerAgentID: %s\n",
						agent.UserId, agent.OrgId, agent.FirstName, agent.LastName, agent.Username, agent.PartnerAgentId)
				}
			}
			return nil
		},
	}
	return cmd
}
