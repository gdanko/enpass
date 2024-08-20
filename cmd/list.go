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
	logger = util.ConfigureLogger(logLevel)
}

func listRunCmd(cmd *cobra.Command, args []string) {
	if vaultPath == "" {
		vaultPath, err = enpass.FindDefaultVaultPath()
		if err != nil {
			logger.Error(err)
			logger.Exit(2)
		}
	}

	err = enpass.ValidateVaultPath(vaultPath)
	if err != nil {
		logger.Error(err)
		logger.Exit(2)
	}

	vault, credentials, err = enpass.OpenVault(logger, pinEnable, nonInteractive, vaultPath, keyFilePath, logLevel)
	if err != nil {
		logger.Error(err)
		logger.Exit(2)
	}

	defer func() {
		vault.Close()
	}()
	if err := vault.Open(credentials, logLevel); err != nil {
		logger.WithError(err).Error("could not open vault")
		logger.Exit(2)
	}
	logger.Debug("opened vault")

	cards, err := vault.GetEntries(cardType, recordCategory, recordTitle, recordLogin, recordUuid, caseSensitive, orderbyFlag)
	if err != nil {
		logger.Error(err)
		logger.Exit(2)
	}

	output.GenerateOutput(logger, "list", jsonFlag, listFlag, tableFlag, trashedFlag, yamlFlag, &cards)
}
