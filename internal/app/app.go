package app

import (
	"github.com/sund3RRR/cure/internal/config"
	"github.com/sund3RRR/cure/internal/modules/hello"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
)

type App struct {
	hello *hello.Hello
}

func NewApp(cfg config.Config, logger *zap.Logger) *App {
	return &App{
		hello: hello.NewHello(cfg, logger),
	}
}

func (app *App) GetCommands() []*cli.Command {
	return []*cli.Command{
		app.hello.Command(),
	}
}
