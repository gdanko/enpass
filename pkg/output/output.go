package output

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/hokaccha/go-prettyjson"
	"github.com/markkurossi/tabulate"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func GenerateOutput(logger *logrus.Logger, cmdType string, jsonFlag, listFlag, tableFlag, trashedFlag, yamlFlag, nocolorFlag bool, cards *[]enpass.Card) {
	var (
		// err        error
		// jsonBytes  []byte
		yamlString string
	)

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

	} else if yamlFlag {
		yamlBytes, _ := yaml.Marshal(cards)
		yamlString = string(yamlBytes)
		yamlString = strings.TrimSpace(yamlString)
		fmt.Println(yamlString)
	} else if listFlag {
		length := 10
		for i, cardItem := range *cards {
			fmt.Printf("%*s = %s\n", length, "uuid", cardItem.UUID)
			fmt.Printf("%*s = %s\n", length, "title", cardItem.Title)
			fmt.Printf("%*s = %s\n", length, "login", cardItem.Subtitle)
			fmt.Printf("%*s = %s\n", length, "category", cardItem.Category)
			fmt.Printf("%*s = %s\n", length, "note", cardItem.Note)
			fmt.Printf("%*s = %v\n", length, "sensitive", cardItem.Sensitive)
			fmt.Printf("%*s = %v\n", length, "raw", cardItem.RawValue)
			if cmdType == "show" {
				fmt.Printf("%*s = %s: %s\n", length, "type", cardItem.Type, cardItem.DecryptedValue)
			}
			if i < len(*cards)-1 {
				fmt.Println()
			}
		}
	} else if tableFlag {
		tab := tabulate.New(tabulate.Simple)
		tab.Header("title").SetAlign(tabulate.ML)
		tab.Header("login").SetAlign(tabulate.ML)
		tab.Header("category").SetAlign(tabulate.ML)
		if cmdType == "show" {
			tab.Header("decrypted").SetAlign(tabulate.ML)
		}
		for _, cardItem := range *cards {
			row := tab.Row()
			row.Column(cardItem.Title)
			row.Column(cardItem.Subtitle)
			row.Column(cardItem.Category)
			if cmdType == "show" {
				password := fmt.Sprintf("%s: %s", cardItem.Type, cardItem.DecryptedValue)
				row.Column(password)
			}
		}
		tab.Print(os.Stdout)
	} else {
		for i, cardItem := range *cards {
			if cmdType == "list" {
				c := color.New(color.FgCyan)
				title := c.Sprintf("[%05d] >", i+1)
				fmt.Printf(
					"%s title: %s, login: %s, category: %s\n",
					title,
					cardItem.Title,
					cardItem.Subtitle,
					cardItem.Category,
				)
			} else if cmdType == "show" {
				c := color.New(color.FgRed)
				title := c.Sprintf("[%05d] >", i+1)
				fmt.Printf(
					"%s title: %s, login: %s, category: %s, %s: %s\n",
					title,
					cardItem.Title,
					cardItem.Subtitle,
					cardItem.Category,
					cardItem.Type,
					cardItem.DecryptedValue,
				)
			}
		}
	}
}
