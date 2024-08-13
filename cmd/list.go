package cmd

import (
	"fmt"

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
	listCmd.Flags().BoolVarP(&sort, "sort", "s", false, "Sort the output by title and username")
	listCmd.Flags().BoolVar(&trashed, "trashed", false, "Show trashed items")
	listCmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Output the data as JSON.")
	listCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "Output the data as list, similar to SQLite line mode.")
	listCmd.Flags().BoolVarP(&yamlFlag, "yaml", "y", false, "Output the data as YAML.")
	rootCmd.AddCommand(listCmd)
}

func listPreRunCmd(cmd *cobra.Command, args []string) error {
	logger = logrus.New()

	logLevel, err = logrus.ParseLevel(logLevelMap[logLevelStr])
	if err != nil {
		logrus.WithError(err).Fatal("invalid log level specified")
	}
	fmt.Println(logLevel)
	logger.SetLevel(logLevel)

	return nil
}

func listRunCmd(cmd *cobra.Command, args []string) error {
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

	for _, card := range cards {
		if card.IsTrashed() && !trashed {
			continue
		}
		logger.Printf(
			"> title: %s"+
				"  login: %s"+
				"  cat.: %s",
			card.Title,
			card.Subtitle,
			card.Category,
		)
	}

	return nil
}
