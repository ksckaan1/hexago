package cli

import (
	"context"
	"fmt"

	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type RootCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
}

func NewRootCommand(i *do.Injector) (*RootCommand, error) {
	return &RootCommand{
		cmd: &cobra.Command{
			Use:           "hexago",
			Short:         "short description",
			Long:          header + "\nhexago is a cli tool for initializing and managing hexagonal Go projects.",
			Version:       "v0.5.0",
			SilenceUsage:  true,
			SilenceErrors: true,
		},
		injector: i,
	}, nil
}

func (c *RootCommand) Execute(ctx context.Context) error {
	err := c.cmd.ExecuteContext(ctx)
	if err != nil {
		return fmt.Errorf("cmd: execute context: %w", err)
	}
	return nil
}

func (c *RootCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}
