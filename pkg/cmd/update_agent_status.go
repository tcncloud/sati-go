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
	"fmt"

	"github.com/spf13/cobra"
	gatev2 "github.com/tcncloud/sati-go/internal/genproto/tcnapi/exile/gate/v2"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func UpdateAgentStatusCmd(configPath *string) *cobra.Command {
	var partnerAgentID, newStateStr, reason string

	cmd := &cobra.Command{
		Use:   "update-agent-status",
		Short: "Call GateService.UpdateAgentStatus",
		RunE: func(cmd *cobra.Command, args []string) error {
			if partnerAgentID == "" || newStateStr == "" {
				return ErrRequiredFieldsMissing
			}
			newStateEnum, ok := gatev2.AgentState_value[newStateStr]
			if !ok {
				return fmt.Errorf("%w: %s", ErrInvalidNewState, newStateStr)
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
			params := saticlient.UpdateAgentStatusParams{
				PartnerAgentID: partnerAgentID,
				NewState:       gatev2.AgentState(newStateEnum),
			}
			if reason != "" {
				params.Reason = &reason
			}

			// Call the client method with custom Params
			_, err = client.UpdateAgentStatus(ctx, params)
			if err != nil {
				return err
			}
			// Response is now an empty struct on success
			fmt.Println("Successfully updated agent status.")

			return nil
		},
	}
	cmd.Flags().StringVar(&partnerAgentID, "partner-agent-id", "", "Partner Agent ID (required)")
	cmd.Flags().StringVar(&newStateStr, "new-state", "", "New State (required, e.g. AGENT_STATE_READY)")
	cmd.Flags().StringVar(&reason, "reason", "", "Reason (optional)")
	markFlagRequired(cmd, "partner-agent-id")
	markFlagRequired(cmd, "new-state")

	return cmd
}
