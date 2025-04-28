package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	gatev2 "github.com/tcncloud/sati-go/internal/genproto/tcnapi/exile/gate/v2"
	"github.com/tcncloud/sati-go/pkg/sati"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func AddScrubListEntriesCmd(configPath *string) *cobra.Command {
	var scrubListId, entriesJSON, countryCode string
	cmd := &cobra.Command{
		Use:   "add-scrub-list-entries",
		Short: "Call GateService.AddScrubListEntries",
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
			var entriesInput []struct {
				Content string `json:"content"`
				Notes   string `json:"notes,omitempty"`
			}
			if err := json.Unmarshal([]byte(entriesJSON), &entriesInput); err != nil {
				return fmt.Errorf("invalid entries JSON: %w", err)
			}
			var entries []*gatev2.AddScrubListEntriesRequest_Entry
			for _, e := range entriesInput {
				entry := &gatev2.AddScrubListEntriesRequest_Entry{
					Content: e.Content,
				}
				if e.Notes != "" {
					entry.Notes = &wrapperspb.StringValue{Value: e.Notes}
				}
				entries = append(entries, entry)
			}
			request := &gatev2.AddScrubListEntriesRequest{
				ScrubListId: scrubListId,
				Entries:     entries,
				CountryCode: countryCode,
			}
			resp, err := client.AddScrubListEntries(ctx, request)
			if err != nil {
				return err
			}
			fmt.Printf("%+v\n", resp)
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
