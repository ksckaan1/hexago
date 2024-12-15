package main

import (
	"fmt"

	"github.com/ksckaan1/hexago/config"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/appcmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/doctorcmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/domaincmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/entrypointcmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/infracmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/initcmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/packagecmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/portcmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/rootcmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/runnercmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/servicecmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/treecmd"
	"github.com/ksckaan1/hexago/internal/domain/core/service/project"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
)

func initCommands(projectService *project.Project, tl *tuilog.TUILog, cfg *config.Config) (*cli.CLI, error) {
	// root
	rootCmd, err := rootcmd.NewRootCommand()
	if err != nil {
		return nil, fmt.Errorf("rootcmd.NewRootCommand: %w", err)
	}

	// init
	initCmd, err := initcmd.NewInitCommand(projectService, tl)
	if err != nil {
		return nil, fmt.Errorf("initcmd.NewInitCommand: %w", err)
	}

	// domain
	domainCmd, err := domaincmd.NewDomainCommand()
	if err != nil {
		return nil, fmt.Errorf("domaincmd.NewDomainCommand: %w", err)
	}

	domainLSCmd, err := domaincmd.NewDomainLSCommand(projectService, tl)
	if err != nil {
		return nil, fmt.Errorf("domaincmd.NewDomainLSCommand: %w", err)
	}

	domainCreateCmd, err := domaincmd.NewDomainCreateCommand(projectService, tl)
	if err != nil {
		return nil, fmt.Errorf("domaincmd.NewDomainCreateCommand: %w", err)
	}

	// service
	serviceCmd, err := servicecmd.NewServiceCommand()
	if err != nil {
		return nil, fmt.Errorf("servicecmd.NewServiceCommand: %w", err)
	}

	serviceLSCmd, err := servicecmd.NewServiceLSCommand(projectService, tl)
	if err != nil {
		return nil, fmt.Errorf("servicecmd.NewServiceLSCommand: %w", err)
	}

	serviceCreateCmd, err := servicecmd.NewServiceCreateCommand(projectService, cfg, tl)
	if err != nil {
		return nil, fmt.Errorf("servicecmd.NewServiceCreateCommand: %w", err)
	}

	// port
	portCmd, err := portcmd.NewPortCommand()
	if err != nil {
		return nil, fmt.Errorf("portcmd.NewPortCommand: %w", err)
	}

	portLSCmd, err := portcmd.NewPortLSCommand(projectService, tl)
	if err != nil {
		return nil, fmt.Errorf("portcmd.NewPortLSCommand: %w", err)
	}

	// app
	appCmd, err := appcmd.NewAppCommand()
	if err != nil {
		return nil, fmt.Errorf("appcmd.NewAppCommand: %w", err)
	}

	appLSCmd, err := appcmd.NewAppLSCommand(projectService, tl)
	if err != nil {
		return nil, fmt.Errorf("appcmd.NewAppLSCommand: %w", err)
	}

	appCreateCmd, err := appcmd.NewAppCreateCommand(projectService, cfg, tl)
	if err != nil {
		return nil, fmt.Errorf("appcmd.NewAppCreateCommand: %w", err)
	}

	// cmd
	entryPointCmd, err := entrypointcmd.NewEntryPointCommand()
	if err != nil {
		return nil, fmt.Errorf("entrypointcmd.NewEntryPointCommand: %w", err)
	}

	entryPointLSCmd, err := entrypointcmd.NewEntryPointLSCommand(projectService, tl)
	if err != nil {
		return nil, fmt.Errorf("entrypointcmd.NewEntryPointLSCommand: %w", err)
	}

	entryPointCreateCmd, err := entrypointcmd.NewEntryPointCreateCommand(projectService, tl)
	if err != nil {
		return nil, fmt.Errorf("entrypointcmd.NewEntryPointCreateCommand: %w", err)
	}

	// infra
	infraCmd, err := infracmd.NewInfraCommand()
	if err != nil {
		return nil, fmt.Errorf("infracmd.NewInfraCommand: %w", err)
	}

	infraLSCmd, err := infracmd.NewInfraLSCommand(projectService, tl)
	if err != nil {
		return nil, fmt.Errorf("infracmd.NewInfraLSCommand: %w", err)
	}

	infraCreateCmd, err := infracmd.NewInfraCreateCommand(projectService, cfg, tl)
	if err != nil {
		return nil, fmt.Errorf("infracmd.NewInfraCreateCommand: %w", err)
	}

	// pkg
	packageCmd, err := packagecmd.NewPackageCommand()
	if err != nil {
		return nil, fmt.Errorf("packagecmd.NewPackageCommand: %w", err)
	}

	packageLSCmd, err := packagecmd.NewPackageLSCommand(projectService, tl)
	if err != nil {
		return nil, fmt.Errorf("packagecmd.NewPackageLSCommand: %w", err)
	}

	packageCreateCmd, err := packagecmd.NewPackageCreateCommand(projectService, cfg, tl)
	if err != nil {
		return nil, fmt.Errorf("packagecmd.NewPackageCreateCommand: %w", err)
	}

	// run
	runnerCmd, err := runnercmd.NewRunnerCommand(projectService, cfg, tl)
	if err != nil {
		return nil, fmt.Errorf("runnercmd.NewRunnerCommand: %w", err)
	}

	// doctor
	doctorCmd, err := doctorcmd.NewDoctorCommand(projectService, tl)
	if err != nil {
		return nil, fmt.Errorf("doctorcmd.NewDoctorCommand: %w", err)
	}

	// tree
	treeCmd, err := treecmd.NewTreeCommand(projectService, tl)
	if err != nil {
		return nil, fmt.Errorf("treecmd.NewTreeCommand: %w", err)
	}

	app, err := cli.New(
		rootCmd,
		initCmd,
		domainCmd,
		domainLSCmd,
		domainCreateCmd,
		serviceCmd,
		serviceLSCmd,
		serviceCreateCmd,
		portCmd,
		portLSCmd,
		appCmd,
		appLSCmd,
		appCreateCmd,
		entryPointCmd,
		entryPointLSCmd,
		entryPointCreateCmd,
		infraCmd,
		infraLSCmd,
		infraCreateCmd,
		packageCmd,
		packageLSCmd,
		packageCreateCmd,
		runnerCmd,
		doctorCmd,
		treeCmd,
	)
	if err != nil {
		return nil, fmt.Errorf("cli.New: %w", err)
	}

	return app, nil
}
