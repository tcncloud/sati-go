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
	gatev2 "github.com/tcncloud/sati-go/internal/genproto/tcnapi/exile/gate/v2"
	"github.com/tcncloud/sati-go/pkg/sati"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
)

func UpdateAgentStatusCmd(configPath *string) *cobra.Command {
	var partnerAgentId, newStateStr, reason string
	cmd := &cobra.Command{
		Use:   "update-agent-status",
		Short: "Call GateService.UpdateAgentStatus",
		RunE: func(cmd *cobra.Command, args []string) error {
			if partnerAgentId == "" || newStateStr == "" {
				return fmt.Errorf("--partner-agent-id and --new-state are required")
			}
			newStateEnum, ok := gatev2.AgentState_value[newStateStr]
			if !ok {
				return fmt.Errorf("invalid --new-state: %s", newStateStr)
			}
			cfg, err := sati.LoadConfig(*configPath)
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

			// Build the request struct
			request := &gatev2.UpdateAgentStatusRequest{
				PartnerAgentId: partnerAgentId,
				NewState:       gatev2.AgentState(newStateEnum),
				Reason:         reason,
			}

			// Call the client method
			resp, err := client.UpdateAgentStatus(ctx, request)
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
				fmt.Printf("%+v\n", resp)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&partnerAgentId, "partner-agent-id", "", "Partner Agent ID (required)")
	cmd.Flags().StringVar(&newStateStr, "new-state", "", "New State (required, e.g. AGENT_STATE_READY)")
	cmd.Flags().StringVar(&reason, "reason", "", "Reason (optional)")
	cmd.MarkFlagRequired("partner-agent-id")
	cmd.MarkFlagRequired("new-state")
	return cmd
}
