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

func SearchVoiceRecordingsCmd(configPath *string) *cobra.Command {
	var (
		startDate    string
		endDate      string
		agentID      string
		callSid      string
		recordingSid string
		searchQuery  string
		pageSize     int32
		pageToken    string
		searchFields []string
	)
	cmd := &cobra.Command{
		Use:   "search-voice-recordings",
		Short: "Call GateService.SearchVoiceRecordings",
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

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			// Build the custom Params struct
			params := saticlient.SearchVoiceRecordingsParams{
				SearchFields: searchFields,
			}

			if startDate != "" {
				params.StartDate = &startDate
			}
			if endDate != "" {
				params.EndDate = &endDate
			}
			if agentID != "" {
				params.AgentID = &agentID
			}
			if callSid != "" {
				params.CallSid = &callSid
			}
			if recordingSid != "" {
				params.RecordingSid = &recordingSid
			}
			if searchQuery != "" {
				params.SearchQuery = &searchQuery
			}
			if pageSize > 0 {
				params.PageSize = &pageSize
			}
			if pageToken != "" {
				params.PageToken = &pageToken
			}

			// Call the client stream method - returns a channel
			resultsChan := client.SearchVoiceRecordings(ctx, params)

			var recordings []*saticlient.VoiceRecording
			for result := range resultsChan {
				if result.Error != nil {
					return fmt.Errorf("error searching recordings: %w", result.Error)
				}
				if result.Recording != nil {
					recordings = append(recordings, result.Recording)
				}
			}

			// Process collected recordings
			if OutputFormat == "json" {
				data, err := json.MarshalIndent(recordings, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
			} else {
				for _, recording := range recordings {
					fmt.Printf("RecordingSid: %s, CallSid: %s, AgentID: %s, StartTime: %s, EndTime: %s, Duration: %d, FileSize: %d, Status: %s\n",
						recording.RecordingSid, recording.CallSid, recording.AgentID, recording.StartTime, recording.EndTime, recording.Duration, recording.FileSize, recording.Status)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&startDate, "start-date", "", "Start date for search (ISO format)")
	cmd.Flags().StringVar(&endDate, "end-date", "", "End date for search (ISO format)")
	cmd.Flags().StringVar(&agentID, "agent-id", "", "Agent ID to filter by")
	cmd.Flags().StringVar(&callSid, "call-sid", "", "Call SID to filter by")
	cmd.Flags().StringVar(&recordingSid, "recording-sid", "", "Recording SID to filter by")
	cmd.Flags().StringVar(&searchQuery, "search-query", "", "Search query text")
	cmd.Flags().Int32Var(&pageSize, "page-size", 0, "Number of results per page")
	cmd.Flags().StringVar(&pageToken, "page-token", "", "Token for pagination")
	cmd.Flags().StringSliceVar(&searchFields, "search-fields", []string{}, "Fields to search in")
	return cmd
}
