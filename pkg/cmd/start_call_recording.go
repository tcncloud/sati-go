package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func StartCallRecordingCmd(configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "start-call-recording",
		Short: "Call GateService.StartCallRecording",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := saticonfig.LoadConfig(*configPath)
			if err != nil {
				return err
			}

			// Use the new client constructor
			client, err := saticlient.NewClient(cfg)
			if err != nil {
				return err
			}
			defer handleClientClose(client) // Ensure connection is closed

			ctx, cancel := createContext(DefaultTimeout)
			defer cancel()

			// Build the params struct
			params := saticlient.StartCallRecordingParams{}

			// Call the client method
			resp, err := client.StartCallRecording(ctx, params)
			if err != nil {
				return err
			}
			fmt.Printf("%+v\n", resp)

			return nil
		},
	}
}
