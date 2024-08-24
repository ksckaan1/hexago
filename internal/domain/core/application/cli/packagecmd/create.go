package packagecmd

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

type PackageCreateCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
	// flags
	flagPkgName         *string
	flagPortName        *string
	flagNoPort          *bool
	flagAssertInterface *bool
	flagIsGlobal        *bool
}

func NewPackageCreateCommand(i *do.Injector) (*PackageCreateCommand, error) {
	return &PackageCreateCommand{
		cmd: &cobra.Command{
			Use:     "new",
			Example: "hexago pkg new <PackageName>",
			Short:   "Create a package",
			Long:    `Create a package`,
			Args:    cobra.ExactArgs(1),
		},
		injector: i,
	}, nil
}

func (c *PackageCreateCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *PackageCreateCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *PackageCreateCommand) init() {
	c.cmd.RunE = c.runner
	c.flagPkgName = c.cmd.Flags().StringP("pkg", "p", "", "hexago pkg new <PackageName> -p <packagename>")
	c.flagPortName = c.cmd.Flags().StringP("impl", "i", "", "hexago pkg new <PackageName> -i <domainname>:<PortName>")
	c.flagNoPort = c.cmd.Flags().BoolP("no-port", "n", false, "hexago pkg new <PackageName> -n")
	c.flagAssertInterface = c.cmd.Flags().BoolP("assert-port", "a", false, "hexago pkg new <PackageName> -i <domainname>:<PortName> -a")
	c.flagIsGlobal = c.cmd.Flags().BoolP("global", "g", false, "hexago pkg new <PackageName> -g")
}

func (c *PackageCreateCommand) runner(cmd *cobra.Command, args []string) error {
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

	packageFile, err := projectService.CreatePackage(
		cmd.Context(),
		dto.CreatePackageParams{
			StructName:      args[0],
			PackageName:     *c.flagPkgName,
			PortParam:       *c.flagPortName,
			AssertInterface: *c.flagAssertInterface,
			IsGlobal:        *c.flagIsGlobal,
		},
	)
	if err != nil {
		return fmt.Errorf("project service: create package: %w", err)
	}

	fmt.Println("")
	util.UILog(util.Success, "Package created\n"+packageFile)
	fmt.Println("")

	return nil
}
