package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	gatev2 "github.com/tcncloud/sati-go/internal/genproto/tcnapi/exile/gate/v2"
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
			defer handleClientClose(client) // Ensure connection is closed

			ctx, cancel := createContext(DefaultTimeout)
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
