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
		Short:        "List vault entries matching FILTER without password",
		Long:         "List vault entries matching FILTER without password",
		PreRunE:      listPreRunCmd,
		RunE:         listRunCmd,
		SilenceUsage: true,
	}
	logger *logrus.Logger
)

func init() {
	listCmd.Flags().BoolVarP(&sort, "sort", "s", false, "Sort the output by title and username.")
	listCmd.Flags().BoolVar(&trashed, "trashed", false, "Show trashed items.")
	listCmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Output the data as JSON.")
	listCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "Output the data as list, similar to SQLite line mode.")
	listCmd.Flags().BoolVarP(&yamlFlag, "yaml", "y", false, "Output the data as YAML.")
	listCmd.Flags().BoolVarP(&tableFlag, "table", "t", false, "Output the data as a table.")
	rootCmd.AddCommand(listCmd)
}

func listPreRunCmd(cmd *cobra.Command, args []string) error {
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
	if err := vault.Open(credentials); err != nil {
		logger.WithError(err).Error("could not open vault")
		logger.Exit(2)
	}
	logger.Debug("opened vault")

	cards, err := vault.GetEntries(cardType, cardCategory, filters)
	if err != nil {
		return err
	}

	if sort {
		util.SortEntries(cards)
	}

	output.GenerateOutput(logger, "list", cmd.Flags(), &cards)

	return nil
}
