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
			conn, err := sati.SetupClient(cfg)
			if err != nil {
				return err
			}
			defer conn.Close()
			client := gatev2.NewGateServiceClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			resp, err := client.UpdateAgentStatus(ctx, &gatev2.UpdateAgentStatusRequest{
				PartnerAgentId: partnerAgentId,
				NewState:       gatev2.AgentState(newStateEnum),
				Reason:         reason,
			})
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
