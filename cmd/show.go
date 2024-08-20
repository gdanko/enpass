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
	logger = util.ConfigureLogger(logLevel, nocolorFlag)
}

func showRunCmd(cmd *cobra.Command, args []string) {
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

	vault, credentials, err = enpass.OpenVault(logger, pinEnable, nonInteractive, vaultPath, keyFilePath, logLevel, nocolorFlag)
	if err != nil {
		logger.Error(err)
		logger.Exit(2)
	}

	defer func() {
		vault.Close()
	}()
	if err := vault.Open(credentials, logLevel, nocolorFlag); err != nil {
		logger.WithError(err).Error("could not open vault")
		logger.Exit(2)
	}
	logger.Debug("opened vault")

	cards, err := vault.GetEntries(cardType, recordCategory, recordTitle, recordLogin, recordUuid, caseSensitive, orderbyFlag)
	if err != nil {
		logger.Error(err)
		logger.Exit(2)
	}

	output.GenerateOutput(logger, "show", jsonFlag, listFlag, tableFlag, trashedFlag, yamlFlag, nocolorFlag, &cards)
}
