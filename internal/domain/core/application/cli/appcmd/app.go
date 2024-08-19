package appcmd

import (
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type Commander interface {
	Command() *cobra.Command
}

type AppCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
}

func NewAppCommand(i *do.Injector) (*AppCommand, error) {
	return &AppCommand{
		cmd: &cobra.Command{
			Use:     "app",
			Example: "hexago app",
			Aliases: []string{"a"},
			Short:   "Application processes",
			Long:    `Application processes`,
		},
		injector: i,
	}, nil
}

func (c *AppCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *AppCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *AppCommand) init() {
}
