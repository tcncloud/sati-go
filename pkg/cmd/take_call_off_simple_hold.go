package cmd

import (
	"context"
	"fmt"
	"time"

	gatev2 "buf.build/gen/go/tcn/exileapi/protocolbuffers/go/tcnapi/exile/gate/v2"
	"github.com/spf13/cobra"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func TakeCallOffSimpleHoldCmd(configPath *string) *cobra.Command {
	var partnerAgentId string
	cmd := &cobra.Command{
		Use:   "take-call-off-simple-hold",
		Short: "Call GateService.TakeCallOffSimpleHold",
		RunE: func(cmd *cobra.Command, args []string) error {
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

			request := &gatev2.TakeCallOffSimpleHoldRequest{
				PartnerAgentId: partnerAgentId,
			}
			resp, err := client.TakeCallOffSimpleHold(ctx, request)
			if err != nil {
				return err
			}
			fmt.Printf("%+v\n", resp)
			return nil
		},
	}
	cmd.Flags().StringVar(&partnerAgentId, "partner-agent-id", "", "Partner Agent ID (required)")
	cmd.MarkFlagRequired("partner-agent-id")
	return cmd
}
