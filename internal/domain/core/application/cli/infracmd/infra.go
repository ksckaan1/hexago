package infracmd

import (
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type Commander interface {
	Command() *cobra.Command
}

type InfraCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
}

func NewInfraCommand(i *do.Injector) (*InfraCommand, error) {
	return &InfraCommand{
		cmd: &cobra.Command{
			Use:     "infra",
			Example: "hexago infra",
			Short:   "Infrastructure processes",
			Long:    `Infrastructure processes`,
		},
		injector: i,
	}, nil
}

func (c *InfraCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *InfraCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *InfraCommand) init() {
}
