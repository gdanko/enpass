package output

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/gdanko/enpass/globals"
	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/hokaccha/go-prettyjson"
	"github.com/sirupsen/logrus"
)

func doJsonOutput(logger *logrus.Logger, cards []enpass.Card, nocolorFlag bool) {
	colorMap := globals.GetColorMap()
	disabledColor := false
	if nocolorFlag {
		disabledColor = true
	}
	formatter := prettyjson.NewFormatter()
	formatter.DisabledColor = disabledColor
	formatter.Indent = 4
	formatter.BoolColor = color.New(colorMap["BoolColor"])
	formatter.KeyColor = color.New(colorMap["KeyColor"])
	formatter.NullColor = color.New(colorMap["NullColor"])
	formatter.NumberColor = color.New(colorMap["NumberColor"])
	formatter.StringColor = color.New(colorMap["StringColor"])

	jsonBytes, err := formatter.Marshal(cards)
	if err != nil {
		logger.Errorf("failed to parse the output to JSON, %s", err)
		logger.Exit(2)
	}
	fmt.Println(string(jsonBytes))
}
