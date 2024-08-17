package cmd

import (
	"fmt"

	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/gdanko/enpass/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	passCmd = &cobra.Command{
		Use:          "pass",
		Short:        "Print the password of a vault entry to STDOUT",
		Long:         "Print the password of a vault entry to STDOUT",
		PreRunE:      passPreRunCmd,
		RunE:         passRunCmd,
		SilenceUsage: true,
	}
)

func init() {
	GetPassFlags(passCmd)
	rootCmd.AddCommand(passCmd)
}

func passPreRunCmd(cmd *cobra.Command, args []string) error {
	logger = logrus.New()
	logLevel = logLevelMap[logLevelStr]
	logger.SetLevel(logLevel)
	return nil
}

func passRunCmd(cmd *cobra.Command, args []string) error {
	if vaultPath == "" {
		vaultPath, err = enpass.FindDefaultVaultPath()
		if err != nil {
			return err
		}
	}

	err = enpass.ValidateVaultPath(vaultPath)
	if err != nil {
		return err
	}

	vault, credentials, err = util.OpenVault(logger, pinEnable, nonInteractive, vaultPath, keyFilePath, logLevel)
	if err != nil {
		return err
	}

	defer func() {
		vault.Close()
	}()
	if err := vault.Open(credentials); err != nil {
		logger.WithError(err).Error("could not open vault")
		logger.Exit(2)
	}
	logger.Debug("opened vault")

	card, err := vault.GetEntry(cardType, recordCategory, recordTitle, recordLogin, recordUuid, caseSensitive, orderbyFlag, true)
	if err != nil {
		return err
	}
	fmt.Println(card.DecryptedValue)

	return nil
}
