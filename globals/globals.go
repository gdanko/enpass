package globals

import (
	"sync"

	"github.com/fatih/color"
)

var (
	colorMap = map[string]color.Attribute{
		"AliasColor":  color.FgHiYellow,
		"AnchorColor": color.FgHiYellow,
		"BoolColor":   color.FgHiYellow,
		"KeyColor":    color.FgHiCyan,
		"NullColor":   color.FgHiBlack,
		"NumberColor": color.FgHiMagenta,
		"StringColor": color.FgHiGreen,
	}
	mu sync.RWMutex
)

func GetColorMap() (x map[string]color.Attribute) {
	mu.Lock()
	x = colorMap
	mu.Unlock()
	return x
}
