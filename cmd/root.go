package cmd

import (
	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cardCategory     []string
	cardLogin        []string
	cardTitle        []string
	cardType         string
	caseSensitive    bool
	clipboardPrimary bool
	credentials      *enpass.VaultCredentials
	defaultLogLevel  = "info"
	err              error
	keyFilePath      string
	jsonFlag         bool
	listFlag         bool
	logLevel         logrus.Level
	logLevelStr      string
	logLevelMap      = map[string]logrus.Level{
		"panic": logrus.PanicLevel,
		"fatal": logrus.FatalLevel,
		"error": logrus.ErrorLevel,
		"warn":  logrus.WarnLevel,
		"info":  logrus.InfoLevel,
		"debug": logrus.DebugLevel,
		"trace": logrus.TraceLevel,
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
