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
	"github.com/tcncloud/sati-go/pkg/ports"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func ListHuntGroupPauseCodesCmd(configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "list-hunt-group-pause-codes",
		Short: "Call GateService.ListHuntGroupPauseCodes",
		RunE: func(cmd *cobra.Command, args []string) error {
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
			params := ports.ListHuntGroupPauseCodesParams{}

			// Call the client method
			resp, err := client.ListHuntGroupPauseCodes(ctx, params)
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
				// Human-readable text output
				fmt.Printf("Hunt Group Pause Codes:\n")
				if len(resp.PauseCodes) == 0 {
					fmt.Println("  No pause codes found")
				} else {
					for _, pauseCode := range resp.PauseCodes {
						fmt.Printf("  - %s", pauseCode.Code)
						if pauseCode.Description != "" {
							fmt.Printf(": %s", pauseCode.Description)
						}
						if pauseCode.Duration > 0 {
							fmt.Printf(" (Duration: %d minutes)", pauseCode.Duration)
						}
						fmt.Println()
					}
				}
			}

			return nil
		},
	}
}
