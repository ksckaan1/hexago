package servicecmd

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/ksckaan1/hexago/internal/domain/core/port"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type ServiceLSCommand struct {
	cmd      *cobra.Command
	injector *do.Injector

	// flags
	flagLine   *bool
	flagDomain *string
}

func NewServiceLSCommand(i *do.Injector) (*ServiceLSCommand, error) {
	return &ServiceLSCommand{
		cmd: &cobra.Command{
			Use:     "ls",
			Example: "hexago service ls -d <domainname>\nhexago service ls (select domain interatively)",
			Short:   "List services",
			Long:    `List services`,
		},
		injector: i,
		// flags
		flagLine: new(bool),
	}, nil
}

func (c *ServiceLSCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *ServiceLSCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *ServiceLSCommand) init() {
	c.cmd.RunE = c.runner
	c.flagLine = c.cmd.Flags().BoolP("line", "l", false, "hexago service ls -l")
	c.flagDomain = c.cmd.Flags().StringP("domain", "d", "", "hexago service ls -d <domainname>")
}

func (c *ServiceLSCommand) runner(cmd *cobra.Command, _ []string) error {
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

			selectList := []huh.Option[string]{
				huh.NewOption("* (All Domains)", "*"),
			}

			selectList = append(selectList, lo.Map(domains, func(d string, _ int) huh.Option[string] {
				return huh.NewOption(d, d)
			})...)

			err2 := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Select a domain.").
						Options(
							selectList...,
						).
						Value(c.flagDomain),
				).WithShowHelp(true),
			).Run()
			if err2 != nil {
				return fmt.Errorf("select a domain: %w", err2)
			}
		}
	} else if *c.flagDomain != "*" && !slices.Contains(domains, *c.flagDomain) {
		return fmt.Errorf("domain not found: %s", *c.flagDomain)
	}

	allServices := make([]string, 0)

	for i := range domains {
		if !(*c.flagDomain == "*" || domains[i] == *c.flagDomain) {
			continue
		}

		services, err2 := projectService.GetAllServices(cmd.Context(), domains[i])
		if err2 != nil {
			return fmt.Errorf("project service: get all services: %w", err2)
		}

		if *c.flagDomain == "*" {
			for j := range services {
				services[j] = domains[i] + ":" + services[j]
			}
		}

		allServices = append(allServices, services...)
	}

	seperator := lo.Ternary(*c.flagLine, "\n", " ")

	fmt.Println(strings.Join(allServices, seperator))

	return nil
}
