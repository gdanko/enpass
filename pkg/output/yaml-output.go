package output

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/gdanko/enpass/globals"
	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/printer"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

const escape = "\x1b"

var (
	err        error
	yamlBytes  []byte
	yamlString string
)

func format(attr color.Attribute) string {
	return fmt.Sprintf("%s[%dm", escape, attr)
}

func doYamlOutput(logger *logrus.Logger, cards []enpass.Card, flagNoColor bool) {
	yamlBytes, err = yaml.Marshal(cards)
	if err != nil {
		logger.Errorf("failed to parse the output to YAML, %s", err)
		logger.Exit(2)
	}

	if flagNoColor {
		yamlString = string(yamlBytes)
		yamlString = strings.TrimSpace(yamlString)
		fmt.Println(yamlString)
	} else {
		tokens := lexer.Tokenize(string(yamlBytes))
		var p printer.Printer
		p.LineNumber = false
		p.LineNumberFormat = func(num int) string {
			fn := color.New(color.Bold, color.FgHiWhite).SprintFunc()
			return fn(fmt.Sprintf("%2d | ", num))
		}
		p.Alias = func() *printer.Property {
			return &printer.Property{
				Prefix: format(colorMap[globals.GetConfig().Colors.AliasColor]),
				Suffix: format(color.Reset),
			}
		}
		p.Anchor = func() *printer.Property {
			return &printer.Property{
				Prefix: format(colorMap[globals.GetConfig().Colors.AnchorColor]),
				Suffix: format(color.Reset),
			}
		}
		p.Bool = func() *printer.Property {
			return &printer.Property{
				Prefix: format(colorMap[globals.GetConfig().Colors.BoolColor]),
				Suffix: format(color.Reset),
			}
		}
		p.MapKey = func() *printer.Property {
			return &printer.Property{
				Prefix: format(colorMap[globals.GetConfig().Colors.KeyColor]),
				Suffix: format(color.Reset),
			}
		}
		p.Number = func() *printer.Property {
			return &printer.Property{
				Prefix: format(colorMap[globals.GetConfig().Colors.NumberColor]),
				Suffix: format(color.Reset),
			}
		}
		p.String = func() *printer.Property {
			return &printer.Property{
				Prefix: format(colorMap[globals.GetConfig().Colors.StringColor]),
				Suffix: format(color.Reset),
			}
		}
		writer := colorable.NewColorableStdout()
		writer.Write([]byte(p.PrintTokens(tokens) + "\n"))
	}
}
