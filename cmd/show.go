package cmd

import (
	"fmt"

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
	fmt.Println("showPreRun")
	return nil
}

func showRunCmd(cmd *cobra.Command, args []string) error {
	fmt.Println("showRun")
	return nil
}
