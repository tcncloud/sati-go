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
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
)

func ListSkillsCmd(configPath *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-skills",
		Short: "Call GateService.ListSkills",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := createClient(configPath)
			if err != nil {
				return err
			}
			defer handleClientClose(client)

			ctx, cancel := createContext(DefaultTimeout)
			defer cancel()

			params := saticlient.ListSkillsParams{}
			resp, err := client.ListSkills(ctx, params)
			if err != nil {
				return fmt.Errorf("failed to list skills: %w", err)
			}

			if OutputFormat == OutputFormatJSON {
				return outputJSON(resp.Skills)
			}

			for _, skill := range resp.Skills {
				fmt.Printf("ID: %s, Name: %s, Description: %s\n",
					skill.ID, skill.Name, skill.Description)
			}

			return nil
		},
	}

	return cmd
}
