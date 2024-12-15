package domaincmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/internal/customerrors"
	"github.com/ksckaan1/hexago/internal/domain/core/model"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/ksckaan1/hexago/internal/port"
)

var _ port.Commander = (*DomainCreateCommand)(nil)

type DomainCreateCommand struct {
	cmd            *cobra.Command
	projectService ProjectService
	tuilog         *tuilog.TUILog
}

const newLong = `new command creates a domain under the "internal/domain/" directory.`

func NewDomainCreateCommand(projectService ProjectService, tl *tuilog.TUILog) (*DomainCreateCommand, error) {
	return &DomainCreateCommand{
		cmd: &cobra.Command{
			Use:     "new",
			Example: "hexago domain new",
			Short:   "Create a domain",
			Long:    newLong,
		},
		projectService: projectService,
		tuilog:         tl,
	}, nil
}

func (c *DomainCreateCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *DomainCreateCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}

func (c *DomainCreateCommand) init() {
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := c.runner(cmd, args)
		if err != nil {
			return customerrors.ErrSuppressed
		}
		return nil
	}
}

func (c *DomainCreateCommand) runner(cmd *cobra.Command, args []string) error {
	var domainName string

	if len(args) > 0 {
		domainName = args[0]
	}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What’s domain name?").
				Placeholder("domainname").
				Validate(c.projectService.ValidatePkgName).
				Description("Domain name must be lowercase").
				Value(&domainName),
		).WithShowHelp(true),
	).Run()
	if err != nil {
		return fmt.Errorf("input domain name: %w", err)
	}

	err = c.projectService.CreateDomain(
		cmd.Context(),
		model.CreateDomainParams{
			DomainName: domainName,
		},
	)
	if err != nil {

		c.tuilog.Error(err.Error())

		return fmt.Errorf("projectService.CreateDomain: %w", err)
	}

	c.tuilog.Success("Domain created")

	return nil
}
