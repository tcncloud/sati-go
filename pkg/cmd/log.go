package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func LogCmd(configPath *string) *cobra.Command {
	var payload string

	cmd := &cobra.Command{
		Use:   "log",
		Short: "Call GateService.Log",
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

			params := saticlient.LogParams{
				Level:   "INFO", // Default level
				Message: payload,
			}
			resp, err := client.Log(ctx, params)
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
				fmt.Println("Log entry submitted successfully")
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&payload, "payload", "", "Log payload (required)")
	markFlagRequired(cmd, "payload")

	return cmd
}
