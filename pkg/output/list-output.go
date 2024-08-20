package output

import (
	"fmt"

	"github.com/gdanko/enpass/pkg/enpass"
)

func doListOutput(cards []enpass.Card, cmdType string) {
	length := 10
	for i, cardItem := range cards {
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
		if i < len(cards)-1 {
			fmt.Println()
		}
	}
}
