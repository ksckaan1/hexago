package servicecmd

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/internal/customerrors"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/ksckaan1/hexago/internal/port"
)

var _ port.Commander = (*ServiceLSCommand)(nil)

type ServiceLSCommand struct {
	cmd            *cobra.Command
	tuilog         *tuilog.TUILog
	projectService ProjectService

	// flags
	flagLine   *bool
	flagDomain *string
}

func NewServiceLSCommand(projectService ProjectService, tl *tuilog.TUILog) (*ServiceLSCommand, error) {
	return &ServiceLSCommand{
		cmd: &cobra.Command{
			Use:     "ls",
			Example: "hexago service ls -d <domainname>\nhexago service ls (select domain interatively)",
			Short:   "List services",
			Long:    `List services`,
		},
		projectService: projectService,
		tuilog:         tl,
	}, nil
}

func (c *ServiceLSCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *ServiceLSCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}

func (c *ServiceLSCommand) init() {
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := c.runner(cmd, args)
		if err != nil {
			return customerrors.ErrSuppressed
		}
		return nil
	}
	c.flagLine = c.cmd.Flags().BoolP("line", "l", false, "hexago service ls -l")
	c.flagDomain = c.cmd.Flags().StringP("domain", "d", "", "hexago service ls -d <domainname>")
}

func (c *ServiceLSCommand) runner(cmd *cobra.Command, _ []string) error {
	domains, err := c.projectService.GetAllDomains(cmd.Context())
	if err != nil {
		return fmt.Errorf("projectService.GetAllDomains: %w", err)
	}

	if len(domains) == 0 {
		c.tuilog.Error("No domains found.\nA domain needs to be created first")
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

		c.tuilog.Error("Domain not found: ", *c.flagDomain)

		return fmt.Errorf("domain not found: %s", *c.flagDomain)
	}

	allServices := make([]string, 0)

	for i := range domains {
		if !(*c.flagDomain == "*" || domains[i] == *c.flagDomain) {
			continue
		}

		services, err2 := c.projectService.GetAllServices(cmd.Context(), domains[i])
		if err2 != nil {

			if errors.Is(err2, customerrors.ErrDomainNotFound) {
				c.tuilog.Error("Domain not found: " + domains[i])
			} else {
				c.tuilog.Error(err2.Error())
			}

			return fmt.Errorf("projectService.GetAllServices: %w", err2)
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
