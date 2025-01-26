// Package app is the cure entrypoint
package app

import (
	"github.com/urfave/cli/v3"
)

// Command interface
type Command interface {
	Command() *cli.Command
}

// App is the entrypoint to
type App struct {
	commands []Command
}

// NewApp created an app
func NewApp(commands ...Command) *App {
	return &App{
		commands: commands,
	}
}

// GetCommands returns the slice with all available commands
func (app *App) GetCommands() []*cli.Command {
	commands := make([]*cli.Command, 0, len(app.commands))

	for _, cmd := range app.commands {
		commands = append(commands, cmd.Command())
	}

	return commands
}
