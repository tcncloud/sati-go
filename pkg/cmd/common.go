// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright 2024 TCN Inc

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	saticlient "github.com/tcncloud/sati-go/pkg/sati/client"
	saticonfig "github.com/tcncloud/sati-go/pkg/sati/config"
)

// Common error messages.
var (
	ErrPartnerAgentIDRequired = errors.New("--partner-agent-id is required")
	ErrSkillIDRequired        = errors.New("--skill-id is required")
	ErrUserIDRequired         = errors.New("--user-id is required")
	ErrCallSidRequired        = errors.New("--call-sid is required")
	ErrRecordingSidRequired   = errors.New("--recording-sid is required")
	ErrPhoneNumberRequired    = errors.New("--phone-number is required")
	ErrNewStateRequired       = errors.New("--new-state is required")
	ErrUsernameRequired       = errors.New("--username is required")
	ErrFirstNameRequired      = errors.New("--first-name is required")
	ErrLastNameRequired       = errors.New("--last-name is required")
	ErrPasswordRequired       = errors.New("--password is required")
	ErrCallTypeRequired       = errors.New("--call-type is required")
	ErrKeyRequired            = errors.New("--key is required")
	ErrValueRequired          = errors.New("--value is required")
	ErrRequiredFieldsMissing  = errors.New("required fields are missing")
	ErrInvalidCallType        = errors.New("invalid call type")
	ErrInvalidNewState        = errors.New("invalid new state")
	ErrAtLeastOneDestination  = errors.New("at least one destination must be provided")
)

// Common constants.
const (
	OutputFormatJSON = "json"
	DefaultTimeout   = 10 * time.Second
	LongTimeout      = 30 * time.Second
)

// createClient creates a new client with proper error handling.
func createClient(configPath *string) (*saticlient.Client, error) {
	cfg, err := saticonfig.LoadConfig(*configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	client, err := saticlient.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client, nil
}

// createContext creates a context with timeout.
func createContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// handleClientClose handles client.Close() with error checking.
func handleClientClose(client *saticlient.Client) {
	if err := client.Close(); err != nil {
		// Log error but don't fail the command
		fmt.Fprintf(os.Stderr, "Warning: failed to close client: %v\n", err)
	}
}

// markFlagRequired marks a flag as required with error handling.
func markFlagRequired(cmd *cobra.Command, flagName string) {
	if err := cmd.MarkFlagRequired(flagName); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to mark flag %s as required: %v\n", flagName, err)
	}
}

// markPersistentFlagRequired marks a persistent flag as required with error handling.
func markPersistentFlagRequired(cmd *cobra.Command, flagName string) {
	if err := cmd.MarkPersistentFlagRequired(flagName); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to mark persistent flag %s as required: %v\n", flagName, err)
	}
}

// outputJSON outputs data in JSON format.
func outputJSON(data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(jsonData))

	return nil
}

// createSkillCommand creates a skill command with common setup.
func createSkillCommand(use, short string, configPath *string, partnerAgentID, skillID *string, operation func(*saticlient.Client, context.Context, saticlient.AssignAgentSkillParams) error, successMsg string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if *partnerAgentID == "" {
				return ErrPartnerAgentIDRequired
			}
			if *skillID == "" {
				return ErrSkillIDRequired
			}

			client, err := createClient(configPath)
			if err != nil {
				return err
			}
			defer handleClientClose(client)

			ctx, cancel := createContext(DefaultTimeout)
			defer cancel()

			params := saticlient.AssignAgentSkillParams{
				PartnerAgentID: *partnerAgentID,
				SkillID:        *skillID,
			}

			err = operation(client, ctx, params)
			if err != nil {
				return err
			}

			fmt.Println(successMsg)

			return nil
		},
	}

	cmd.Flags().StringVar(partnerAgentID, "partner-agent-id", "", "Partner Agent ID (required)")
	cmd.Flags().StringVar(skillID, "skill-id", "", "Skill ID (required)")
	markFlagRequired(cmd, "partner-agent-id")
	markFlagRequired(cmd, "skill-id")

	return cmd
}
