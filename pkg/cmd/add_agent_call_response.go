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
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func AddAgentCallResponseCmd(configPath *string) *cobra.Command {
	var partnerAgentId, callSid, callTypeStr, key, value string
	var currentSessionId int64
	cmd := &cobra.Command{
		Use:   "add-agent-call-response",
		Short: "Call GateService.AddAgentCallResponse",
		RunE: func(cmd *cobra.Command, args []string) error {
			if partnerAgentId == "" || callSid == "" || callTypeStr == "" || key == "" || value == "" {
				return fmt.Errorf("--partner-agent-id, --call-sid, --call-type, --key, and --value are required")
			}
			callTypeEnum, ok := gatev2.CallType_value[callTypeStr]
			if !ok {
				return fmt.Errorf("invalid --call-type: %s", callTypeStr)
			}
			cfg, err := saticonfig.LoadConfig(*configPath)
			if err != nil {
				return err
			}
			client, err := saticlient.NewClient(cfg)
			if err != nil {
				return err
			}
			defer client.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			request := &gatev2.AddAgentCallResponseRequest{
				PartnerAgentId:   partnerAgentId,
				CallSid:          callSid,
				CallType:         gatev2.CallType(callTypeEnum),
				CurrentSessionId: currentSessionId,
				Key:              key,
				Value:            value,
			}

			resp, err := client.AddAgentCallResponse(ctx, request)
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
	cmd.Flags().StringVar(&partnerAgentId, "partner-agent-id", "", "Partner Agent ID (required)")
	cmd.Flags().StringVar(&callSid, "call-sid", "", "Call SID (required)")
	cmd.Flags().StringVar(&callTypeStr, "call-type", "", "Call Type (required, e.g. CALL_TYPE_INBOUND)")
	cmd.Flags().Int64Var(&currentSessionId, "current-session-id", 0, "Current Session ID (optional)")
	cmd.Flags().StringVar(&key, "key", "", "Key (required)")
	cmd.Flags().StringVar(&value, "value", "", "Value (required)")
	cmd.MarkFlagRequired("partner-agent-id")
	cmd.MarkFlagRequired("call-sid")
	cmd.MarkFlagRequired("call-type")
	cmd.MarkFlagRequired("key")
	cmd.MarkFlagRequired("value")
	return cmd
}
