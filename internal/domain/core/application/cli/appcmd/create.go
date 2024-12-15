package appcmd

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

var _ port.Commander = (*AppCreateCommand)(nil)

type AppCreateCommand struct {
	cmd            *cobra.Command
	tuilog         *tuilog.TUILog
	projectService ProjectService
	cfg            *config.Config
}

func NewAppCreateCommand(projectService ProjectService, cfg *config.Config, tl *tuilog.TUILog) (*AppCreateCommand, error) {
	return &AppCreateCommand{
		cmd: &cobra.Command{
			Use:     "new",
			Example: "hexago app new",
			Short:   "Create an application",
			Long:    `Create an application`,
		},
		projectService: projectService,
		tuilog:         tl,
		cfg:            cfg,
	}, nil
}

func (c *AppCreateCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *AppCreateCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}

func (c *AppCreateCommand) init() {
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := c.runner(cmd, args)
		if err != nil {
			return customerrors.ErrSuppressed
		}
		return nil
	}
}

func (c *AppCreateCommand) runner(cmd *cobra.Command, args []string) error {
	err := c.cfg.Load()
	if err != nil {
		c.tuilog.Error(err.Error())
		return fmt.Errorf("cfg.Load: %w", err)
	}

	domains, err := c.projectService.GetAllDomains(cmd.Context())
	if err != nil {
		c.tuilog.Error(err.Error())
		return fmt.Errorf("projectService.GetAllDomains: %w", err)
	}

	if len(domains) == 0 {
		c.tuilog.Error("No domains found.\nA domain needs to be created first")
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
				Validate(c.projectService.ValidateInstanceName).
				Description("Application name must be PascalCase").
				Value(&appName),
		).WithShowHelp(true),
	).Run()
	if err != nil {
		return fmt.Errorf("input application name: %w", err)
	}

	pkgName, err := c.selectPkgName(appName)
	if err != nil {
		return fmt.Errorf("selectPkgName: %w", err)
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
			c.tuilog.Error("Select a domain: " + err2.Error())
			return fmt.Errorf("select a domain: %w", err2)
		}
	}

	allPorts, err := c.projectService.GetAllPorts(cmd.Context())
	if err != nil {
		c.tuilog.Error(err.Error())
		return fmt.Errorf("get all ports: %w", err)
	}

	portInfo, err := c.selectPort(allPorts, appName)
	if err != nil {
		return fmt.Errorf("select port: %w", err)
	}

	applicationFile, err := c.projectService.CreateApplication(
		cmd.Context(),
		model.CreateApplicationParams{
			TargetDomain:    domainName,
			StructName:      appName,
			PackageName:     pkgName,
			PortParam:       portInfo.portName,
			AssertInterface: portInfo.assertInterface,
		},
	)
	if err != nil {
		if errors.Is(err, customerrors.ErrInvalidInstanceName) {
			c.tuilog.Error("Application name not valid\nMust be <PascalCase>")
		} else if errors.Is(err, customerrors.ErrInvalidPkgName) {
			c.tuilog.Error("Folder name not valid\nMust be <lowercase>")
		} else if errors.Is(err, customerrors.ErrDomainNotFound) {
			c.tuilog.Error("Domain not found")
		} else if errors.Is(err, customerrors.ErrTemplateCanNotParsed) {
			c.tuilog.Error("Template can not parsed")
		} else if err2, ok1 := lo.ErrorsAs[customerrors.ErrTemplateCanNotExecute](err); ok1 {
			c.tuilog.Error("Template can not execute\n" + err2.Message)
		} else if err2, ok2 := lo.ErrorsAs[customerrors.ErrFormatGoFile](err); ok2 {
			c.tuilog.Error("Go file doesn't formatted\n" + err2.Message)
		} else {
			c.tuilog.Error(err.Error())
		}
		return fmt.Errorf("projectService.CreateApplication: %w", err)
	}

	c.tuilog.Success("Application created\n" + applicationFile)

	return nil
}

func (c *AppCreateCommand) selectPkgName(instanceName string) (string, error) {
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
