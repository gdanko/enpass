package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	passCmd = &cobra.Command{
		Use:          "pass",
		Short:        "Print the password of a vault entry matching FILTER to stdout",
		Long:         "Print the password of a vault entry matching FILTER to stdout",
		PreRunE:      passPreRunCmd,
		RunE:         passRunCmd,
		SilenceUsage: true,
	}
)

func init() {
	rootCmd.AddCommand(passCmd)
}

func passPreRunCmd(cmd *cobra.Command, args []string) error {
	fmt.Println("passPreRun")
	return nil
}

func passRunCmd(cmd *cobra.Command, args []string) error {
	fmt.Println("passRun")
	return nil
}
