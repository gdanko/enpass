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
	GetListFlags(listCmd)
	rootCmd.AddCommand(listCmd)
}

func listPreRunCmd(cmd *cobra.Command, args []string) {
	logLevel = logLevelMap[logLevelStr]
	logger = util.ConfigureLogger(logLevel, nocolorFlag)
}

func listRunCmd(cmd *cobra.Command, args []string) {
	vaultPath := enpass.DetermineVaultPath(logger, vaultPathFlag)
	vault, credentials, err = enpass.OpenVault(logger, pinEnable, nonInteractive, vaultPath, keyFilePath, logLevel, nocolorFlag)
	if err != nil {
		logger.Error(err)
		logger.Exit(2)
	}

	defer func() {
		vault.Close()
	}()
	if err := vault.Open(credentials, logLevel, nocolorFlag); err != nil {
		logger.Error(err)
		logger.Exit(2)
	}
	logger.Debug("opened vault")

	cards, err := vault.GetEntries(cardType, recordCategory, recordTitle, recordLogin, recordUuid, recordFieldLabel, caseSensitive, orderbyFlag)
	if err != nil {
		logger.Error(err)
		logger.Exit(2)
	}

	output.GenerateOutput(logger, "list", jsonFlag, listFlag, tableFlag, trashedFlag, yamlFlag, nocolorFlag, &cards)
}
