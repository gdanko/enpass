package output

import (
	"fmt"
	"os"

	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/sirupsen/logrus"
)

const escape = "\x1b"

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
		doListOutput(*cards, cmdType)
	} else if tableFlag {
		doTableOutput(*cards, cmdType)
	} else {
		doDefaultOutput(*cards, cmdType)
	}
}
