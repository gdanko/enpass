package output

import (
	"fmt"
	"os"

	"github.com/gdanko/enpass/pkg/enpass"
	"github.com/markkurossi/tabulate"
)

func doTableOutput(cards []enpass.Card, cmdType string) {
	tab := tabulate.New(tabulate.Simple)
	tab.Header("title").SetAlign(tabulate.ML)
	tab.Header("login").SetAlign(tabulate.ML)
	tab.Header("category").SetAlign(tabulate.ML)
	if cmdType == "show" {
		tab.Header("decrypted").SetAlign(tabulate.ML)
	}
	for _, cardItem := range cards {
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
