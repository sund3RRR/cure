package main

import (
	"fmt"
	"log"
	"os/user"

	"github.com/sund3RRR/cure/internal/config"
	"github.com/sund3RRR/cure/pkg/adapters/nix"
)

var home = getUserHome()
var configPaths = []string{home + "/.config/cure/cure.yaml", "/etc/cure/cure.yaml"}

func main() {
	// Create main config
	cfg := config.NewConfig(configPaths...)

	// Init logger
	logger, err := cfg.Logger.Build()
	if err != nil {
		log.Fatal("failed to create logger: ", err)
	}
	defer logger.Sync() //nolint

	// Create adapters
	nixAdapter := nix.NewNix(logger)

	pi, err := nixAdapter.GetPackage("nixpkgs", "micro")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pi)
}

func getUserHome() string {
	u, err := user.Current()
	if err != nil {
		log.Fatal("failed to get user information: ", err)
	}

	return u.HomeDir
}
