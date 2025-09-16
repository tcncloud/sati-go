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

func UpsertAgentCmd(configPath *string) *cobra.Command {
	var username, partnerAgentID, firstName, lastName, password string

	cmd := &cobra.Command{
		Use:   "upsert-agent",
		Short: "Call GateService.UpsertAgent",
		RunE: func(cmd *cobra.Command, args []string) error {
			if username == "" || partnerAgentID == "" || firstName == "" || lastName == "" || password == "" {
				return ErrRequiredFieldsMissing
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
			params := saticlient.UpsertAgentParams{
				Username:       username,
				PartnerAgentID: partnerAgentID,
				FirstName:      firstName,
				LastName:       lastName,
			}

			// Call the client method
			resp, err := client.UpsertAgent(ctx, params)
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
				fmt.Println("Agent upserted successfully")
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&username, "username", "", "Username (required)")
	cmd.Flags().StringVar(&partnerAgentID, "partner-agent-id", "", "Partner Agent ID (required)")
	cmd.Flags().StringVar(&firstName, "first-name", "", "First Name (required)")
	cmd.Flags().StringVar(&lastName, "last-name", "", "Last Name (required)")
	cmd.Flags().StringVar(&password, "password", "", "Password (required)")
	markFlagRequired(cmd, "username")
	markFlagRequired(cmd, "partner-agent-id")
	markFlagRequired(cmd, "first-name")
	markFlagRequired(cmd, "last-name")
	markFlagRequired(cmd, "password")

	return cmd
}
