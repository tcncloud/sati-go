package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func AddScrubListEntriesCmd(configPath *string) *cobra.Command {
	var scrubListId, entriesJSON, countryCode string
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
			defer client.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			var entriesInput []struct {
				Content string `json:"content"`
				Notes   string `json:"notes,omitempty"`
			}
			if err := json.Unmarshal([]byte(entriesJSON), &entriesInput); err != nil {
				return fmt.Errorf("invalid entries JSON: %w", err)
			}
			// Build custom Params struct
			customEntries := make([]saticlient.ScrubListEntryInput, 0, len(entriesInput))
			for _, e := range entriesInput {
				entry := saticlient.ScrubListEntryInput{Content: e.Content}
				if e.Notes != "" {
					notesCopy := e.Notes // Create copy for pointer
					entry.Notes = &notesCopy
				}
				customEntries = append(customEntries, entry)
			}

			params := saticlient.AddScrubListEntriesParams{
				ScrubListID: scrubListId,
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
	cmd.Flags().StringVar(&scrubListId, "scrub-list-id", "", "Scrub List ID (required)")
	cmd.Flags().StringVar(&entriesJSON, "entries", "", "Entries as JSON array (required)")
	cmd.Flags().StringVar(&countryCode, "country-code", "", "Country code (optional)")
	cmd.MarkFlagRequired("scrub-list-id")
	cmd.MarkFlagRequired("entries")
	return cmd
}
