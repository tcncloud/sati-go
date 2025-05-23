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

func LogCmd(configPath *string) *cobra.Command {
	var payload string
	cmd := &cobra.Command{
		Use:   "log",
		Short: "Call GateService.Log",
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

			request := &gatev2.LogRequest{
				Payload: payload,
			}
			resp, err := client.Log(ctx, request)
			if err != nil {
				return err
			}
			fmt.Printf("%+v\n", resp)
			return nil
		},
	}
	cmd.Flags().StringVar(&payload, "payload", "", "Log payload (required)")
	cmd.MarkFlagRequired("payload")
	return cmd
}
