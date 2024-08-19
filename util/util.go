package util

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/gdanko/enpass/pkg/unlock"
	"github.com/miquella/ask"
	"github.com/sirupsen/logrus"
)

const (
	pinMinLength           = 8
	pinDefaultKdfIterCount = 100000
)

func prompt(logger *logrus.Logger, nonInteractive bool, msg string) string {
	if !nonInteractive {
		if response, err := ask.HiddenAsk("Enter " + msg + ": "); err != nil {
			logger.WithError(err).Fatal("could not prompt for " + msg)
		} else {
			return response
		}
	}
	return ""
}

func InitializeStore(logger *logrus.Logger, vaultPath string, nonInteractive bool) *unlock.SecureStore {
	vaultPath, _ = filepath.EvalSymlinks(vaultPath)
	store, err := unlock.NewSecureStore(filepath.Base(vaultPath), logger.Level)
	if err != nil {
		logger.WithError(err).Fatal("could not create store")
	}

	pin := os.Getenv("ENP_PIN")
	if pin == "" {
		pin = prompt(logger, nonInteractive, "PIN")
	}
	if len(pin) < pinMinLength {
		logger.Fatal("PIN too short")
	}

	pepper := os.Getenv("ENP_PIN_PEPPER")

	pinKdfIterCount, err := strconv.ParseInt(os.Getenv("ENP_PIN_ITER_COUNT"), 10, 32)
	if err != nil {
		pinKdfIterCount = pinDefaultKdfIterCount
	}

	if err := store.GeneratePassphrase(pin, pepper, int(pinKdfIterCount)); err != nil {
		logger.WithError(err).Fatal("could not initialize store")
	}

	return store
}

func AssembleVaultCredentials(logger *logrus.Logger, vaultPath string, keyFilePath string, nonInteractive bool, store *unlock.SecureStore) *enpass.VaultCredentials {
	credentials := &enpass.VaultCredentials{
		Password:    os.Getenv("MASTERPW"),
		KeyfilePath: keyFilePath,
	}

	if !credentials.IsComplete() && store != nil {
		var err error
		if credentials.DBKey, err = store.Read(); err != nil {
			logger.WithError(err).Fatal("could not read credentials from store")
		}
		logger.Debug("read credentials from store")
	}

	if !credentials.IsComplete() {
		credentials.Password = prompt(logger, nonInteractive, "vault password")
	}

	return credentials
}

func OpenVault(logger *logrus.Logger, pinEnable bool, nonInteractive bool, vaultPath string, keyFilePath string, logLevel logrus.Level) (vault *enpass.Vault, credentials *enpass.VaultCredentials, err error) {
	vault, err = enpass.NewVault(vaultPath, logLevel)
	if err != nil {
		panic(err)
	}

	var store *unlock.SecureStore
	if !pinEnable {
		logger.Debug("PIN disabled")
	} else {
		logger.Debug("PIN enabled, using store")
		store = InitializeStore(logger, vaultPath, nonInteractive)
		logger.Debug("initialized store")
	}
	credentials = AssembleVaultCredentials(logger, vaultPath, keyFilePath, nonInteractive, store)

	return vault, credentials, nil
}

func ReturnLogLevels(levelMap map[string]logrus.Level) string {
	logLevels := make([]string, 0, len(levelMap))
	for k := range levelMap {
		logLevels = append(logLevels, k)
	}
	sort.Strings(logLevels)
	return strings.Join(logLevels, ", ")
}
