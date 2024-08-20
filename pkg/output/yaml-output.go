package output

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/printer"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

var (
	err        error
	yamlBytes  []byte
	yamlString string
)

func doYamlOutput(logger *logrus.Logger, cards []enpass.Card, nocolorFlag bool) {
	yamlBytes, err = yaml.Marshal(cards)
	if err != nil {
		logger.Errorf("failed to parse the output to YAML, %s", err)
		logger.Exit(2)
	}

	if nocolorFlag {
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
		p.Bool = func() *printer.Property {
			return &printer.Property{
				Prefix: format(color.FgHiMagenta),
				Suffix: format(color.Reset),
			}
		}
		p.Number = func() *printer.Property {
			return &printer.Property{
				Prefix: format(color.FgHiMagenta),
				Suffix: format(color.Reset),
			}
		}
		p.MapKey = func() *printer.Property {
			return &printer.Property{
				Prefix: format(color.FgHiCyan),
				Suffix: format(color.Reset),
			}
		}
		p.Anchor = func() *printer.Property {
			return &printer.Property{
				Prefix: format(color.FgHiYellow),
				Suffix: format(color.Reset),
			}
		}
		p.Alias = func() *printer.Property {
			return &printer.Property{
				Prefix: format(color.FgHiYellow),
				Suffix: format(color.Reset),
			}
		}
		p.String = func() *printer.Property {
			return &printer.Property{
				Prefix: format(color.FgHiGreen),
				Suffix: format(color.Reset),
			}
		}
		writer := colorable.NewColorableStdout()
		writer.Write([]byte(p.PrintTokens(tokens) + "\n"))
	}
}
