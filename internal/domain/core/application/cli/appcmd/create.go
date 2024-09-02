package appcmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/ksckaan1/hexago/internal/domain/core/dto"
	"github.com/ksckaan1/hexago/internal/domain/core/port"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type AppCreateCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
	tuilog   *tuilog.TUILog
}

func NewAppCreateCommand(i *do.Injector) (*AppCreateCommand, error) {
	return &AppCreateCommand{
		cmd: &cobra.Command{
			Use:     "new",
			Example: "hexago app new",
			Short:   "Create an application",
			Long:    `Create an application`,
		},
		injector: i,
		tuilog:   do.MustInvoke[*tuilog.TUILog](i),
	}, nil
}

func (c *AppCreateCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *AppCreateCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *AppCreateCommand) init() {
	c.cmd.RunE = c.runner
}

func (c *AppCreateCommand) runner(cmd *cobra.Command, args []string) error {
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
		fmt.Println("")
		c.tuilog.Error(err.Error())
		fmt.Println("")
		return fmt.Errorf("load config: %w", err)
	}

	domains, err := projectService.GetAllDomains(cmd.Context())
	if err != nil {
		fmt.Println("")
		c.tuilog.Error(err.Error())
		fmt.Println("")
		return fmt.Errorf("project service: get all domains: %w", err)
	}

	if len(domains) == 0 {
		fmt.Println("")
		c.tuilog.Error("No domains found.\nA domain needs to be created first")
		fmt.Println("")
		return fmt.Errorf("No domains found.\nA domain needs to be created first")
	}

	var appName string

	if len(args) > 0 {
		appName = args[0]
	}

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What’s application name?").
				Placeholder("AppName").
				Validate(projectService.ValidateInstanceName).
				Description("Application name must be PascalCase").
				Value(&appName),
		).WithShowHelp(true),
	).Run()
	if err != nil {
		return fmt.Errorf("input application name: %w", err)
	}

	pkgName, err := c.selectPkgName(projectService, appName)
	if err != nil {
		return fmt.Errorf("select pkg name: %w", err)
	}

	var domainName string
	if len(domains) == 1 {
		domainName = domains[0]
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
					Value(&domainName),
			).WithShowHelp(true),
		).Run()
		if err2 != nil {
			fmt.Println("")
			c.tuilog.Error("Select a domain: " + err2.Error())
			fmt.Println("")
			return fmt.Errorf("select a domain: %w", err2)
		}
	}

	allPorts := make([]string, 0)

	for i := range domains {
		ports, err := projectService.GetAllPorts(cmd.Context(), domains[i])
		if err != nil {
			fmt.Println("")
			c.tuilog.Error(err.Error())
			fmt.Println("")
			return fmt.Errorf("get all ports: %w", err)
		}
		for j := range ports {
			allPorts = append(allPorts, domains[i]+":"+ports[j])
		}
	}

	portInfo, err := c.selectPort(allPorts, appName)
	if err != nil {
		return fmt.Errorf("select port: %w", err)
	}

	applicationFile, err := projectService.CreateApplication(
		cmd.Context(),
		dto.CreateApplicationParams{
			TargetDomain:    domainName,
			StructName:      appName,
			PackageName:     pkgName,
			PortParam:       portInfo.portName,
			AssertInterface: portInfo.assertInterface,
		},
	)
	if err != nil {
		fmt.Println("")
		if errors.Is(err, dto.ErrInvalidInstanceName) {
			c.tuilog.Error("Application name not valid\nMust be <PascalCase>")
		} else if errors.Is(err, dto.ErrInvalidPkgName) {
			c.tuilog.Error("Folder name not valid\nMust be <lowercase>")
		} else if errors.Is(err, dto.ErrDomainNotFound) {
			c.tuilog.Error("Domain not found")
		} else if errors.Is(err, dto.ErrTemplateCanNotParsed) {
			c.tuilog.Error("Template can not parsed")
		} else if err2, ok := lo.ErrorsAs[dto.ErrTemplateCanNotExecute](err); ok {
			c.tuilog.Error("Template can not execute\n" + err2.Message)
		} else if err2, ok := lo.ErrorsAs[dto.ErrFormatGoFile](err); ok {
			c.tuilog.Error("Go file doesn't formatted\n" + err2.Message)
		} else {
			c.tuilog.Error(err.Error())
		}
		fmt.Println("")
		return fmt.Errorf("project service: create application: %w", err)
	}

	fmt.Println("")
	c.tuilog.Success("Application created\n" + applicationFile)
	fmt.Println("")

	return nil
}

func (c *AppCreateCommand) selectPkgName(projectService port.ProjectService, instanceName string) (string, error) {
	var pkgName string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What’s folder (pkg) name?").
				Placeholder(strings.ToLower(instanceName)).
				Validate(func(s string) error {
					if s == "" {
						return nil
					}
					return projectService.ValidatePkgName(s)
				}).
				Description("Folder name must be lowercase").
				Value(&pkgName),
		).WithShowHelp(true),
	).Run()
	if err != nil {
		return "", fmt.Errorf("input pkg name: %w", err)
	}
	return pkgName, nil
}

type portInfo struct {
	portName        string
	assertInterface bool
}

func (c *AppCreateCommand) selectPort(allPorts []string, instanceName string) (*portInfo, error) {
	if len(allPorts) == 0 {
		return &portInfo{}, nil
	}

	selectPortList := []huh.Option[string]{
		huh.NewOption[string]("Do not implement!", ""),
	}

	selectPortList = append(selectPortList, lo.Map(allPorts, func(d string, _ int) huh.Option[string] {
		return huh.NewOption(d, d)
	})...)

	var portName string

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a port.").
				Options(
					selectPortList...,
				).
				Value(&portName),
		).WithShowHelp(true),
	).Run()
	if err != nil {
		return nil, fmt.Errorf("select a port: %w", err)
	}

	assertInterface := false

	if portName != "" {
		err = huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Do you want to assert port?").
					Description(
						fmt.Sprintf(
							"var _ %sport.%s = (*%s)(nil)",
							strings.Split(portName, ":")[0],
							strings.Split(portName, ":")[1],
							instanceName,
						),
					).
					Affirmative("Yes").
					Negative("No").
					Value(&assertInterface),
			).WithShowHelp(true),
		).Run()
		if err != nil {
			return nil, fmt.Errorf("confirm assert port: %w", err)
		}
	}

	return &portInfo{
		portName:        portName,
		assertInterface: assertInterface,
	}, nil
}
