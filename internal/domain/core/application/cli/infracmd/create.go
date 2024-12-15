package infracmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/config"
	"github.com/ksckaan1/hexago/internal/customerrors"
	"github.com/ksckaan1/hexago/internal/domain/core/model"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/ksckaan1/hexago/internal/port"
)

var _ port.Commander = (*InfraCreateCommand)(nil)

type InfraCreateCommand struct {
	cmd            *cobra.Command
	tuilog         *tuilog.TUILog
	projectService ProjectService
	cfg            *config.Config
}

func NewInfraCreateCommand(projectService ProjectService, cfg *config.Config, tl *tuilog.TUILog) (*InfraCreateCommand, error) {
	return &InfraCreateCommand{
		cmd: &cobra.Command{
			Use:     "new",
			Example: "hexago infra new",
			Short:   "Create a infrastructure",
			Long:    `Create a infrastructure`,
		},
		tuilog:         tl,
		projectService: projectService,
		cfg:            cfg,
	}, nil
}

func (c *InfraCreateCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *InfraCreateCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}

func (c *InfraCreateCommand) init() {
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := c.runner(cmd, args)
		if err != nil {
			return customerrors.ErrSuppressed
		}
		return nil
	}
}

func (c *InfraCreateCommand) runner(cmd *cobra.Command, args []string) error {
	err := c.cfg.Load()
	if err != nil {

		c.tuilog.Error(err.Error())

		return fmt.Errorf("cfg.Load: %w", err)
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
				Validate(c.projectService.ValidateInstanceName).
				Description("Infrastructure name must be PascalCase").
				Value(&infraName),
		).WithShowHelp(true),
	).Run()
	if err != nil {
		return fmt.Errorf("input infrastructure name: %w", err)
	}

	pkgName, err := c.selectPkgName(infraName)
	if err != nil {
		return fmt.Errorf("select pkg name: %w", err)
	}

	allPorts, err := c.projectService.GetAllPorts(cmd.Context())
	if err != nil {

		c.tuilog.Error(err.Error())

		return fmt.Errorf("projectService.GetAllPorts: %w", err)
	}

	portInfo, err := c.selectPort(allPorts, infraName)
	if err != nil {
		return fmt.Errorf("select port: %w", err)
	}

	infraFile, err := c.projectService.CreateInfrastructure(
		cmd.Context(),
		model.CreateInfraParams{
			StructName:      infraName,
			PackageName:     pkgName,
			PortParam:       portInfo.portName,
			AssertInterface: portInfo.assertInterface,
		},
	)
	if err != nil {

		if errors.Is(err, customerrors.ErrInvalidInstanceName) {
			c.tuilog.Error("Infrastructure name not valid\nMust be <PascalCase>")
		} else if errors.Is(err, customerrors.ErrInvalidPkgName) {
			c.tuilog.Error("Folder name not valid\nMust be <lowercase>")
		} else if errors.Is(err, customerrors.ErrTemplateCanNotParsed) {
			c.tuilog.Error("Template can not parsed")
		} else if err2, ok1 := lo.ErrorsAs[customerrors.ErrTemplateCanNotExecute](err); ok1 {
			c.tuilog.Error("Template can not execute\n" + err2.Message)
		} else if err2, ok2 := lo.ErrorsAs[customerrors.ErrFormatGoFile](err); ok2 {
			c.tuilog.Error("Go file doesn't formatted\n" + err2.Message)
		} else {
			c.tuilog.Error(err.Error())
		}

		return fmt.Errorf("projectService.CreateInfrastructure: %w", err)
	}

	c.tuilog.Success("infrastructure created\n" + infraFile)

	return nil
}

func (c *InfraCreateCommand) selectPkgName(instanceName string) (string, error) {
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
					return c.projectService.ValidatePkgName(s)
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
