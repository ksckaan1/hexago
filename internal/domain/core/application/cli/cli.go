package cli

import (
	"context"
	"fmt"

	"github.com/samber/do"
	"github.com/spf13/cobra"
)

type Commander interface {
	Command() *cobra.Command
}

type CLI struct {
	rootCmd *RootCommand
	initCmd *InitCommand
}

func NewCLI(i *do.Injector) (*CLI, error) {
	return &CLI{
		rootCmd: do.MustInvoke[*RootCommand](i),
		initCmd: do.MustInvoke[*InitCommand](i),
	}, nil
}

func (c *CLI) Run(ctx context.Context) error {
	c.rootCmd.AddCommand(c.initCmd)

	err := c.rootCmd.Execute(ctx)
	if err != nil {
		return fmt.Errorf("root cmd: execute: %w", err)
	}
	return nil
}
