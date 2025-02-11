package main

import (
	"log"
	"os"
	"os/exec"
	"os/user"
)

var home = getUserHome()
var configPaths = []string{home + "/.config/cure/cure.yaml", "/etc/cure/cure.yaml"}

func main() {
	file, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	cmd := exec.Command("nix", "--extra-experimental-features", "nix-command",
		"--extra-experimental-features", "flakes", "build", "--no-link", "--json", "nixgl#nixGLIntel")
	cmd.Stdout = file
	cmd.Stderr = file

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func getUserHome() string {
	u, err := user.Current()
	if err != nil {
		log.Fatal("failed to get user information: ", err)
	}

	return u.HomeDir
}
