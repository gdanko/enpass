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
		PreRunE:      listPreRunCmd,
		RunE:         listRunCmd,
		SilenceUsage: true,
	}
	logger *logrus.Logger
)

func init() {
	GetListFlags(listCmd)
	rootCmd.AddCommand(listCmd)
}

func listPreRunCmd(cmd *cobra.Command, args []string) error {
	logger = logrus.New()
	logLevel = logLevelMap[logLevelStr]
	logger.SetLevel(logLevel)
	return nil
}

func listRunCmd(cmd *cobra.Command, args []string) error {
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
	if err := vault.Open(credentials, logLevel); err != nil {
		logger.WithError(err).Error("could not open vault")
		logger.Exit(2)
	}
	logger.Debug("opened vault")

	cards, err := vault.GetEntries(cardType, recordCategory, recordTitle, recordLogin, recordUuid, caseSensitive, orderbyFlag)
	if err != nil {
		return err
	}

	output.GenerateOutput(logger, "list", jsonFlag, listFlag, tableFlag, trashedFlag, yamlFlag, &cards)

	return nil
}
