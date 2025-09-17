package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tcncloud/sati-go/pkg/ports"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func PollEventsCmd(configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "poll-events",
		Short: "Call GateService.PollEvents",
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
			params := ports.PollEventsParams{}

			// Call the client method
			resp, err := client.PollEvents(ctx, params)
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
				fmt.Printf("Events:\n")
				if len(resp.Events) == 0 {
					fmt.Println("  No events found")
				} else {
					for i, event := range resp.Events {
						fmt.Printf("  Event %d:\n", i+1)
						fmt.Printf("    Type: %s\n", event.Type)
						if event.AgentCall != nil {
							fmt.Printf("    Agent Call: %+v\n", event.AgentCall)
						}
						if event.Telephony != nil {
							fmt.Printf("    Telephony: %+v\n", event.Telephony)
						}
						if event.AgentResponse != nil {
							fmt.Printf("    Agent Response: %+v\n", event.AgentResponse)
						}
					}
				}
			}

			return nil
		},
	}
}
