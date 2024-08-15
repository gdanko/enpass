package cmd

import "github.com/spf13/cobra"

func GetListFlags(cmd *cobra.Command) {
	getListShowFlags(cmd)
}

func GetShowFlags(cmd *cobra.Command) {
	getListShowFlags(cmd)
}

func getListShowFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&trashedFlag, "trashed", false, "Show trashed items.")
	cmd.Flags().StringArrayVarP(&orderbyFlag, "orderby", "o", []string{"title"}, "Specify fields to sort by. Can be used multiple times.")
	cmd.Flags().BoolVar(&jsonFlag, "json", false, "Output the data as JSON.")
	cmd.Flags().BoolVar(&listFlag, "list", false, "Output the data as list, similar to SQLite line mode.")
	cmd.Flags().BoolVar(&yamlFlag, "yaml", false, "Output the data as YAML.")
	cmd.Flags().BoolVar(&tableFlag, "table", false, "Output the data as a table.")
}

func GetPersistenFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&vaultPath, "vault", "v", "", "Path to your Enpass vault.")
	cmd.PersistentFlags().StringVar(&cardType, "type", "password", "The type of your card. (password, ...)")
	cmd.PersistentFlags().StringArrayVarP(&cardTitle, "title", "t", []string{}, "Filter based on record title. Wildcards (%) are allowed. Can be used multiple times.")
	cmd.PersistentFlags().StringArrayVarP(&cardCategory, "category", "c", []string{}, "Filter based on record category. Wildcards (%) are allowed. Can be used multiple times.")
	cmd.PersistentFlags().StringArrayVarP(&cardLogin, "login", "l", []string{}, "Filter based on record login. Wildcards (%) are allowed. Can be used multiple times.")
	cmd.PersistentFlags().StringVarP(&keyFilePath, "keyfile", "k", "", "Path to your Enpass vault keyfile.")
	cmd.PersistentFlags().StringVar(&defaultLogLevel, "log", "4", "The log level from debug (5) to panic (1).")
	cmd.PersistentFlags().BoolVarP(&nonInteractive, "nonInteractive", "n", false, "Disable prompts and fail instead.")
	cmd.PersistentFlags().BoolVar(&caseSensitive, "sensitive", false, "Force category and title searches to be case-sensitive.")
	cmd.PersistentFlags().BoolVarP(&pinEnable, "pin", "p", false, "Enable PIN.")
}

func GetCopyFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&clipboardPrimary, "clipboardPrimary", false, "Use primary X selection instead of clipboard.")
	cmd.Flags().StringArrayVarP(&orderbyFlag, "orderby", "o", []string{"title"}, "Specify fields to sort by. Can be used multiple times.")
}

func GetPassFlags(cmd *cobra.Command) {
	cmd.Flags().StringArrayVarP(&orderbyFlag, "orderby", "o", []string{"title"}, "Specify fields to sort by. Can be used multiple times.")
}
