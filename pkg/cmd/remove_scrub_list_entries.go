package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	gatev2 "buf.build/gen/go/tcn/exileapi/protocolbuffers/go/tcnapi/exile/gate/v2"
	"github.com/spf13/cobra"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func RemoveScrubListEntriesCmd(configPath *string) *cobra.Command {
	var scrubListId, entries string
	cmd := &cobra.Command{
		Use:   "remove-scrub-list-entries",
		Short: "Call GateService.RemoveScrubListEntries",
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

			var entriesList []string
			if entries != "" {
				for _, e := range strings.Split(entries, ",") {
					trimmed := strings.TrimSpace(e)
					if trimmed != "" {
						entriesList = append(entriesList, trimmed)
					}
				}
			}
			request := &gatev2.RemoveScrubListEntriesRequest{
				ScrubListId: scrubListId,
				Entries:     entriesList,
			}

			// Call the client method
			resp, err := client.RemoveScrubListEntries(ctx, request)
			if err != nil {
				return err
			}
			fmt.Printf("%+v\n", resp)
			return nil
		},
	}
	cmd.Flags().StringVar(&scrubListId, "scrub-list-id", "", "Scrub List ID (required)")
	cmd.Flags().StringVar(&entries, "entries", "", "Comma-separated list of entries to remove (required)")
	cmd.MarkFlagRequired("scrub-list-id")
	cmd.MarkFlagRequired("entries")
	return cmd
}
