package infracmd

import (
	"errors"
	"fmt"
	"github.com/ksckaan1/hexago/internal/port"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/ksckaan1/hexago/internal/domain/core/dto"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type InfraCreateCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
	tuilog   *tuilog.TUILog
}

func NewInfraCreateCommand(i *do.Injector) (*InfraCreateCommand, error) {
	return &InfraCreateCommand{
		cmd: &cobra.Command{
			Use:     "new",
			Example: "hexago infra new",
			Short:   "Create a infrastructure",
			Long:    `Create a infrastructure`,
		},
		injector: i,
		tuilog:   do.MustInvoke[*tuilog.TUILog](i),
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
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := c.runner(cmd, args)
		if err != nil {
			return dto.ErrSuppressed
		}
		return nil
	}
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
		fmt.Println("")
		c.tuilog.Error(err.Error())
		fmt.Println("")
		return fmt.Errorf("load config: %w", err)
	}

	var infraName string

	if len(args) > 0 {
		infraName = args[0]
	}

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What’s infrastructure name?").
				Placeholder("InfraName").
				Validate(projectService.ValidateInstanceName).
				Description("Infrastructure name must be PascalCase").
				Value(&infraName),
		).WithShowHelp(true),
	).Run()
	if err != nil {
		return fmt.Errorf("input infrastructure name: %w", err)
	}

	pkgName, err := c.selectPkgName(projectService, infraName)
	if err != nil {
		return fmt.Errorf("select pkg name: %w", err)
	}

	allPorts, err := projectService.GetAllPorts(cmd.Context())
	if err != nil {
		fmt.Println("")
		c.tuilog.Error(err.Error())
		fmt.Println("")
		return fmt.Errorf("get all ports: %w", err)
	}

	portInfo, err := c.selectPort(allPorts, infraName)
	if err != nil {
		return fmt.Errorf("select port: %w", err)
	}

	infraFile, err := projectService.CreateInfrastructure(
		cmd.Context(),
		dto.CreateInfraParams{
			StructName:      infraName,
			PackageName:     pkgName,
			PortParam:       portInfo.portName,
			AssertInterface: portInfo.assertInterface,
		},
	)
	if err != nil {
		fmt.Println("")
		if errors.Is(err, dto.ErrInvalidInstanceName) {
			c.tuilog.Error("Infrastructure name not valid\nMust be <PascalCase>")
		} else if errors.Is(err, dto.ErrInvalidPkgName) {
			c.tuilog.Error("Folder name not valid\nMust be <lowercase>")
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
		return fmt.Errorf("project service: create infrastructure: %w", err)
	}

	fmt.Println("")
	c.tuilog.Success("infrastructure created\n" + infraFile)
	fmt.Println("")

	return nil
}

func (c *InfraCreateCommand) selectPkgName(projectService port.ProjectService, instanceName string) (string, error) {
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

func (c *InfraCreateCommand) selectPort(allPorts []string, instanceName string) (*portInfo, error) {
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
							"var _ port.%s = (*%s)(nil)",
							portName,
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
