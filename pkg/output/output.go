package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/markkurossi/tabulate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

func GenerateOutput(logger *logrus.Logger, cmdType string, flags *pflag.FlagSet, cards *[]enpass.Card) {
	var (
		err         error
		jsonFlag    bool
		listFlag    bool
		trashedFlag bool
		yamlFlag    bool
		yamlString  string
	)

	if len(*cards) <= 0 {
		fmt.Println("No cards found matching the specified criteria")
		os.Exit(0)
	}

	jsonFlag, err = flags.GetBool("json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	listFlag, err = flags.GetBool("list")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	yamlFlag, err = flags.GetBool("yaml")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	trashedFlag, err = flags.GetBool("trashed")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
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

	cards = &cardsPruned

	if jsonFlag {
		jsonBytes, _ := json.Marshal(cards)
		var prettyJSON bytes.Buffer
		if err = json.Indent(&prettyJSON, jsonBytes, "", "    "); err != nil {
			fmt.Printf("failed to parse the output as JSON: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(prettyJSON.String())

	} else if yamlFlag {
		yamlBytes, _ := yaml.Marshal(cards)
		yamlString = string(yamlBytes)
		yamlString = strings.TrimSpace(yamlString)
		fmt.Println(yamlString)
	} else if listFlag {
		length := 10
		for i, cardItem := range *cards {
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
	} else {
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
	}
}
