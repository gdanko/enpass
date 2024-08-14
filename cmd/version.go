package cmd

import (
	"fmt"

	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:          "version",
		Short:        "Print the current enpass version",
		Long:         "Print the current enpass version",
		RunE:         runVersionCmd,
		SilenceUsage: true,
	}
)

func init() {
	// This is only here to override the required vault flag for other commands
	versionCmd.Flags().StringVarP(&vaultPath, "vault", "v", "", "Path to your Enpass vault")
	rootCmd.AddCommand(versionCmd)
}

func runVersionCmd(cmd *cobra.Command, args []string) error {
	fmt.Fprintf(cmd.OutOrStdout(), "%s\n", enpass.Version(true, true))

	return nil
}
