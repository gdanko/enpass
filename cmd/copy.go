package cmd

import (
	"fmt"
	"strings"

	"github.com/gdanko/enpass/pkg/clipboard"
	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/gdanko/enpass/util"
	"github.com/spf13/cobra"
)

var (
	copyCmd = &cobra.Command{
		Use:          "copy",
		Short:        "Copy the password of a vault entry to the clipboard",
		Long:         "Copy the password of a vault entry to the clipboard",
		PreRun:       copyPreRunCmd,
		Run:          copyRunCmd,
		SilenceUsage: true,
	}
)

func init() {
	GetCopyFlags(copyCmd)
	rootCmd.AddCommand(copyCmd)
}

func copyPreRunCmd(cmd *cobra.Command, args []string) {
	logLevel = logLevelMap[logLevelStr]
	logger = util.ConfigureLogger(logLevel, flagNoColor)
}

func copyRunCmd(cmd *cobra.Command, args []string) {
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

	card, err := vault.GetEntry(flagCardType, flagRecordCategory, flagRecordTitle, flagRecordLogin, flagRecordUuid, flagLabel, flagCaseSensitive, flagOrderBy, validOrderBy, true)
	if err != nil {
		logger.Error(err)
		logger.Exit(2)
	}

	if flagClipboardPrimary {
		clipboard.Primary = true
		logger.Debug("primary X selection enabled")
	}

	if err = clipboard.WriteAll(card.DecryptedValue); err != nil {
		logger.WithError(err).Fatal("could not copy password to clipboard")
	}

	fmt.Printf("The password for \"%s\" was copied to the clipboard\n", strings.TrimSpace(card.Title))
}
