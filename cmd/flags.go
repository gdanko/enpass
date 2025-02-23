package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gdanko/enpass/util"
	"github.com/spf13/cobra"
)

func GetflagLists(cmd *cobra.Command) {
	getListShowFlags(cmd)
}

func GetShowFlags(cmd *cobra.Command) {
	getListShowFlags(cmd)
}

func getListShowFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&flagTrashed, "trashed", false, "Show trashed items.")
	cmd.Flags().StringArrayVarP(&flagOrderBy, "orderby", "o", []string{}, fmt.Sprintf("Specify fields to sort by. Can be used multiple times. Valid: %s", strings.Join(sort.StringSlice(validOrderBy), ", ")))
	cmd.Flags().BoolVar(&flagList, "list", false, "Output the data as list, similar to SQLite line mode.")
	cmd.Flags().BoolVar(&flagYaml, "yaml", false, "Output the data as YAML.")
	cmd.Flags().BoolVar(&flagTable, "table", false, "Output the data as a table.")
}

func GetPersistenFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&flagVaultPath, "vault", "v", "", "Path to your Enpass vault.")
	cmd.PersistentFlags().StringVar(&flagCardType, "type", "password", "The type of your card. (password, ...)")
	cmd.PersistentFlags().StringArrayVarP(&flagRecordTitle, "title", "t", []string{}, "Filter based on record title. Wildcards (%) are allowed. Can be used multiple times.")
	cmd.PersistentFlags().StringArrayVarP(&flagRecordCategory, "category", "c", []string{}, "Filter based on record category. Wildcards (%) are allowed. Can be used multiple times.")
	cmd.PersistentFlags().StringArrayVarP(&flagRecordLogin, "login", "l", []string{}, "Filter based on record login. Wildcards (%) are allowed. Can be used multiple times.")
	cmd.PersistentFlags().StringArrayVarP(&flagLabel, "label", "y", []string{}, "Filter based on record field label. Can be used multiple times")
	cmd.PersistentFlags().StringArrayVarP(&flagRecordUuid, "uuid", "u", []string{}, "Filter based on record uuid. Can be used multiple times.")
	cmd.PersistentFlags().StringVarP(&flagKeyFilePath, "keyfile", "k", "", "Path to your Enpass vault keyfile.")
	cmd.PersistentFlags().StringVar(&logLevelStr, "log", defaultLogLevel, fmt.Sprintf("The log level, one of: %s", util.ReturnLogLevels(logLevelMap)))
	cmd.PersistentFlags().BoolVarP(&flagNonInteractive, "non-interactive", "n", false, "Disable prompts and fail instead.")
	cmd.PersistentFlags().BoolVar(&flagCaseSensitive, "sensitive", false, "Force category and title searches to be case-sensitive.")
	cmd.PersistentFlags().BoolVar(&flagNoColor, "nocolor", false, "Disable colorized output and logging.")
	cmd.PersistentFlags().BoolVarP(&flagEnablePin, "pin", "p", false, "Enable PIN.")
}

func GetCopyFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&flagClipboardPrimary, "flagClipboardPrimary", false, "Use primary X selection instead of clipboard.")
	cmd.Flags().StringArrayVarP(&flagOrderBy, "orderby", "o", []string{"title"}, "Specify fields to sort by. Can be used multiple times.")
}

func GetPassFlags(cmd *cobra.Command) {
	cmd.Flags().StringArrayVarP(&flagOrderBy, "orderby", "o", []string{"title"}, "Specify fields to sort by. Can be used multiple times.")
}
