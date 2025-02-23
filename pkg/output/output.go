package output

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gdanko/enpass/globals"
	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/sirupsen/logrus"
)

var (
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

func GenerateOutput(logger *logrus.Logger, cmdType string, flagList, flagTable, flagTrashed, flagYaml, flagNoColor bool, cards *[]enpass.Card) {
	if len(*cards) <= 0 {
		fmt.Println("No records found matching the specified criteria")
		os.Exit(0)
	}

	// Loop through all of the cards and exclude trashed items unless we specify --trashed
	cardsPruned := []enpass.Card{}
	for _, cardItem := range *cards {
		if cardItem.IsTrashed() {
			if flagTrashed {
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

	if flagList {
		doListOutput(*cards, cmdType, flagNoColor)
	} else if flagTable {
		doTableOutput(*cards, cmdType)
	} else if flagYaml {
		doYamlOutput(logger, *cards, flagNoColor)
	} else {
		outputStyle := globals.GetConfig().OutputStyle
		switch outputStyle {
		case "list":
			doListOutput(*cards, cmdType, flagNoColor)
		case "table":
			doTableOutput(*cards, cmdType)
		case "yaml":
			doYamlOutput(logger, *cards, flagNoColor)
		default:
			doDefaultOutput(*cards, cmdType, flagNoColor)
		}
	}
}
