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

func GetVoiceRecordingDownloadLinkCmd(configPath *string) *cobra.Command {
	var recordingSid string

	cmd := &cobra.Command{
		Use:   "get-voice-recording-download-link",
		Short: "Call GateService.GetVoiceRecordingDownloadLink",
		RunE: func(cmd *cobra.Command, args []string) error {
			if recordingSid == "" {
				return ErrRecordingSidRequired
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

			// Build the custom Params struct
			params := ports.GetVoiceRecordingDownloadLinkParams{RecordingSid: recordingSid}

			// Call the client method with custom Params
			resp, err := client.GetVoiceRecordingDownloadLink(ctx, params)
			if err != nil {
				return err
			}

			// Use the custom Result struct
			if OutputFormat == OutputFormatJSON {
				data, err := json.MarshalIndent(resp, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
			} else {
				fmt.Printf("Download URL: %s\nExpires At: %s\n", resp.DownloadURL, resp.ExpiresAt)
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&recordingSid, "recording-sid", "", "Recording SID (required)")
	markFlagRequired(cmd, "recording-sid")

	return cmd
}
