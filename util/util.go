package util

import (
	"os"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func ConfigureLogger(logLevel logrus.Level, nocolorFlag bool) (logger *logrus.Logger) {
	disableColors := false
	if nocolorFlag {
		disableColors = true
	}
	logger = &logrus.Logger{
		Out:   os.Stderr,
		Level: logLevel,
		// Formatter: &easy.Formatter{
		// 	TimestampFormat: "2006-01-02 15:04:05",
		// 	LogFormat:       "[%lvl%]: %time% - %msg%",
		// },
		Formatter: &prefixed.TextFormatter{
			DisableColors:    disableColors,
			DisableTimestamp: true,
			TimestampFormat:  "2006-01-02 15:04:05",
			FullTimestamp:    true,
			ForceFormatting:  true,
		},
	}

	logger.SetLevel(logLevel)
	return logger
}

func ReturnLogLevels(levelMap map[string]logrus.Level) string {
	logLevels := make([]string, 0, len(levelMap))
	for k := range levelMap {
		logLevels = append(logLevels, k)
	}
	sort.Strings(logLevels)
	return strings.Join(logLevels, ", ")
}
