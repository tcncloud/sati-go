package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

func RotateCertificateCmd(configPath *string) *cobra.Command {
	var certificateHash string

	cmd := &cobra.Command{
		Use:   "rotate-certificate",
		Short: "Call GateService.RotateCertificate",
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

			params := saticlient.RotateCertificateParams{}
			resp, err := client.RotateCertificate(ctx, params)
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
				fmt.Println("Certificate rotated successfully")
				if resp.Certificate != "" {
					fmt.Printf("New certificate: %s\n", resp.Certificate)
				}
				if resp.PrivateKey != "" {
					fmt.Printf("New private key: %s\n", resp.PrivateKey)
				}
				if resp.CACertificate != "" {
					fmt.Printf("New CA certificate: %s\n", resp.CACertificate)
				}
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&certificateHash, "certificate-hash", "", "Certificate hash (required)")
	markFlagRequired(cmd, "certificate-hash")

	return cmd
}
