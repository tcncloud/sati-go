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

func SubmitJobResultsCmd(configPath *string) *cobra.Command {
	var jobId, resultJSON string
	var endOfTransmission bool
	cmd := &cobra.Command{
		Use:   "submit-job-results",
		Short: "Call GateService.SubmitJobResults",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := sati.LoadConfig(*configPath)
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
			request := &gatev2.SubmitJobResultsRequest{
				JobId:             jobId,
				EndOfTransmission: endOfTransmission,
			}
			if resultJSON != "" {
				var resultMap map[string]interface{}
				if err := json.Unmarshal([]byte(resultJSON), &resultMap); err != nil {
					return fmt.Errorf("invalid result JSON: %w", err)
				}
				// For demonstration, only support error_result (with message)
				if msg, ok := resultMap["error_result"].(map[string]interface{}); ok {
					if message, ok := msg["message"].(string); ok {
						request.Result = &gatev2.SubmitJobResultsRequest_ErrorResult_{
							ErrorResult: &gatev2.SubmitJobResultsRequest_ErrorResult{
								Message: message,
							},
						}
					}
				}
				// Add more result types as needed
			}

			// Call the client method
			resp, err := client.SubmitJobResults(ctx, request)
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
	cmd.Flags().StringVar(&jobId, "job-id", "", "Job ID (required)")
	cmd.Flags().BoolVar(&endOfTransmission, "end-of-transmission", false, "End of transmission (optional)")
	cmd.Flags().StringVar(&resultJSON, "result", "", "Result as JSON (optional, e.g. '{\"error_result\":{\"message\":\"fail\"}}')")
	cmd.MarkFlagRequired("job-id")
	return cmd
}
