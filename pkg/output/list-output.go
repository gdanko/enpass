package output

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/gdanko/enpass/globals"
	"github.com/gdanko/enpass/pkg/enpass"
)

func doListOutput(cards []enpass.Card, cmdType string, flagNoColor bool) {
	for i, cardItem := range cards {
		if flagNoColor {
			fmt.Printf("%s = %s\n", "           uuid", cardItem.UUID)
			fmt.Printf("%s = %s\n", "        created", cardItem.Created)
			fmt.Printf("%s = %s\n", "        updated", cardItem.Updated)
			fmt.Printf("%s = %s\n", "      card_type", cardItem.Type)
			fmt.Printf("%s = %s\n", "          title", cardItem.Title)
			fmt.Printf("%s = %s\n", "          login", cardItem.Subtitle)
			if cardItem.Note != "" {
				fmt.Printf("%s = %s\n", "           note", cardItem.Note)
			}
			fmt.Printf("%s = %s\n", "       category", cardItem.Category)
			fmt.Printf("%s = %s\n", "          label", cardItem.Label)
			fmt.Printf("%s = %s\n", "      last_used", cardItem.LastUsed)
			fmt.Printf("%s = %v\n", "      sensitive", cardItem.Sensitive)
			fmt.Printf("%s = %v\n", "           icon", cardItem.Icon)
			if cmdType == "show" {
				fmt.Printf("%s = %s: %s\n", "decrypted_value", cardItem.Type, cardItem.DecryptedValue)
			}
		} else {
			var (
				boolColor   = color.New(colorMap[globals.GetConfig().Colors.BoolColor]).SprintFunc()
				keyColor    = color.New(colorMap[globals.GetConfig().Colors.KeyColor]).SprintFunc()
				numberColor = color.New(colorMap[globals.GetConfig().Colors.NumberColor]).SprintFunc()
				stringColor = color.New(colorMap[globals.GetConfig().Colors.StringColor]).SprintFunc()
			)
			fmt.Printf("%s = %s\n", keyColor("           uuid"), stringColor(cardItem.UUID))
			fmt.Printf("%s = %s\n", keyColor("        created"), numberColor(cardItem.Created))
			fmt.Printf("%s = %s\n", keyColor("        updated"), numberColor(cardItem.Updated))
			fmt.Printf("%s = %s\n", keyColor("      card_type"), stringColor(cardItem.Type))
			fmt.Printf("%s = %s\n", keyColor("          title"), stringColor(cardItem.Title))
			fmt.Printf("%s = %s\n", keyColor("       subtitle"), stringColor(cardItem.Subtitle))
			if cardItem.Note != "" {
				fmt.Printf("%s = %s\n", keyColor("           note"), stringColor(cardItem.Note))
			}
			fmt.Printf("%s = %s\n", keyColor("       category"), stringColor(cardItem.Category))
			fmt.Printf("%s = %s\n", keyColor("          label"), stringColor(cardItem.Label))
			fmt.Printf("%s = %s\n", keyColor("      last_used"), numberColor(cardItem.LastUsed))
			fmt.Printf("%s = %s\n", keyColor("      sensitive"), boolColor(cardItem.Sensitive))
			fmt.Printf("%s = %s\n", keyColor("           icon"), stringColor(cardItem.Icon))
			if cmdType == "show" {
				fmt.Printf("%s = %s\n", keyColor("decrypted_value"), stringColor(cardItem.DecryptedValue))
			}
		}
		if i < len(cards)-1 {
			fmt.Println()
		}
	}
}
