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
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func ListSearchableRecordingFieldsCmd(configPath *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-searchable-recording-fields",
		Short: "Call GateService.ListSearchableRecordingFields",
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
			defer client.Close() // Ensure connection is closed

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Build the custom Params struct
			params := saticlient.ListSearchableRecordingFieldsParams{}

			// Call the client method with custom Params
			resp, err := client.ListSearchableRecordingFields(ctx, params)
			if err != nil {
				return err
			}

			// Use the custom Result struct
			if OutputFormat == "json" {
				data, err := json.MarshalIndent(resp.Fields, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
			} else {
				for _, field := range resp.Fields {
					fmt.Printf("Name: %s, DisplayName: %s, Type: %s\n",
						field.Name, field.DisplayName, field.Type)
				}
			}
			return nil
		},
	}
	return cmd
}
