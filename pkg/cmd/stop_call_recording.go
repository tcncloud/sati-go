package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	gatev2 "github.com/tcncloud/sati-go/internal/genproto/tcnapi/exile/gate/v2"
	"github.com/tcncloud/sati-go/pkg/sati"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
)

func StopCallRecordingCmd(configPath *string) *cobra.Command {
	var partnerAgentId string
	cmd := &cobra.Command{
		Use:   "stop-call-recording",
		Short: "Call GateService.StopCallRecording",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := sati.LoadConfig(*configPath)
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

			request := &gatev2.StopCallRecordingRequest{
				PartnerAgentId: partnerAgentId,
			}
			resp, err := client.StopCallRecording(ctx, request)
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
