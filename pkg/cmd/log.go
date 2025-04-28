package cmd

import (
	"context"
	"fmt"
	"time"

	gatev2 "git.tcncloud.net/experiments/sati-client/internal/genproto/tcnapi/exile/gate/v2"
	"git.tcncloud.net/experiments/sati-client/pkg/sati"
	"github.com/spf13/cobra"
)

func LogCmd(configPath *string) *cobra.Command {
	var payload string
	cmd := &cobra.Command{
		Use:   "log",
		Short: "Call GateService.Log",
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
