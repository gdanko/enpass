package cmd

import (
	"github.com/gdanko/enpass/pkg/output"
	"github.com/gdanko/enpass/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	showCmd = &cobra.Command{
		Use:          "show",
		Short:        "List vault entries matching FILTER with password",
		Long:         "List vault entries matching FILTER with password",
		PreRunE:      showPreRunCmd,
		RunE:         showRunCmd,
		SilenceUsage: true,
	}
)

func init() {
	showCmd.Flags().BoolVarP(&sort, "sort", "s", false, "Sort the output by title and username")
	showCmd.Flags().BoolVar(&trashed, "trashed", false, "Show trashed items")
	showCmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Output the data as JSON.")
	showCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "Output the data as list, similar to SQLite line mode.")
	showCmd.Flags().BoolVarP(&yamlFlag, "yaml", "y", false, "Output the data as YAML.")
	rootCmd.AddCommand(showCmd)
}

func showPreRunCmd(cmd *cobra.Command, args []string) error {
	logger = logrus.New()

	logLevel, err = logrus.ParseLevel(logLevelMap[logLevelStr])
	if err != nil {
		logrus.WithError(err).Fatal("invalid log level specified")
	}
	logger.SetLevel(logLevel)

	return nil
}

func showRunCmd(cmd *cobra.Command, args []string) error {
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

	cards, err := vault.GetEntries(cardType, []string{})
	if err != nil {
		panic(err)
	}

	if sort {
		util.SortEntries(cards)
	}

	output.GenerateOutput(logger, "show", cmd.Flags(), &cards)

	return nil
}
