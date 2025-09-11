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
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func UpsertAgentCmd(configPath *string) *cobra.Command {
	var username, partnerAgentId, firstName, lastName, password string
	cmd := &cobra.Command{
		Use:   "upsert-agent",
		Short: "Call GateService.UpsertAgent",
		RunE: func(cmd *cobra.Command, args []string) error {
			if username == "" || partnerAgentId == "" || firstName == "" || lastName == "" || password == "" {
				return fmt.Errorf("--username, --partner-agent-id, --first-name, --last-name, and --password are required")
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

			// Build the request struct
			request := &gatev2.UpsertAgentRequest{
				Username:       username,
				PartnerAgentId: partnerAgentId,
				FirstName:      firstName,
				LastName:       lastName,
				Password:       password,
			}

			// Call the client method
			resp, err := client.UpsertAgent(ctx, request)
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
	cmd.Flags().StringVar(&username, "username", "", "Username (required)")
	cmd.Flags().StringVar(&partnerAgentId, "partner-agent-id", "", "Partner Agent ID (required)")
	cmd.Flags().StringVar(&firstName, "first-name", "", "First Name (required)")
	cmd.Flags().StringVar(&lastName, "last-name", "", "Last Name (required)")
	cmd.Flags().StringVar(&password, "password", "", "Password (required)")
	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("partner-agent-id")
	cmd.MarkFlagRequired("first-name")
	cmd.MarkFlagRequired("last-name")
	cmd.MarkFlagRequired("password")
	return cmd
}
