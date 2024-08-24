package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gdanko/enpass/globals"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"gopkg.in/yaml.v3"
)

// FileOrDirectoryExists : Determine if a file or directory exists
func FileOrDirectoryExists(path string) (exists bool, err error) {
	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, err
	}

	return false, err
}

// ReturnLogLevels : Return a comma-delimited list of log levels
func ReturnLogLevels(levelMap map[string]logrus.Level) string {
	logLevels := make([]string, 0, len(levelMap))
	for k := range levelMap {
		logLevels = append(logLevels, k)
	}
	sort.Strings(logLevels)

	return strings.Join(logLevels, ", ")
}

// ConfigureLogger : Configure the logger
func ConfigureLogger(logLevel logrus.Level, nocolorFlag bool) (logger *logrus.Logger) {
	disableColors := false
	if nocolorFlag {
		disableColors = true
	}
	logger = &logrus.Logger{
		Out:   os.Stderr,
		Level: logLevel,
		Formatter: &prefixed.TextFormatter{
			DisableColors:    disableColors,
			DisableTimestamp: true,
			TimestampFormat:  "2006-01-02 15:04:05",
			FullTimestamp:    true,
			ForceFormatting:  false,
		},
	}
	logger.SetLevel(logLevel)

	return logger
}

// ExpandPath : Expand paths starting with ~/
func ExpandPath(path string) (expanded string) {
	if path == "~" {
		expanded = globals.GetHomeDirectory()
	} else if strings.HasPrefix(path, "~/") {
		expanded = filepath.Join(globals.GetHomeDirectory(), path[2:])
	} else {
		expanded = path
	}

	return expanded
}

// ParseConfig : Read the config file and return it as an EnpassConfig object
func ParseConfig(path string) (enpassConfig globals.EnpassConfig, err error) {
	expanded := ExpandPath(path)
	exists, err := FileOrDirectoryExists(ExpandPath(expanded))
	if !exists && err != nil {
		return globals.EnpassConfig{}, fmt.Errorf("the config file %s does not exist", expanded)
	}

	data, err := ioutil.ReadFile(expanded)
	if err != nil {
		return globals.EnpassConfig{}, fmt.Errorf("failed to read the config file %s: %s", expanded, err)
	}

	err = yaml.Unmarshal(data, &enpassConfig)
	if err != nil {
		return globals.EnpassConfig{}, fmt.Errorf("failed to parse the config file %s: %s", expanded, err)
	}

	return enpassConfig, nil
}
