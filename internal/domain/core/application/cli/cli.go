package cli

import (
	"context"
	"fmt"

	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/domaincmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/servicecmd"
	"github.com/samber/do"
	"github.com/spf13/cobra"
)

type Commander interface {
	Command() *cobra.Command
}

type CLI struct {
	rootCmd          *RootCommand
	initCmd          *InitCommand
	domainCmd        *domaincmd.DomainCommand
	domainLSCmd      *domaincmd.DomainLSCommand
	domainCreateCmd  *domaincmd.DomainCreateCommand
	serviceCmd       *servicecmd.ServiceCommand
	serviceLSCmd     *servicecmd.ServiceLSCommand
	serviceCreateCmd *servicecmd.ServiceCreateCommand
}

func New(i *do.Injector) (*CLI, error) {
	return &CLI{
		rootCmd:          do.MustInvoke[*RootCommand](i),
		initCmd:          do.MustInvoke[*InitCommand](i),
		domainCmd:        do.MustInvoke[*domaincmd.DomainCommand](i),
		domainLSCmd:      do.MustInvoke[*domaincmd.DomainLSCommand](i),
		domainCreateCmd:  do.MustInvoke[*domaincmd.DomainCreateCommand](i),
		serviceCmd:       do.MustInvoke[*servicecmd.ServiceCommand](i),
		serviceLSCmd:     do.MustInvoke[*servicecmd.ServiceLSCommand](i),
		serviceCreateCmd: do.MustInvoke[*servicecmd.ServiceCreateCommand](i),
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

	err := c.rootCmd.Execute(ctx)
	if err != nil {
		return fmt.Errorf("root cmd: execute: %w", err)
	}
	return nil
}
