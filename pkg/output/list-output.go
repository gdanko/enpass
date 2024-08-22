package output

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gdanko/enpass/pkg/enpass"
)

func doListOutput(cards []enpass.Card, cmdType string, nocolorFlag bool) {
	for i, cardItem := range cards {
		if nocolorFlag {
			fmt.Printf("%s = %s\n", "      uuid", cardItem.UUID)
			fmt.Printf("%s = %d\n", "   created", cardItem.CreatedAt)
			fmt.Printf("%s = %s\n", " card_type", cardItem.Type)
			fmt.Printf("%s = %s\n", "     title", cardItem.Title)
			fmt.Printf("%s = %s\n", "     login", cardItem.Subtitle)
			if cardItem.Note != "" {
				fmt.Printf("%s = %s\n", "      note", cardItem.Note)
			}
			fmt.Printf("%s = %s\n", "  category", cardItem.Category)
			fmt.Printf("%s = %s\n", "     label", cardItem.Label)
			fmt.Printf("%s = %d\n", " last_used", cardItem.LastUsed)
			fmt.Printf("%s = %v\n", " sensitive", cardItem.Sensitive)
			fmt.Printf("%s = %v\n", "      icon", cardItem.Icon)
			if cmdType == "show" {
				fmt.Printf("%s = %s: %s\n", "type", cardItem.Type, cardItem.DecryptedValue)
			}
		} else {
			var (
				keyColor    = color.New(color.FgHiCyan).SprintFunc()
				valueString = color.New(color.FgHiGreen).SprintFunc()
				valueBool   = color.New(color.FgHiYellow).SprintFunc()
				valueNumber = color.New(color.FgHiMagenta).SprintFunc()
			)
			fmt.Printf("%s = %s\n", keyColor("      uuid"), valueString(cardItem.UUID))
			fmt.Printf("%s = %s\n", keyColor("   created"), valueNumber(cardItem.CreatedAt))
			fmt.Printf("%s = %s\n", keyColor(" card_type"), valueString(cardItem.Type))
			fmt.Printf("%s = %s\n", keyColor("     title"), valueString(cardItem.Title))
			fmt.Printf("%s = %s\n", keyColor("  subtitle"), valueString(cardItem.Subtitle))
			if cardItem.Note != "" {
				fmt.Printf("%s = %s\n", keyColor("      note"), valueString(cardItem.Note))
			}
			fmt.Printf("%s = %s\n", keyColor("  category"), valueString(cardItem.Category))
			fmt.Printf("%s = %s\n", keyColor("     label"), valueString(cardItem.Label))
			fmt.Printf("%s = %s\n", keyColor(" last_used"), valueNumber(cardItem.LastUsed))
			fmt.Printf("%s = %s\n", keyColor(" sensitive"), valueBool(cardItem.Sensitive))
			fmt.Printf("%s = %s\n", keyColor("      icon"), valueString(cardItem.Icon))
		}
		if i < len(cards)-1 {
			fmt.Println()
		}
	}
	os.Exit(0)
}
