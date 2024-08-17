package servicecmd

import (
	"fmt"
	"slices"

	"github.com/charmbracelet/huh"
	"github.com/ksckaan1/hexago/internal/domain/core/port"
	"github.com/ksckaan1/hexago/internal/util"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type ServiceCreateCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
	// flags
	flagDomain  *string
	flagPkgName *string
}

func NewServiceCreateCommand(i *do.Injector) (*ServiceCreateCommand, error) {
	return &ServiceCreateCommand{
		cmd: &cobra.Command{
			Use:     "create",
			Example: "hexago service create <ServiceName>\nhexago s c <ServiceName>\nhexago service new <ServiceName>\nhexago s n <ServiceName>",
			Aliases: []string{"c", "new", "n"},
			Short:   "Create a service",
			Long:    `Create a service`,
			Args:    cobra.ExactArgs(1),
		},
		injector: i,
	}, nil
}

func (c *ServiceCreateCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *ServiceCreateCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *ServiceCreateCommand) init() {
	c.cmd.RunE = c.runner
	c.flagDomain = c.cmd.Flags().StringP("domain", "d", "", "hexago service new ServiceName -d core")
	c.flagPkgName = c.cmd.Flags().StringP("pkg-name", "p", "", "hexago service new ServiceName -p servicename")
}

func (c *ServiceCreateCommand) runner(cmd *cobra.Command, args []string) error {
	projectService, err := do.Invoke[port.ProjectService](c.injector)
	if err != nil {
		return fmt.Errorf("invoke project service: %w", err)
	}

	domains, err := projectService.GetAllDomains(cmd.Context())
	if err != nil {
		return fmt.Errorf("project service: get all domains: %w", err)
	}

	if len(domains) == 0 {
		return fmt.Errorf("No domains found.\nA domain needs to be created first")
	}

	if *c.flagDomain == "" {
		if len(domains) == 1 {
			*c.flagDomain = domains[0]
		} else {

			selectList := lo.Map(domains, func(d string, _ int) huh.Option[string] {
				return huh.NewOption(d, d)
			})

			err2 := huh.NewSelect[string]().
				Title("Select a domain.").
				Options(
					selectList...,
				).
				Value(c.flagDomain).Run()
			if err2 != nil {
				return fmt.Errorf("select a domain: %w", err2)
			}
		}
	} else if !slices.Contains(domains, *c.flagDomain) {
		return fmt.Errorf("domain not found: %s", *c.flagDomain)
	}

	err = projectService.CreateService(cmd.Context(), *c.flagDomain, args[0], *c.flagPkgName)
	if err != nil {
		return fmt.Errorf("project service: create service: %w", err)
	}

	fmt.Println("")
	util.UILog(util.Success, "service created")
	fmt.Println("")

	return nil
}
