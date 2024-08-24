package infracmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/ksckaan1/hexago/internal/domain/core/dto"
	"github.com/ksckaan1/hexago/internal/domain/core/port"
	"github.com/ksckaan1/hexago/internal/util"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type InfraCreateCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
	// flags
	flagPkgName         *string
	flagPortName        *string
	flagNoPort          *bool
	flagAssertInterface *bool
}

func NewInfraCreateCommand(i *do.Injector) (*InfraCreateCommand, error) {
	return &InfraCreateCommand{
		cmd: &cobra.Command{
			Use:     "new",
			Example: "hexago infra new <InfraName>",
			Short:   "Create a infrastructure",
			Long:    `Create a infrastructure`,
			Args:    cobra.ExactArgs(1),
		},
		injector: i,
	}, nil
}

func (c *InfraCreateCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *InfraCreateCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *InfraCreateCommand) init() {
	c.cmd.RunE = c.runner
	c.flagPkgName = c.cmd.Flags().StringP("pkg", "p", "", "hexago infra new <InfraName> -p <infraname>")
	c.flagPortName = c.cmd.Flags().StringP("impl", "i", "", "hexago infra new <InfraName> -i <domainname>:<PortName>")
	c.flagNoPort = c.cmd.Flags().BoolP("no-port", "n", false, "hexago infra new <InfraName> -n")
	c.flagAssertInterface = c.cmd.Flags().BoolP("assert-port", "a", false, "hexago infra new <InfraName> -i <domainname>:<PortName> -a")
}

func (c *InfraCreateCommand) runner(cmd *cobra.Command, args []string) error {
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

	if *c.flagNoPort {
		*c.flagPortName = ""
	}

	if !*c.flagNoPort && *c.flagPortName == "" {
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
		if err2 != nil {
			return fmt.Errorf("select a port: %w", err2)
		}
	}

	infraFile, err := projectService.CreateInfrastructure(
		cmd.Context(),
		dto.CreateInfraParams{
			StructName:      args[0],
			PackageName:     *c.flagPkgName,
			PortParam:       *c.flagPortName,
			AssertInterface: *c.flagAssertInterface,
		},
	)
	if err != nil {
		return fmt.Errorf("project service: create infrastructure: %w", err)
	}

	fmt.Println("")
	util.UILog(util.Success, "infrastructure created\n"+infraFile)
	fmt.Println("")

	return nil
}
