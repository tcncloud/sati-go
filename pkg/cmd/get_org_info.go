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

func GetOrgInfoCmd(configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "get-org-info",
		Short: "Call GateService.GetOrganizationInfo",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := sati.LoadConfig(*configPath)
			if err != nil {
				return err
			}

			// Use the new client constructor
			client, err := saticlient.NewClient(cfg)
			if err != nil {
				fmt.Printf("Error setting up client: %v\n", err)
				return err
			}
			defer client.Close() // Ensure connection is closed

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Build the request struct
			request := &gatev2.GetOrganizationInfoRequest{}

			// Call the client method
			resp, err := client.GetOrganizationInfo(ctx, request)
			if err != nil {
				fmt.Printf("Error getting organization info: %v\n", err)
				return err
			}
			if OutputFormat == "json" {
				data, err := json.MarshalIndent(resp, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
			} else {
				fmt.Printf("OrgID: %s\nOrgName: %s\n", resp.GetOrgId(), resp.GetOrgName())
			}
			return nil
		},
	}
}
