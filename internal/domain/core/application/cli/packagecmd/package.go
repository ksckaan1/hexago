package packagecmd

import (
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type Commander interface {
	Command() *cobra.Command
}

type PackageCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
}

func NewPackageCommand(i *do.Injector) (*PackageCommand, error) {
	return &PackageCommand{
		cmd: &cobra.Command{
			Use:     "pkg",
			Example: "hexago pkg",
			Short:   "Package processes",
			Long:    `Package processes`,
		},
		injector: i,
	}, nil
}

func (c *PackageCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *PackageCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *PackageCommand) init() {
}
