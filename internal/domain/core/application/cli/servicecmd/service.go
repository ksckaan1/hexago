package servicecmd

import (
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type Commander interface {
	Command() *cobra.Command
}

type ServiceCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
}

func NewServiceCommand(i *do.Injector) (*ServiceCommand, error) {
	return &ServiceCommand{
		cmd: &cobra.Command{
			Use:     "service",
			Example: "hexago service",
			Aliases: []string{"s"},
			Short:   "Service processes",
			Long:    `Service processes`,
		},
		injector: i,
	}, nil
}

func (c *ServiceCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *ServiceCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *ServiceCommand) init() {
}
