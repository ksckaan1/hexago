package domaincmd

import (
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type Commander interface {
	Command() *cobra.Command
}

type DomainCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
}

func NewDomainCommand(i *do.Injector) (*DomainCommand, error) {
	return &DomainCommand{
		cmd: &cobra.Command{
			Use:     "domain",
			Example: "hexago domain",
			Aliases: []string{"d"},
			Short:   "Domain processes",
			Long:    `Domain processes`,
		},
		injector: i,
	}, nil
}

func (c *DomainCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *DomainCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *DomainCommand) init() {
}
