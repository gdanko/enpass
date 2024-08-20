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
	logger = util.ConfigureLogger(logLevel)
}

func passRunCmd(cmd *cobra.Command, args []string) {
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

	card, err := vault.GetEntry(cardType, recordCategory, recordTitle, recordLogin, recordUuid, caseSensitive, orderbyFlag, true)
	if err != nil {
		logger.Error(err)
		logger.Exit(2)
	}
	fmt.Println(card.DecryptedValue)
}
