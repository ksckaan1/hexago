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
	flagDomain     *string
	flagPkgName    *string
	flagPortName   *string
	flagSelectPort *bool
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
	c.flagPortName = c.cmd.Flags().StringP("implement", "i", "", "hexago service new ServiceName -i core:Example")
	c.flagSelectPort = c.cmd.Flags().BoolP("select-port", "s", false, "hexago service new ServiceName -s")
}

func (c *ServiceCreateCommand) runner(cmd *cobra.Command, args []string) error {
	projectService, err := do.Invoke[port.ProjectService](c.injector)
	if err != nil {
		return fmt.Errorf("invoke project service: %w", err)
	}

	cfg, err := do.Invoke[port.ConfigService](c.injector)
	if err != nil {
		return fmt.Errorf("invoke config service: %w", err)
	}

	err = cfg.Load(".hexago/config.yaml")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
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
	} else if !slices.Contains(domains, *c.flagDomain) {
		return fmt.Errorf("domain not found: %s", *c.flagDomain)
	}

	if *c.flagSelectPort {
		allPorts := make([]string, 0)

		for i := range domains {
			ports, err := projectService.GetAllPorts(cmd.Context(), domains[i])
			if err != nil {
				return fmt.Errorf("get all ports: %w", err)
			}
			for j := range ports {
				allPorts = append(allPorts, domains[i]+":"+ports[j])
			}
		}

		selectPortList := []huh.Option[string]{
			huh.NewOption[string]("Do not implement!", ""),
		}

		selectPortList = append(selectPortList, lo.Map(allPorts, func(d string, _ int) huh.Option[string] {
			return huh.NewOption(d, d)
		})...)

		err2 := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select a port.").
					Options(
						selectPortList...,
					).
					Value(c.flagPortName),
			).WithShowHelp(true),
		).Run()
		if err != nil {
			return fmt.Errorf("select a port: %w", err2)
		}
	}

	serviceFile, err := projectService.CreateService(cmd.Context(), *c.flagDomain, args[0], *c.flagPkgName, *c.flagPortName)
	if err != nil {
		return fmt.Errorf("project service: create service: %w", err)
	}

	fmt.Println("")
	util.UILog(util.Success, "service created\n"+serviceFile)
	fmt.Println("")

	return nil
}
