package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func RemoveScrubListEntriesCmd(configPath *string) *cobra.Command {
	var scrubListID, entries string

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
			defer handleClientClose(client) // Ensure connection is closed

			ctx, cancel := createContext(DefaultTimeout)
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
			params := saticlient.RemoveScrubListEntriesParams{
				ScrubListID: scrubListID,
				EntryIDs:    entriesList,
			}

			// Call the client method
			resp, err := client.RemoveScrubListEntries(ctx, params)
			if err != nil {
				return err
			}
			fmt.Printf("%+v\n", resp)

			return nil
		},
	}
	cmd.Flags().StringVar(&scrubListID, "scrub-list-id", "", "Scrub List ID (required)")
	cmd.Flags().StringVar(&entries, "entries", "", "Comma-separated list of entries to remove (required)")
	markFlagRequired(cmd, "scrub-list-id")
	markFlagRequired(cmd, "entries")

	return cmd
}
