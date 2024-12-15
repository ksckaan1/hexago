package packagecmd

import (
	"errors"
	"fmt"
	"path/filepath"
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

var _ port.Commander = (*PackageCreateCommand)(nil)

type PackageCreateCommand struct {
	cmd            *cobra.Command
	tuilog         *tuilog.TUILog
	projectService ProjectService
	cfg            *config.Config
}

func NewPackageCreateCommand(projectService ProjectService, cfg *config.Config, tl *tuilog.TUILog) (*PackageCreateCommand, error) {
	return &PackageCreateCommand{
		cmd: &cobra.Command{
			Use:     "new",
			Example: "hexago pkg new",
			Short:   "Create a package",
			Long:    `Create a package`,
		},
		projectService: projectService,
		tuilog:         tl,
		cfg:            cfg,
	}, nil
}

func (c *PackageCreateCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *PackageCreateCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}

func (c *PackageCreateCommand) init() {
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := c.runner(cmd, args)
		if err != nil {
			return customerrors.ErrSuppressed
		}
		return nil
	}
}

func (c *PackageCreateCommand) runner(cmd *cobra.Command, args []string) error {
	err := c.cfg.Load()
	if err != nil {

		c.tuilog.Error(err.Error())

		return fmt.Errorf("cfg.Load: %w", err)
	}

	var packageName string

	if len(args) > 0 {
		packageName = args[0]
	}

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What’s package name?").
				Placeholder("PackageName").
				Validate(c.projectService.ValidateInstanceName).
				Description("Package name must be PascalCase").
				Value(&packageName),
		).WithShowHelp(true),
	).Run()
	if err != nil {
		return fmt.Errorf("input package name: %w", err)
	}

	pkgName, err := c.selectPkgName(packageName)
	if err != nil {
		return fmt.Errorf("select pkg name: %w", err)
	}

	allPorts, err := c.projectService.GetAllPorts(cmd.Context())
	if err != nil {

		c.tuilog.Error(err.Error())

		return fmt.Errorf("projectService.GetAllPorts: %w", err)
	}

	portInfo, err := c.selectPort(allPorts, packageName)
	if err != nil {
		return fmt.Errorf("select port: %w", err)
	}

	var isGlobal bool

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[bool]().
				Title("Select package scope").
				Options(
					huh.NewOption(
						fmt.Sprintf("internal (%q)", filepath.Join("internal", "pkg", "*")),
						false),
					huh.NewOption(
						fmt.Sprintf("global (%q)", filepath.Join("pkg", "*")),
						true),
				).
				Value(&isGlobal),
		).WithShowHelp(true),
	).Run()
	if err != nil {

		c.tuilog.Error("Select a port: ", err.Error())

		return fmt.Errorf("select is global: %w", err)
	}

	packageFile, err := c.projectService.CreatePackage(
		cmd.Context(),
		model.CreatePackageParams{
			StructName:      packageName,
			PackageName:     pkgName,
			PortParam:       portInfo.portName,
			AssertInterface: portInfo.assertInterface,
			IsGlobal:        isGlobal,
		},
	)
	if err != nil {

		if errors.Is(err, customerrors.ErrInvalidInstanceName) {
			c.tuilog.Error("Package name not valid\nMust be <PascalCase>")
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

		return fmt.Errorf("projectService.CreatePackage: %w", err)
	}

	c.tuilog.Success("Package created\n" + packageFile)

	return nil
}

func (c *PackageCreateCommand) selectPkgName(instanceName string) (string, error) {
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

func (c *PackageCreateCommand) selectPort(allPorts []string, instanceName string) (*portInfo, error) {
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
