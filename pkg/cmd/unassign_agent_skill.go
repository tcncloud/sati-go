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
	"fmt"
	"time"

	"github.com/spf13/cobra"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func UnassignAgentSkillCmd(configPath *string) *cobra.Command {
	var partnerAgentID, skillID string
	cmd := &cobra.Command{
		Use:   "unassign-agent-skill",
		Short: "Call GateService.UnassignAgentSkill",
		RunE: func(cmd *cobra.Command, args []string) error {
			if partnerAgentID == "" {
				return fmt.Errorf("--partner-agent-id is required")
			}
			if skillID == "" {
				return fmt.Errorf("--skill-id is required")
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
			defer client.Close() // Ensure connection is closed

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Build the custom Params struct
			params := saticlient.UnassignAgentSkillParams{
				PartnerAgentID: partnerAgentID,
				SkillID:        skillID,
			}

			// Call the client method with custom Params
			_, err = client.UnassignAgentSkill(ctx, params)
			if err != nil {
				return err
			}

			fmt.Println("Skill unassigned successfully")
			return nil
		},
	}
	cmd.Flags().StringVar(&partnerAgentID, "partner-agent-id", "", "Partner Agent ID (required)")
	cmd.Flags().StringVar(&skillID, "skill-id", "", "Skill ID (required)")
	cmd.MarkFlagRequired("partner-agent-id")
	cmd.MarkFlagRequired("skill-id")
	return cmd
}
