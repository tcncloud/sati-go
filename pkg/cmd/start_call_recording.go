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

func StartCallRecordingCmd(configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "start-call-recording",
		Short: "Call GateService.StartCallRecording",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := sati.LoadConfig(*configPath)
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

			// TODO: Add flags for StartCallRecordingRequest fields
			// Build the request struct
			request := &gatev2.StartCallRecordingRequest{}

			// Call the client method
			resp, err := client.StartCallRecording(ctx, request)
			if err != nil {
				return err
			}
			fmt.Printf("%+v\n", resp)
			return nil
		},
	}
}
