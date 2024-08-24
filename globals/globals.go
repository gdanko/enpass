package globals

import (
	"fmt"
	"os/user"
	"sync"
)

type Colors struct {
	AliasColor  string `yaml:"alias_color"`
	AnchorColor string `yaml:"anchor_color"`
	BoolColor   string `yaml:"bool_color"`
	KeyColor    string `yaml:"key_color"`
	NullColor   string `yaml:"null_color"`
	NumberColor string `yaml:"number_color"`
	StringColor string `yaml:"string_color"`
}

type EnpassConfig struct {
	Colors        Colors `yaml:"colors"`
	OutputStyle   string `yaml:"output_style"`
	VaultPassword string `yaml:"vault_password"`
	VaultPath     string `yaml:"vault_path"`
}

var (
	enpassConfig = EnpassConfig{
		Colors: Colors{
			AliasColor:  "yellow-bold",
			AnchorColor: "yellow-bold",
			BoolColor:   "yellow-bold",
			KeyColor:    "cyan-bold",
			NullColor:   "black-bold",
			NumberColor: "magenta-bold",
			StringColor: "green-bold",
		},
		VaultPath: "~/Documents/Enpass/Vaults/primary",
	}
	homeDir string
	mu      sync.RWMutex
)

// Set and get pairs
func SetConfig(x EnpassConfig) {
	mu.Lock()
	enpassConfig = x
	mu.Unlock()
}

func GetConfig() (x EnpassConfig) {
	mu.Lock()
	x = enpassConfig
	mu.Unlock()
	return x
}

func SetHomeDirectory() (err error) {
	mu.Lock()
	userObj, err := user.Current()
	if err != nil {
		mu.Unlock()
		return fmt.Errorf("failed to determine the path of your home directory: %s", err)
	}
	homeDir = userObj.HomeDir
	mu.Unlock()
	return nil
}

func GetHomeDirectory() (x string) {
	mu.Lock()
	x = homeDir
	mu.Unlock()
	return x
}
