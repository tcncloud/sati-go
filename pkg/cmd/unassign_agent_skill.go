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

	"github.com/spf13/cobra"
	"github.com/tcncloud/sati-go/pkg/ports"
)

func UnassignAgentSkillCmd(configPath *string) *cobra.Command {
	var partnerAgentID, skillID string

	return createSkillCommand(
		"unassign-agent-skill",
		"Call GateService.UnassignAgentSkill",
		configPath,
		&partnerAgentID,
		&skillID,
		func(client ports.ClientInterface, ctx context.Context, params ports.AssignAgentSkillParams) error {
			// Convert to UnassignAgentSkillParams
			unassignParams := ports.UnassignAgentSkillParams(params)

			_, err := client.UnassignAgentSkill(ctx, unassignParams)
			if err != nil {
				return fmt.Errorf("failed to unassign skill: %w", err)
			}

			return nil
		},
		"Skill unassigned successfully",
	)
}
