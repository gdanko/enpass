package output

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/gdanko/enpass/pkg/enpass"
)

func doDefaultOutput(cards []enpass.Card, cmdType string, noolorString bool) {
	var title string
	for i, cardItem := range cards {
		if cmdType == "list" {
			if noolorString {
				title = fmt.Sprintf("[%05d] >", i+1)
			} else {
				c := color.New(color.FgCyan)
				title = c.Sprintf("[%05d] >", i+1)
			}
			fmt.Printf(
				"%s title: %s, login: %s, category: %s\n",
				title,
				cardItem.Title,
				cardItem.Subtitle,
				cardItem.Category,
			)
		} else if cmdType == "show" {
			if noolorString {
				title = fmt.Sprintf("[%05d] >", i+1)
			} else {
				c := color.New(color.FgRed)
				title = c.Sprintf("[%05d] >", i+1)
			}
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
