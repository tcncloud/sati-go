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

func ListScrubListsCmd(configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "list-scrub-lists",
		Short: "Call GateService.ListScrubLists",
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
			defer client.Close() // Ensure connection is closed

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Build the request struct
			request := &gatev2.ListScrubListsRequest{}

			// Call the client method
			resp, err := client.ListScrubLists(ctx, request)
			if err != nil {
				return err
			}
			fmt.Printf("%+v\n", resp)
			return nil
		},
	}
}
