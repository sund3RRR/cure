// Package app represents cure api
package app

import (
	"github.com/sund3RRR/cure/internal/config"
	"github.com/sund3RRR/cure/internal/modules/hello"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
)

// App is the entrypoint to
type App struct {
	hello *hello.Hello
}

// NewApp created an app
func NewApp(cfg config.Config, logger *zap.Logger) *App {
	return &App{
		hello: hello.NewHello(cfg, logger),
	}
}

// GetCommands returns the slice with all available commands
func (app *App) GetCommands() []*cli.Command {
	return []*cli.Command{
		app.hello.Command(),
	}
}
