package cmd

import (
	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/gdanko/enpass/pkg/output"
	"github.com/gdanko/enpass/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:          "list",
		Short:        "List vault entries without displaying the password",
		Long:         "List vault entries without displaying the password",
		PreRun:       listPreRunCmd,
		Run:          listRunCmd,
		SilenceUsage: true,
	}
	logger *logrus.Logger
)

func init() {
	GetflagLists(listCmd)
	rootCmd.AddCommand(listCmd)
}

func listPreRunCmd(cmd *cobra.Command, args []string) {
	logLevel = logLevelMap[logLevelStr]
	logger = util.ConfigureLogger(logLevel, flagNoColor)
}

func listRunCmd(cmd *cobra.Command, args []string) {
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

	output.GenerateOutput(logger, "list", flagList, flagTable, flagTrashed, flagYaml, flagNoColor, &cards)
}
