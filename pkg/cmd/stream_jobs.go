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

	gatev2 "git.tcncloud.net/experiments/sati-client/internal/genproto/tcnapi/exile/gate/v2"
	"git.tcncloud.net/experiments/sati-client/pkg/sati"
	"github.com/spf13/cobra"
)

func StreamJobsCmd(configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "stream-jobs",
		Short: "Call GateService.StreamJobs (streaming)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := sati.LoadConfig(*configPath)
			if err != nil {
				return err
			}
			conn, err := sati.SetupClient(cfg)
			if err != nil {
				return err
			}
			defer conn.Close()
			client := gatev2.NewGateServiceClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			stream, err := client.StreamJobs(ctx, &gatev2.StreamJobsRequest{})
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
