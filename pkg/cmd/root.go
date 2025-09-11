package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var OutputFormat string
var rootCmd = &cobra.Command{
	Use:   "sati-client",
	Short: "Sati Client - CLI for exile gateway that exposes the API",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	var configPath string
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to base64-encoded JSON config file")
	rootCmd.MarkPersistentFlagRequired("config")
	rootCmd.PersistentFlags().StringVarP(&OutputFormat, "output", "o", "text", "Output format: json or text")

	rootCmd.AddCommand(
		GetClientConfigCmd(&configPath),
		GetOrgInfoCmd(&configPath),
		RotateCertificateCmd(&configPath),
		PollEventsCmd(&configPath),
		StreamJobsCmd(&configPath),
		SubmitJobResultsCmd(&configPath),
		GetAgentStatusCmd(&configPath),
		UpdateAgentStatusCmd(&configPath),
		ListAgentsCmd(&configPath),
		UpsertAgentCmd(&configPath),
		GetAgentByIdCmd(&configPath),
		GetAgentByPartnerIdCmd(&configPath),
		AddAgentCallResponseCmd(&configPath),
		ListHuntGroupPauseCodesCmd(&configPath),
		PutCallOnSimpleHoldCmd(&configPath),
		TakeCallOffSimpleHoldCmd(&configPath),
		DialCmd(&configPath),
		ListNCLRulesetNamesCmd(&configPath),
		StartCallRecordingCmd(&configPath),
		StopCallRecordingCmd(&configPath),
		GetRecordingStatusCmd(&configPath),
		ListScrubListsCmd(&configPath),
		AddScrubListEntriesCmd(&configPath),
		UpdateScrubListEntryCmd(&configPath),
		RemoveScrubListEntriesCmd(&configPath),
		ListSkillsCmd(&configPath),
		ListAgentSkillsCmd(&configPath),
		AssignAgentSkillCmd(&configPath),
		UnassignAgentSkillCmd(&configPath),
		LogCmd(&configPath),
		SearchVoiceRecordingsCmd(&configPath),
		GetVoiceRecordingDownloadLinkCmd(&configPath),
		ListSearchableRecordingFieldsCmd(&configPath),
		TransferCmd(&configPath),
	)
}
