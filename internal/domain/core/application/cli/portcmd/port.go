package portcmd

import (
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type Commander interface {
	Command() *cobra.Command
}

type PortCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
}

func NewPortCommand(i *do.Injector) (*PortCommand, error) {
	return &PortCommand{
		cmd: &cobra.Command{
			Use:     "port",
			Example: "hexago port",
			Short:   "Port processes",
			Long:    `Port processes`,
		},
		injector: i,
	}, nil
}

func (c *PortCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *PortCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *PortCommand) init() {
}
