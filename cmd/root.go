package cmd

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gdanko/enpass/globals"
	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/gdanko/enpass/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cardType         string
	caseSensitive    bool
	clipboardPrimary bool
	configPath       string
	credentials      *enpass.VaultCredentials
	defaultLogLevel  = "info"
	enpassConfig     globals.EnpassConfig
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
	nocolorFlag      bool
	nonInteractive   bool
	orderbyFlag      []string
	pinEnable        bool
	recordCategory   []string
	recordFieldLabel []string
	recordLogin      []string
	recordTitle      []string
	recordUuid       []string
	rootCmd          = &cobra.Command{
		Use:   "enpass",
		Short: "enpass is a command line interface for the Enpass password manager",
		Long:  "enpass is a command line interface for the Enpass password manager",
	}
	tableFlag     bool
	trashedFlag   bool
	vault         *enpass.Vault
	vaultPathFlag string
	versionFull   bool
	yamlFlag      bool
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	GetPersistenFlags(rootCmd)
	logLevel = logLevelMap[logLevelStr]
	logger = util.ConfigureLogger(logLevel, nocolorFlag)

	// Set the home directory in globals
	err = globals.SetHomeDirectory()
	if err != nil {
		logger.Error(err)
		logger.Exit(2)
	}

	// Parse the config file and set the config object in globals
	configPath = filepath.Join(globals.GetHomeDirectory(), ".enpass.yml")
	enpassConfig, err = util.ParseConfig(configPath)
	if err == nil {
		globals.SetConfig(enpassConfig)
	} else {
		if strings.Contains(err.Error(), "does not exist") {
			logger.Infof("%s, using the default configuration", err)
		} else if strings.Contains(err.Error(), "failed to read") || strings.Contains(err.Error(), "failed to parse") {
			logger.Warningf("%s, using the default configuration", err)
		}

		// Need to set default vault path for different OSes - there should be a better way to do this
		enpassConfig = globals.GetConfig()
		if runtime.GOOS == "darwin" {
			enpassConfig.VaultPath = "~/Library/Containers/in.sinew.Enpass-Desktop/Data/Documents/Vaults/find "
		} else if runtime.GOOS == "linux" {
			enpassConfig.VaultPath = "~/Documents/Enpass/Vaults/primary"
		}
		globals.SetConfig(enpassConfig)
	}
}
