package cmd

import (
	"fmt"

	"github.com/gdanko/enpass/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	passCmd = &cobra.Command{
		Use:          "pass",
		Short:        "Print the password of a vault entry matching FILTER to stdout",
		Long:         "Print the password of a vault entry matching FILTER to stdout",
		PreRunE:      passPreRunCmd,
		RunE:         passRunCmd,
		SilenceUsage: true,
	}
)

func init() {
	rootCmd.AddCommand(passCmd)
}

func passPreRunCmd(cmd *cobra.Command, args []string) error {
	logger = logrus.New()

	logLevel, err = logrus.ParseLevel(logLevelMap[logLevelStr])
	if err != nil {
		logrus.WithError(err).Fatal("invalid log level specified")
	}
	logger.SetLevel(logLevel)

	if len(cmd.Flags().Args()) > 0 {
		filters = cmd.Flags().Args()[0:]
	} else {
		filters = []string{}
	}

	return nil
}

func passRunCmd(cmd *cobra.Command, args []string) error {
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

	card, err := vault.GetEntry(cardType, filters, true)
	if err != nil {
		return err
	}
	fmt.Println(card.DecryptedValue)

	return nil
}
