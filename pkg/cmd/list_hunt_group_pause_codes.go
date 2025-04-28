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

func ListHuntGroupPauseCodesCmd(configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "list-hunt-group-pause-codes",
		Short: "Call GateService.ListHuntGroupPauseCodes",
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
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			resp, err := client.ListHuntGroupPauseCodes(ctx, &gatev2.ListHuntGroupPauseCodesRequest{})
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
}
