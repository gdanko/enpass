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
		Short:        "Copy the password of a vault entry to the clipboard",
		Long:         "Copy the password of a vault entry to the clipboard",
		PreRunE:      copyPreRunCmd,
		RunE:         copyRunCmd,
		SilenceUsage: true,
	}
)

func init() {
	GetCopyFlags(copyCmd)
	rootCmd.AddCommand(copyCmd)
}

func copyPreRunCmd(cmd *cobra.Command, args []string) error {
	logger = logrus.New()
	logLevel = logLevelMap[logLevelStr]
	logger.SetLevel(logLevel)
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
	if err := vault.Open(credentials, logLevel); err != nil {
		logger.WithError(err).Error("could not open vault")
		logger.Exit(2)
	}
	logger.Debug("opened vault")

	card, err := vault.GetEntry(cardType, recordCategory, recordTitle, recordLogin, recordUuid, caseSensitive, orderbyFlag, true)
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
