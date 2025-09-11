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
	"fmt"
	"time"

	"github.com/spf13/cobra"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func TransferCmd(configPath *string) *cobra.Command {
	var (
		callSid                 string
		receivingPartnerAgentID string
		outboundPhoneNumber     string
		outboundCallerID        string
		outboundPoolID          string
		outboundRecordID        string
		queueID                 string
	)

	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "Call GateService.Transfer",
		RunE: func(cmd *cobra.Command, args []string) error {
			if callSid == "" {
				return ErrCallSidRequired
			}

			// Validate that at least one destination is provided
			if receivingPartnerAgentID == "" && outboundPhoneNumber == "" && queueID == "" {
				return ErrAtLeastOneDestination
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
			defer client.Close() // Ensure connection is closed

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Build the custom Params struct
			params := saticlient.TransferParams{
				CallSid: callSid,
			}

			if receivingPartnerAgentID != "" {
				params.ReceivingPartnerAgentID = &receivingPartnerAgentID
			}

			if outboundPhoneNumber != "" {
				outbound := &saticlient.TransferOutbound{
					PhoneNumber: outboundPhoneNumber,
				}
				if outboundCallerID != "" {
					outbound.CallerID = &outboundCallerID
				}
				if outboundPoolID != "" {
					outbound.PoolID = &outboundPoolID
				}
				if outboundRecordID != "" {
					outbound.RecordID = &outboundRecordID
				}
				params.Outbound = outbound
			}

			if queueID != "" {
				params.Queue = &saticlient.TransferQueue{
					QueueID: queueID,
				}
			}

			// Call the client method with custom Params
			_, err = client.Transfer(ctx, params)
			if err != nil {
				return err
			}

			fmt.Println("Transfer initiated successfully")

			return nil
		},
	}
	cmd.Flags().StringVar(&callSid, "call-sid", "", "Call SID (required)")
	cmd.Flags().StringVar(&receivingPartnerAgentID, "receiving-partner-agent-id", "", "Receiving Partner Agent ID")
	cmd.Flags().StringVar(&outboundPhoneNumber, "outbound-phone-number", "", "Outbound phone number")
	cmd.Flags().StringVar(&outboundCallerID, "outbound-caller-id", "", "Outbound caller ID")
	cmd.Flags().StringVar(&outboundPoolID, "outbound-pool-id", "", "Outbound pool ID")
	cmd.Flags().StringVar(&outboundRecordID, "outbound-record-id", "", "Outbound record ID")
	cmd.Flags().StringVar(&queueID, "queue-id", "", "Queue ID")
	markFlagRequired(cmd, "call-sid")

	return cmd
}
