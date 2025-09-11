package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	gatev2 "github.com/tcncloud/sati-go/internal/genproto/tcnapi/exile/gate/v2"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

			// Build the request struct
			request := &gatev2.UpdateScrubListEntryRequest{
				ScrubListId: scrubListID,
				Content:     content,
			}
			if notes != "" {
				request.Notes = &wrapperspb.StringValue{Value: notes}
			}
			if countryCode != "" {
				request.CountryCode = &wrapperspb.StringValue{Value: countryCode}
			}
			if expiration != "" {
				t, err := time.Parse(time.RFC3339, expiration)
				if err != nil {
					return fmt.Errorf("invalid expiration format: %w", err)
				}
				request.Expiration = timestamppb.New(t)
			}

			// Call the client method
			resp, err := client.UpdateScrubListEntry(ctx, request)
			if err != nil {
				return err
			}
			fmt.Printf("%+v\n", resp)

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
