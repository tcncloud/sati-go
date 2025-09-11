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
	gatev2 "github.com/tcncloud/sati-go/internal/genproto/tcnapi/exile/gate/v2"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func StreamJobsCmd(configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "stream-jobs",
		Short: "Call GateService.StreamJobs (streaming)",
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

			ctx, cancel := createContext(LongTimeout)
			defer cancel()

			// Build the request struct
			request := &gatev2.StreamJobsRequest{}

			// Call the client stream method
			stream, err := client.StreamJobs(ctx, request)
			if err != nil {
				return err
			}
			for {
				msg, err := stream.Recv()
				if err != nil {
					break
				}
				if OutputFormat == "json" {
					data, err := json.MarshalIndent(msg, "", "  ")
					if err != nil {
						return err
					}
					fmt.Println(string(data))
				} else {
					fmt.Printf("%+v\n", msg)
				}
			}

			return nil
		},
	}
}
