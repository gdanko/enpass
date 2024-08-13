package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

func GenerateOutput(logger *logrus.Logger, cmdType string, flags *pflag.FlagSet, cards *[]enpass.Card) {
	var (
		err        error
		jsonFlag   bool
		listFlag   bool
		yamlFlag   bool
		yamlString string
	)
	// json
	jsonFlag, err = flags.GetBool("json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// list
	listFlag, err = flags.GetBool("list")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// yaml
	yamlFlag, err = flags.GetBool("yaml")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

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
		length := 8
		for i, cardItem := range *cards {
			fmt.Printf("%*s = %s\n", length, "title", cardItem.Title)
			fmt.Printf("%*s = %s\n", length, "login", cardItem.Subtitle)
			fmt.Printf("%*s = %s\n", length, "category", cardItem.Category)
			if cmdType == "show" {
				decrypted, err := cardItem.Decrypt()
				if err != nil {
					logger.WithError(err).Error("could not decrypt " + cardItem.Title)
				} else {
					fmt.Printf("%*s = %s: %s\n", length, "type", cardItem.Type, decrypted)
				}
			}
			if i < len(*cards)-1 {
				fmt.Println()
			}
		}
	}
}
