package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/markkurossi/tabulate"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

func GenerateOutput(flags *pflag.FlagSet, output []map[string]interface{}) {
	var (
		err        error
		fields     []string
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
		jsonBytes, _ := json.Marshal(output)
		var prettyJSON bytes.Buffer
		if err = json.Indent(&prettyJSON, jsonBytes, "", "    "); err != nil {
			fmt.Printf("failed to parse the output as JSON: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(prettyJSON.String())

	} else if yamlFlag {
		yamlBytes, _ := yaml.Marshal(output)
		yamlString = string(yamlBytes)
		yamlString = strings.TrimSpace(yamlString)
		fmt.Println(yamlString)

	} else if listFlag {
		fields = globals.GetFields()
		length := 0
		for key, _ := range output[0] {
			if len(key) > length {
				length = len(key)
			}
		}
		for i, amiItem := range output {
			for _, field := range fields {
				valueType := fmt.Sprint(reflect.TypeOf(amiItem[field]))
				if field == "created" || field == "updated" {
					humanDate, _ := amiItem[field].(int64)
					fmt.Printf("%*s = %s\n", length, field, util.TimestampToHuman(humanDate))
				} else if valueType == "<nil>" {
					fmt.Printf("%*s = %s\n", length, field, globals.GetUnknown())
				} else if valueType == "string" {
					if amiItem[field] == "" {
						fmt.Printf("%*s = %s\n", length, field, globals.GetUnknown())
					} else {
						fmt.Printf("%*s = %s\n", length, field, fmt.Sprint(amiItem[field]))
					}
				} else if valueType == "int64" {
					fmt.Printf("%*s = %d\n", length, field, amiItem[field].(int64))
				}
			}
			if i < len(output)-1 {
				fmt.Println()
			}
		}

	} else {
		fields = globals.GetFields()
		tab := tabulate.New(tabulate.Simple)
		for _, field := range fields {
			tab.Header(field).SetAlign(tabulate.ML)
		}
		for _, amiItem := range output {
			row := tab.Row()
			for _, field := range fields {
				valueType := fmt.Sprint(reflect.TypeOf(amiItem[field]))
				if field == "created" || field == "updated" {
					humanDate, _ := amiItem[field].(int64)
					row.Column(fmt.Sprint(util.TimestampToHuman(humanDate)))
				} else if valueType == "<nil>" {
					row.Column(globals.GetUnknown())
				} else if valueType == "string" {
					if amiItem[field] == "" {
						row.Column(globals.GetUnknown())
					} else {
						row.Column(fmt.Sprint(amiItem[field]))
					}
				} else {
					row.Column(fmt.Sprint(amiItem[field]))
				}
			}
		}
		tab.Print(os.Stdout)
	}
}
