package cmd

import (
	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/gdanko/enpass/pkg/output"
	"github.com/gdanko/enpass/util"
	"github.com/spf13/cobra"
)

var (
	showCmd = &cobra.Command{
		Use:          "show",
		Short:        "List vault entries, displaying the password",
		Long:         "List vault entries, displaying the password",
		PreRun:       showPreRunCmd,
		Run:          showRunCmd,
		SilenceUsage: true,
	}
)

func init() {
	GetShowFlags(showCmd)
	rootCmd.AddCommand(showCmd)
}

func showPreRunCmd(cmd *cobra.Command, args []string) {
	logLevel = logLevelMap[logLevelStr]
	logger = util.ConfigureLogger(logLevel, flagNoColor)
}

func showRunCmd(cmd *cobra.Command, args []string) {
	vaultPath := enpass.DetermineVaultPath(logger, flagVaultPath)
	vault, credentials, err = enpass.OpenVault(logger, flagEnablePin, flagNonInteractive, vaultPath, flagKeyFilePath, logLevel, flagNoColor)
	if err != nil {
		logger.Error(err)
		logger.Exit(2)
	}

	defer func() {
		vault.Close()
	}()
	if err := vault.Open(credentials, logLevel, flagNoColor); err != nil {
		logger.Error(err)
		logger.Exit(2)
	}
	logger.Debug("opened vault")

	cards, err := vault.GetEntries(flagCardType, flagRecordCategory, flagRecordTitle, flagRecordLogin, flagRecordUuid, flagLabel, flagCaseSensitive, flagOrderBy, validOrderBy)
	if err != nil {
		logger.Error(err)
		logger.Exit(2)
	}

	output.GenerateOutput(logger, "show", flagList, flagTable, flagTrashed, flagYaml, flagNoColor, &cards)
}
