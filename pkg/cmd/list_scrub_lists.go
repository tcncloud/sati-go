package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tcncloud/sati-go/pkg/ports"
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

			// Build the params struct
			params := ports.ListScrubListsParams{}

			// Call the client method
			resp, err := client.ListScrubLists(ctx, params)
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
				fmt.Printf("Scrub Lists:\n")
				if len(resp.ScrubLists) == 0 {
					fmt.Println("  No scrub lists found")
				} else {
					for _, scrubList := range resp.ScrubLists {
						fmt.Printf("  - ID: %s, Name: %s\n", scrubList.ID, scrubList.Name)
						if scrubList.Description != "" {
							fmt.Printf("    Description: %s\n", scrubList.Description)
						}
					}
				}
			}

			return nil
		},
	}
}
