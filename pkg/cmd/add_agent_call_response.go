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
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tcncloud/sati-go/pkg/ports"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func AddAgentCallResponseCmd(configPath *string) *cobra.Command {
	var (
		partnerAgentID, callSid, callTypeStr, key, value string
		currentSessionID                                 int64
	)

	cmd := &cobra.Command{
		Use:   "add-agent-call-response",
		Short: "Call GateService.AddAgentCallResponse",
		RunE: func(cmd *cobra.Command, args []string) error {
			if partnerAgentID == "" || callSid == "" || callTypeStr == "" || key == "" || value == "" {
				return ErrRequiredFieldsMissing
			}
			// Validate call type (not used in the new implementation)
			if callTypeStr == "" {
				return fmt.Errorf("%w: %s", ErrInvalidCallType, callTypeStr)
			}
			cfg, err := saticonfig.LoadConfig(*configPath)
			if err != nil {
				return err
			}
			client, err := saticlient.NewClient(cfg)
			if err != nil {
				return err
			}
			defer handleClientClose(client)

			ctx, cancel := createContext(DefaultTimeout)
			defer cancel()

			// Convert callSid from string to int64
			callSidInt, err := strconv.ParseInt(callSid, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid call SID: %w", err)
			}

			params := ports.AddAgentCallResponseParams{
				PartnerAgentID: partnerAgentID,
				CallSid:        callSidInt,
				ResponseKey:    key,
				ResponseValue:  value,
				AgentSid:       currentSessionID,
			}

			resp, err := client.AddAgentCallResponse(ctx, params)
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
				fmt.Println("Agent call response added successfully")
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&partnerAgentID, "partner-agent-id", "", "Partner Agent ID (required)")
	cmd.Flags().StringVar(&callSid, "call-sid", "", "Call SID (required)")
	cmd.Flags().StringVar(&callTypeStr, "call-type", "", "Call Type (required, e.g. CALL_TYPE_INBOUND)")
	cmd.Flags().Int64Var(&currentSessionID, "current-session-id", 0, "Current Session ID (optional)")
	cmd.Flags().StringVar(&key, "key", "", "Key (required)")
	cmd.Flags().StringVar(&value, "value", "", "Value (required)")
	markFlagRequired(cmd, "partner-agent-id")
	markFlagRequired(cmd, "call-sid")
	markFlagRequired(cmd, "call-type")
	markFlagRequired(cmd, "key")
	markFlagRequired(cmd, "value")

	return cmd
}
