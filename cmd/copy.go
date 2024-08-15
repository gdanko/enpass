package cmd

import (
	"fmt"
	"strings"

	"github.com/gdanko/enpass/pkg/clipboard"
	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/gdanko/enpass/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	copyCmd = &cobra.Command{
		Use:          "copy",
		Short:        "Copy the password of a vault entry matching FILTER to the clipboard",
		Long:         "Copy the password of a vault entry matching FILTER to the clipboard",
		PreRunE:      copyPreRunCmd,
		RunE:         copyRunCmd,
		SilenceUsage: true,
	}
)

func init() {
	// copyCmd.Flags().BoolVar(&trashed, "trashed", false, "Show trashed items")
	copyCmd.Flags().BoolVar(&clipboardPrimary, "clipboardPrimary", false, "Use primary X selection instead of clipboard.")
	rootCmd.AddCommand(copyCmd)
}

func copyPreRunCmd(cmd *cobra.Command, args []string) error {
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

func copyRunCmd(cmd *cobra.Command, args []string) error {
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

	card, err := vault.GetEntry(cardType, cardCategory, filters, true)
	if err != nil {
		return err
	}

	if clipboardPrimary {
		clipboard.Primary = true
		logger.Debug("primary X selection enabled")
	}

	if err = clipboard.WriteAll(card.DecryptedValue); err != nil {
		logger.WithError(err).Fatal("could not copy password to clipboard")
	}

	fmt.Printf("The password for \"%s\" was copied to the clipboard\n", strings.TrimSpace(card.Title))

	return nil
}
