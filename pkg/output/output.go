package output

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/sirupsen/logrus"
)

var (
	// colorMap = map[string]color.Attribute{
	// 	"AliasColor":  color.FgHiYellow,
	// 	"AnchorColor": color.FgHiYellow,
	// 	"BoolColor":   color.FgHiYellow,
	// 	"KeyColor":    color.FgHiCyan,
	// 	"NullColor":   color.FgHiBlack,
	// 	"NumberColor": color.FgHiMagenta,
	// 	"StringColor": color.FgHiGreen,
	// }
	colorMap = map[string]color.Attribute{
		"black-bold":   color.FgHiBlack,
		"black":        color.FgBlack,
		"blue-bold":    color.FgHiBlue,
		"blue":         color.FgBlue,
		"cyan-bold":    color.FgHiCyan,
		"cyan":         color.FgCyan,
		"green-bold":   color.FgHiGreen,
		"green":        color.FgGreen,
		"magenta-bold": color.FgHiMagenta,
		"magenta":      color.FgMagenta,
		"red-bold":     color.FgHiRed,
		"red":          color.FgRed,
		"white-bold":   color.FgHiWhite,
		"white":        color.FgWhite,
		"yellow-bold":  color.FgHiYellow,
		"yellow":       color.FgYellow,
	}
)

func GenerateOutput(logger *logrus.Logger, cmdType string, jsonFlag, listFlag, tableFlag, trashedFlag, yamlFlag, nocolorFlag bool, cards *[]enpass.Card) {
	if len(*cards) <= 0 {
		fmt.Println("No cards found matching the specified criteria")
		os.Exit(0)
	}

	// Loop through all of the cards and exclude trashed items unless we specify --trashed
	cardsPruned := []enpass.Card{}
	for _, cardItem := range *cards {
		if cardItem.IsTrashed() {
			if trashedFlag {
				cardsPruned = append(cardsPruned, cardItem)
			}
		} else {
			cardsPruned = append(cardsPruned, cardItem)
		}
	}

	// If it's a list operation, DecryptedValue should be empty
	for i := range cardsPruned {
		if cmdType == "list" {
			(cardsPruned)[i].DecryptedValue = ""
		}
	}

	for i := range cardsPruned {
		(cardsPruned)[i].Key = []byte{}
	}

	cards = &cardsPruned

	if jsonFlag {
		doJsonOutput(logger, *cards, nocolorFlag)
	} else if yamlFlag {
		doYamlOutput(logger, *cards, nocolorFlag)
	} else if listFlag {
		doListOutput(*cards, cmdType, nocolorFlag)
	} else if tableFlag {
		doTableOutput(*cards, cmdType)
	} else {
		doDefaultOutput(*cards, cmdType, nocolorFlag)
	}
}
