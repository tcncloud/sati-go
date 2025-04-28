package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	gatev2 "github.com/tcncloud/sati-go/internal/genproto/tcnapi/exile/gate/v2"
	"github.com/tcncloud/sati-go/pkg/sati"
)

func RemoveScrubListEntriesCmd(configPath *string) *cobra.Command {
	var scrubListId, entries string
	cmd := &cobra.Command{
		Use:   "remove-scrub-list-entries",
		Short: "Call GateService.RemoveScrubListEntries",
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
