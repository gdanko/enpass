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
	disabledColor := false
	if nocolorFlag {
		disabledColor = true
	}
	formatter := prettyjson.NewFormatter()
	formatter.DisabledColor = disabledColor
	formatter.Indent = 4
	formatter.BoolColor = color.New(colorMap[globals.GetConfig().Colors.BoolColor])
	formatter.KeyColor = color.New(colorMap[globals.GetConfig().Colors.KeyColor])
	formatter.NullColor = color.New(colorMap[globals.GetConfig().Colors.NullColor])
	formatter.NumberColor = color.New(colorMap[globals.GetConfig().Colors.NumberColor])
	formatter.StringColor = color.New(colorMap[globals.GetConfig().Colors.StringColor])

	jsonBytes, err := formatter.Marshal(cards)
	if err != nil {
		logger.Errorf("failed to parse the output to JSON, %s", err)
		logger.Exit(2)
	}
	fmt.Println(string(jsonBytes))
}
