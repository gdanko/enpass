package cmd

import (
	"fmt"

	"github.com/gdanko/enpass/util"
	"github.com/spf13/cobra"
)

func GetListFlags(cmd *cobra.Command) {
	getListShowFlags(cmd)
}

func GetShowFlags(cmd *cobra.Command) {
	getListShowFlags(cmd)
}

func getListShowFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&trashedFlag, "trashed", false, "Show trashed items.")
	cmd.Flags().StringArrayVarP(&orderbyFlag, "orderby", "o", []string{"title"}, "Specify fields to sort by. Can be used multiple times.")
	cmd.Flags().BoolVar(&listFlag, "list", false, "Output the data as list, similar to SQLite line mode.")
	cmd.Flags().BoolVar(&yamlFlag, "yaml", false, "Output the data as YAML.")
	cmd.Flags().BoolVar(&tableFlag, "table", false, "Output the data as a table.")
}

func GetPersistenFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&vaultPathFlag, "vault", "v", "", "Path to your Enpass vault.")
	cmd.PersistentFlags().StringVar(&cardType, "type", "password", "The type of your card. (password, ...)")
	cmd.PersistentFlags().StringArrayVarP(&recordTitle, "title", "t", []string{}, "Filter based on record title. Wildcards (%) are allowed. Can be used multiple times.")
	cmd.PersistentFlags().StringArrayVarP(&recordCategory, "category", "c", []string{}, "Filter based on record category. Wildcards (%) are allowed. Can be used multiple times.")
	cmd.PersistentFlags().StringArrayVarP(&recordLogin, "login", "l", []string{}, "Filter based on record login. Wildcards (%) are allowed. Can be used multiple times.")
	cmd.PersistentFlags().StringArrayVarP(&recordFieldLabel, "label", "y", defaultLabels, "Filter based on record field label. Can be used multiple times")
	cmd.PersistentFlags().StringArrayVarP(&recordUuid, "uuid", "u", []string{}, "Filter based on record uuid. Can be used multiple times.")
	cmd.PersistentFlags().StringVarP(&keyFilePath, "keyfile", "k", "", "Path to your Enpass vault keyfile.")
	cmd.PersistentFlags().StringVar(&logLevelStr, "log", defaultLogLevel, fmt.Sprintf("The log level, one of: %s", util.ReturnLogLevels(logLevelMap)))
	cmd.PersistentFlags().BoolVarP(&nonInteractive, "non-interactive", "n", false, "Disable prompts and fail instead.")
	cmd.PersistentFlags().BoolVar(&caseSensitive, "sensitive", false, "Force category and title searches to be case-sensitive.")
	cmd.PersistentFlags().BoolVar(&nocolorFlag, "nocolor", false, "Disable colorized output and logging.")
	cmd.PersistentFlags().BoolVarP(&pinEnable, "pin", "p", false, "Enable PIN.")
}

func GetCopyFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&clipboardPrimary, "clipboardPrimary", false, "Use primary X selection instead of clipboard.")
	cmd.Flags().StringArrayVarP(&orderbyFlag, "orderby", "o", []string{"title"}, "Specify fields to sort by. Can be used multiple times.")
}

func GetPassFlags(cmd *cobra.Command) {
	cmd.Flags().StringArrayVarP(&orderbyFlag, "orderby", "o", []string{"title"}, "Specify fields to sort by. Can be used multiple times.")
}
