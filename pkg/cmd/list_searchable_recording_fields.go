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
	"github.com/tcncloud/sati-go/pkg/ports"
)

func ListSearchableRecordingFieldsCmd(configPath *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-searchable-recording-fields",
		Short: "Call GateService.ListSearchableRecordingFields",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := createClient(configPath)
			if err != nil {
				return err
			}
			defer handleClientClose(client)

			ctx, cancel := createContext(DefaultTimeout)
			defer cancel()

			params := ports.ListSearchableRecordingFieldsParams{}
			resp, err := client.ListSearchableRecordingFields(ctx, params)
			if err != nil {
				return fmt.Errorf("failed to list searchable recording fields: %w", err)
			}

			if OutputFormat == OutputFormatJSON {
				return outputJSON(resp.Fields)
			}

			for _, field := range resp.Fields {
				fmt.Printf("Name: %s, DisplayName: %s, Type: %s\n",
					field.Name, field.DisplayName, field.Type)
			}

			return nil
		},
	}

	return cmd
}
