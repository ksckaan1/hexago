package domaincmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/ksckaan1/hexago/internal/domain/core/dto"
	"github.com/ksckaan1/hexago/internal/domain/core/port"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type DomainCreateCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
	tuilog   *tuilog.TUILog
}

const newLong = `new command creates a domain under the "internal/domain/" directory.`

func NewDomainCreateCommand(i *do.Injector) (*DomainCreateCommand, error) {
	return &DomainCreateCommand{
		cmd: &cobra.Command{
			Use:     "new",
			Example: "hexago domain new",
			Short:   "Create a domain",
			Long:    newLong,
		},
		injector: i,
		tuilog:   do.MustInvoke[*tuilog.TUILog](i),
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

	var domainName string

	if len(args) > 0 {
		domainName = args[0]
	}

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What’s domain name?").
				Placeholder("domainname").
				Validate(projectService.ValidatePkgName).
				Description("Domain name must be lowercase").
				Value(&domainName),
		).WithShowHelp(true),
	).Run()
	if err != nil {
		return fmt.Errorf("input domain name: %w", err)
	}

	err = projectService.CreateDomain(
		cmd.Context(),
		dto.CreateDomainParams{
			DomainName: domainName,
		},
	)
	if err != nil {
		fmt.Println("")
		c.tuilog.Error(err.Error())
		fmt.Println("")
		return fmt.Errorf("project service: get all domains: %w", err)
	}

	fmt.Println("")
	c.tuilog.Success("Domain created")
	fmt.Println("")

	return nil
}
