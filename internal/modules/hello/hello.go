// Package hello is for testing purposes
package hello

import (
	"context"
	"fmt"

	"github.com/sund3RRR/cure/internal/config"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
)

// Hello command
type Hello struct {
	logger *zap.Logger
	cmd    *cli.Command
}

// NewHello creates new hello
func NewHello(_ config.Config, logger *zap.Logger) *Hello {
	return &Hello{
		logger: logger,
		cmd: &cli.Command{
			Name:  "hello",
			Usage: "simple",
			Action: func(_ context.Context, cmd *cli.Command) error {
				logger.Debug("New hello command")

				if cmd.NArg() > 0 {
					fmt.Printf("Hello, %s\n", cmd.Args().Get(0))
				} else {
					fmt.Println("Empty hello")
				}

				return nil
			},
		},
	}
}

// Command is the method returning hello command
func (h *Hello) Command() *cli.Command {
	return h.cmd
}
