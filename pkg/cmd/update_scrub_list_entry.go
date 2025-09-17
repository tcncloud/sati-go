package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tcncloud/sati-go/pkg/ports"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func UpdateScrubListEntryCmd(configPath *string) *cobra.Command {
	var scrubListID, notes, content, expiration, countryCode string

	cmd := &cobra.Command{
		Use:   "update-scrub-list-entry",
		Short: "Call GateService.UpdateScrubListEntry",
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
			params := ports.UpdateScrubListEntryParams{
				ScrubListID: scrubListID,
				Content:     content,
			}
			if notes != "" {
				params.Notes = &notes
			}
			// Note: CountryCode and Expiration are not available in the current params structure

			// Call the client method
			resp, err := client.UpdateScrubListEntry(ctx, params)
			if err != nil {
				return err
			}
			if OutputFormat == OutputFormatJSON {
				data, err := json.MarshalIndent(resp, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
			} else {
				fmt.Println("Scrub list entry updated successfully")
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&scrubListID, "scrub-list-id", "", "Scrub List ID (required)")
	cmd.Flags().StringVar(&notes, "notes", "", "Notes (optional)")
	cmd.Flags().StringVar(&content, "content", "", "Content to block (required)")
	cmd.Flags().StringVar(&expiration, "expiration", "", "Expiration timestamp (RFC3339, optional)")
	cmd.Flags().StringVar(&countryCode, "country-code", "", "Country code (optional)")
	markFlagRequired(cmd, "scrub-list-id")
	markFlagRequired(cmd, "content")

	return cmd
}
