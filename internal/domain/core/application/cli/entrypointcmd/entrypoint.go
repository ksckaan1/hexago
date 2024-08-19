package entrypointcmd

import (
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type Commander interface {
	Command() *cobra.Command
}

type EntryPointCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
}

func NewEntryPointCommand(i *do.Injector) (*EntryPointCommand, error) {
	return &EntryPointCommand{
		cmd: &cobra.Command{
			Use:     "cmd",
			Example: "hexago cmd",
			Short:   "Entry Point processes",
			Long:    `Entry Point processes`,
		},
		injector: i,
	}, nil
}

func (c *EntryPointCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *EntryPointCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *EntryPointCommand) init() {
}
