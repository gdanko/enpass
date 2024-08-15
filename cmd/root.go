package cmd

import (
	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cardCategory     []string
	cardTitle        []string
	cardType         string
	caseSensitive    bool
	clipboardPrimary bool
	credentials      *enpass.VaultCredentials
	defaultLogLevel  string
	err              error
	keyFilePath      string
	jsonFlag         bool
	listFlag         bool
	logLevel         logrus.Level
	logLevelStr      = "4"
	logLevelMap      = map[string]string{
		"5": "debug",
		"4": "info",
		"3": "warn",
		"2": "error",
		"1": "panic",
	}
	nonInteractive bool
	orderbyFlag    []string
	pinEnable      bool
	rootCmd        = &cobra.Command{
		Use:   "enpass",
		Short: "enpass is a command line interface for the Enpass password manager",
		Long:  "enpass is a command line interface for the Enpass password manager",
	}
	tableFlag   bool
	trashedFlag bool
	validFields = []string{
		"category",
		"login",
		"title",
	}
	vault       *enpass.Vault
	vaultPath   string
	versionFull bool
	yamlFlag    bool
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	GetPersistenFlags(rootCmd)
}
