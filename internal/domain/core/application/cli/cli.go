package cli

import (
	"context"
	"fmt"

	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/appcmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/domaincmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/entrypointcmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/infracmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/packagecmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/portcmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/servicecmd"
	"github.com/samber/do"
	"github.com/spf13/cobra"
)

type Commander interface {
	Command() *cobra.Command
}

type CLI struct {
	rootCmd             *RootCommand
	initCmd             *InitCommand
	domainCmd           *domaincmd.DomainCommand
	domainLSCmd         *domaincmd.DomainLSCommand
	domainCreateCmd     *domaincmd.DomainCreateCommand
	serviceCmd          *servicecmd.ServiceCommand
	serviceLSCmd        *servicecmd.ServiceLSCommand
	serviceCreateCmd    *servicecmd.ServiceCreateCommand
	portCmd             *portcmd.PortCommand
	portLSCmd           *portcmd.PortLSCommand
	appCmd              *appcmd.AppCommand
	appLSCmd            *appcmd.AppLSCommand
	appCreateCmd        *appcmd.AppCreateCommand
	entryPointCmd       *entrypointcmd.EntryPointCommand
	entryPointLSCmd     *entrypointcmd.EntryPointLSCommand
	entryPointCreateCmd *entrypointcmd.EntryPointCreateCommand
	infraCmd            *infracmd.InfraCommand
	infraLSCmd          *infracmd.InfraLSCommand
	infraCreateCmd      *infracmd.InfraCreateCommand
	packageCmd          *packagecmd.PackageCommand
	packageLSCmd        *packagecmd.PackageLSCommand
	packageCreateCmd    *packagecmd.PackageCreateCommand
}

func New(i *do.Injector) (*CLI, error) {
	return &CLI{
		rootCmd:             do.MustInvoke[*RootCommand](i),
		initCmd:             do.MustInvoke[*InitCommand](i),
		domainCmd:           do.MustInvoke[*domaincmd.DomainCommand](i),
		domainLSCmd:         do.MustInvoke[*domaincmd.DomainLSCommand](i),
		domainCreateCmd:     do.MustInvoke[*domaincmd.DomainCreateCommand](i),
		serviceCmd:          do.MustInvoke[*servicecmd.ServiceCommand](i),
		serviceLSCmd:        do.MustInvoke[*servicecmd.ServiceLSCommand](i),
		serviceCreateCmd:    do.MustInvoke[*servicecmd.ServiceCreateCommand](i),
		portCmd:             do.MustInvoke[*portcmd.PortCommand](i),
		portLSCmd:           do.MustInvoke[*portcmd.PortLSCommand](i),
		appCmd:              do.MustInvoke[*appcmd.AppCommand](i),
		appLSCmd:            do.MustInvoke[*appcmd.AppLSCommand](i),
		appCreateCmd:        do.MustInvoke[*appcmd.AppCreateCommand](i),
		entryPointCmd:       do.MustInvoke[*entrypointcmd.EntryPointCommand](i),
		entryPointLSCmd:     do.MustInvoke[*entrypointcmd.EntryPointLSCommand](i),
		entryPointCreateCmd: do.MustInvoke[*entrypointcmd.EntryPointCreateCommand](i),
		infraCmd:            do.MustInvoke[*infracmd.InfraCommand](i),
		infraLSCmd:          do.MustInvoke[*infracmd.InfraLSCommand](i),
		infraCreateCmd:      do.MustInvoke[*infracmd.InfraCreateCommand](i),
		packageCmd:          do.MustInvoke[*packagecmd.PackageCommand](i),
		packageLSCmd:        do.MustInvoke[*packagecmd.PackageLSCommand](i),
		packageCreateCmd:    do.MustInvoke[*packagecmd.PackageCreateCommand](i),
	}, nil
}

func (c *CLI) Run(ctx context.Context) error {
	c.rootCmd.AddCommand(c.initCmd)
	c.rootCmd.AddCommand(c.domainCmd)
	c.domainCmd.AddCommand(c.domainLSCmd)
	c.domainCmd.AddCommand(c.domainCreateCmd)
	c.rootCmd.AddCommand(c.serviceCmd)
	c.serviceCmd.AddCommand(c.serviceLSCmd)
	c.serviceCmd.AddCommand(c.serviceCreateCmd)
	c.rootCmd.AddCommand(c.portCmd)
	c.portCmd.AddCommand(c.portLSCmd)
	c.rootCmd.AddCommand(c.appCmd)
	c.appCmd.AddCommand(c.appLSCmd)
	c.appCmd.AddCommand(c.appCreateCmd)
	c.rootCmd.AddCommand(c.entryPointCmd)
	c.entryPointCmd.AddCommand(c.entryPointLSCmd)
	c.entryPointCmd.AddCommand(c.entryPointCreateCmd)
	c.rootCmd.AddCommand(c.infraCmd)
	c.infraCmd.AddCommand(c.infraLSCmd)
	c.infraCmd.AddCommand(c.infraCreateCmd)
	c.rootCmd.AddCommand(c.packageCmd)
	c.packageCmd.AddCommand(c.packageLSCmd)
	c.packageCmd.AddCommand(c.packageCreateCmd)

	err := c.rootCmd.Execute(ctx)
	if err != nil {
		return fmt.Errorf("root cmd: execute: %w", err)
	}
	return nil
}
