package output

import (
	"fmt"

	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/hokaccha/go-prettyjson"
	"github.com/sirupsen/logrus"
)

func doJsonOutput(logger *logrus.Logger, cards []enpass.Card, nocolorFlag bool) {
	disabledColor := false
	if nocolorFlag {
		disabledColor = true
	}
	formatter := prettyjson.NewFormatter()
	formatter.DisabledColor = disabledColor
	formatter.Indent = 4

	jsonBytes, err := formatter.Marshal(cards)
	if err != nil {
		logger.Errorf("failed to parse the output to JSON, %s", err)
		logger.Exit(2)
	}
	fmt.Println(string(jsonBytes))
}
