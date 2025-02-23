package cmd

import (
	"fmt"

	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/gdanko/enpass/util"
	"github.com/spf13/cobra"
)

var (
	passCmd = &cobra.Command{
		Use:          "pass",
		Short:        "Print the password of a vault entry to STDOUT",
		Long:         "Print the password of a vault entry to STDOUT",
		PreRun:       passPreRunCmd,
		Run:          passRunCmd,
		SilenceUsage: true,
	}
)

func init() {
	GetPassFlags(passCmd)
	rootCmd.AddCommand(passCmd)
}

func passPreRunCmd(cmd *cobra.Command, args []string) {
	logLevel = logLevelMap[logLevelStr]
	logger = util.ConfigureLogger(logLevel, flagNoColor)
}

func passRunCmd(cmd *cobra.Command, args []string) {
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

	card, err := vault.GetEntry(flagCardType, flagRecordCategory, flagRecordTitle, flagRecordLogin, flagRecordUuid, flagLabel, flagCaseSensitive, flagOrderBy, validOrderBy, true)
	if err != nil {
		logger.Error(err)
		logger.Exit(2)
	}
	fmt.Println(card.DecryptedValue)
}
