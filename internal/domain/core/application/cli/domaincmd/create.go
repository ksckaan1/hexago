package domaincmd

import (
	"fmt"

	"github.com/ksckaan1/hexago/internal/domain/core/port"
	"github.com/ksckaan1/hexago/internal/util"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type DomainCreateCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
}

func NewDomainCreateCommand(i *do.Injector) (*DomainCreateCommand, error) {
	return &DomainCreateCommand{
		cmd: &cobra.Command{
			Use:     "create",
			Example: "hexago domain create example\nhexago d c example\nhexago domain new example\nhexago d n example",
			Aliases: []string{"c", "new", "n"},
			Short:   "Create a domain",
			Long:    `Create a domain`,
			Args:    cobra.ExactArgs(1),
		},
		injector: i,
	}, nil
}

func (c *DomainCreateCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *DomainCreateCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *DomainCreateCommand) init() {
	c.cmd.RunE = c.runner
}

func (c *DomainCreateCommand) runner(cmd *cobra.Command, args []string) error {
	projectService, err := do.Invoke[port.ProjectService](c.injector)
	if err != nil {
		return fmt.Errorf("invoke project service: %w", err)
	}

	err = projectService.CreateDomain(cmd.Context(), args[0])
	if err != nil {
		return fmt.Errorf("project service: get all domains: %w", err)
	}

	fmt.Println("")
	util.UILog(util.Success, "domain created")
	fmt.Println("")

	return nil
}
