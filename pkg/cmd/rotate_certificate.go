package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	gatev2 "buf.build/gen/go/tcn/exileapi/protocolbuffers/go/tcnapi/exile/gate/v2"
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
			defer client.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			request := &gatev2.RotateCertificateRequest{
				CertificateHash: certificateHash,
			}
			resp, err := client.RotateCertificate(ctx, request)
			if err != nil {
				return err
			}
			if OutputFormat == "json" {
				data, err := json.MarshalIndent(resp, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
			} else {
				fmt.Printf("%+v\n", resp)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&certificateHash, "certificate-hash", "", "Certificate hash (required)")
	cmd.MarkFlagRequired("certificate-hash")
	return cmd
}
