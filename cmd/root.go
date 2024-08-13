package cmd

import (
	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cardType        string
	credentials     *enpass.VaultCredentials
	debug           bool
	defaultLogLevel string
	err             error
	keyFilePath     string
	jsonFlag        bool
	listFlag        bool
	logLevel        logrus.Level
	logLevelStr     = "4"
	logLevelMap     = map[string]string{
		"5": "debug",
		"4": "info",
		"3": "warn",
		"2": "error",
		"1": "panic",
	}
	nonInteractive bool
	pinEnable      bool
	rootCmd        = &cobra.Command{
		Use:   "enpass",
		Short: "enpass is a command line interface for the Enpass password manager",
		Long:  "enpass is a command line interface for the Enpass password manager",
	}
	sort      bool
	trashed   bool
	vault     *enpass.Vault
	vaultPath = "/home/gdanko/Documents/Enpass/Vaults/primary"
	yamlFlag  bool
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Run enpass in debug mode")
	rootCmd.PersistentFlags().StringVarP(&cardType, "type", "t", "password", "The type of your card. (password, ...)")
	rootCmd.PersistentFlags().StringVarP(&vaultPath, "vault", "v", "", "Path to your Enpass vault")
	rootCmd.PersistentFlags().StringVarP(&keyFilePath, "keyfile", "k", "", "Path to your Enpass vault keyfile")
	rootCmd.PersistentFlags().BoolVar(&nonInteractive, "nonInteractive", false, "Disable prompts and fail instead.")
	rootCmd.PersistentFlags().StringVar(&defaultLogLevel, "log", "4", "The log level from debug (5) to panic (1)")
	rootCmd.PersistentFlags().BoolVar(&pinEnable, "pin", false, "Enable PIN")
	rootCmd.MarkPersistentFlagRequired("vault")
}
