package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tcncloud/sati-go/pkg/ports"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func AddScrubListEntriesCmd(configPath *string) *cobra.Command {
	var scrubListID, entriesJSON, countryCode string

	cmd := &cobra.Command{
		Use:   "add-scrub-list-entries",
		Short: "Call GateService.AddScrubListEntries",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := saticonfig.LoadConfig(*configPath)
			if err != nil {
				return err
			}

			client, err := saticlient.NewClient(cfg)
			if err != nil {
				return err
			}
			defer handleClientClose(client)

			ctx, cancel := createContext(DefaultTimeout)
			defer cancel()

			var entriesInput []struct {
				Content string `json:"content"`
				Notes   string `json:"notes,omitempty"`
			}
			if err := json.Unmarshal([]byte(entriesJSON), &entriesInput); err != nil {
				return fmt.Errorf("invalid entries JSON: %w", err)
			}
			// Build custom Params struct
			customEntries := make([]ports.ScrubListEntryInput, 0, len(entriesInput))
			for _, e := range entriesInput {
				entry := ports.ScrubListEntryInput{Content: e.Content}
				if e.Notes != "" {
					notesCopy := e.Notes // Create copy for pointer
					entry.Notes = &notesCopy
				}
				customEntries = append(customEntries, entry)
			}

			params := ports.AddScrubListEntriesParams{
				ScrubListID: scrubListID,
				Entries:     customEntries,
			}
			if countryCode != "" {
				params.CountryCode = &countryCode
			}

			// Call the client method with custom Params
			_, err = client.AddScrubListEntries(ctx, params)
			if err != nil {
				return err
			}
			// Response is now an empty struct on success
			fmt.Println("Successfully added scrub list entries.") // Provide feedback

			return nil
		},
	}
	cmd.Flags().StringVar(&scrubListID, "scrub-list-id", "", "Scrub List ID (required)")
	cmd.Flags().StringVar(&entriesJSON, "entries", "", "Entries as JSON array (required)")
	cmd.Flags().StringVar(&countryCode, "country-code", "", "Country code (optional)")
	markFlagRequired(cmd, "scrub-list-id")
	markFlagRequired(cmd, "entries")

	return cmd
}
