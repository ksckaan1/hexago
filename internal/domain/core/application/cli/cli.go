package cli

import (
	"context"
	"fmt"

	"github.com/ksckaan1/hexago/internal/port"
)

type CLI struct {
	rootCmd             port.Commander
	initCmd             port.Commander
	domainCmd           port.Commander
	domainLSCmd         port.Commander
	domainCreateCmd     port.Commander
	serviceCmd          port.Commander
	serviceLSCmd        port.Commander
	serviceCreateCmd    port.Commander
	portCmd             port.Commander
	portLSCmd           port.Commander
	appCmd              port.Commander
	appLSCmd            port.Commander
	appCreateCmd        port.Commander
	entryPointCmd       port.Commander
	entryPointLSCmd     port.Commander
	entryPointCreateCmd port.Commander
	infraCmd            port.Commander
	infraLSCmd          port.Commander
	infraCreateCmd      port.Commander
	packageCmd          port.Commander
	packageLSCmd        port.Commander
	packageCreateCmd    port.Commander
	runnerCmd           port.Commander
	doctorCmd           port.Commander
	treeCmd             port.Commander
}

func New(
	rootCmd port.Commander,
	initCmd port.Commander,
	domainCmd port.Commander,
	domainLSCmd port.Commander,
	domainCreateCmd port.Commander,
	serviceCmd port.Commander,
	serviceLSCmd port.Commander,
	serviceCreateCmd port.Commander,
	portCmd port.Commander,
	portLSCmd port.Commander,
	appCmd port.Commander,
	appLSCmd port.Commander,
	appCreateCmd port.Commander,
	entryPointCmd port.Commander,
	entryPointLSCmd port.Commander,
	entryPointCreateCmd port.Commander,
	infraCmd port.Commander,
	infraLSCmd port.Commander,
	infraCreateCmd port.Commander,
	packageCmd port.Commander,
	packageLSCmd port.Commander,
	packageCreateCmd port.Commander,
	runnerCmd port.Commander,
	doctorCmd port.Commander,
	treeCmd port.Commander,
) (*CLI, error) {
	return &CLI{
		rootCmd:             rootCmd,
		initCmd:             initCmd,
		domainCmd:           domainCmd,
		domainLSCmd:         domainLSCmd,
		domainCreateCmd:     domainCreateCmd,
		serviceCmd:          serviceCmd,
		serviceLSCmd:        serviceLSCmd,
		serviceCreateCmd:    serviceCreateCmd,
		portCmd:             portCmd,
		portLSCmd:           portLSCmd,
		appCmd:              appCmd,
		appLSCmd:            appLSCmd,
		appCreateCmd:        appCreateCmd,
		entryPointCmd:       entryPointCmd,
		entryPointLSCmd:     entryPointLSCmd,
		entryPointCreateCmd: entryPointCreateCmd,
		infraCmd:            infraCmd,
		infraLSCmd:          infraLSCmd,
		infraCreateCmd:      infraCreateCmd,
		packageCmd:          packageCmd,
		packageLSCmd:        packageLSCmd,
		packageCreateCmd:    packageCreateCmd,
		runnerCmd:           runnerCmd,
		doctorCmd:           doctorCmd,
		treeCmd:             treeCmd,
	}, nil
}

func (c *CLI) Run(ctx context.Context) error {
	// init
	c.rootCmd.AddSubCommand(c.initCmd)

	// domain
	c.rootCmd.AddSubCommand(c.domainCmd)
	c.domainCmd.AddSubCommand(c.domainLSCmd)
	c.domainCmd.AddSubCommand(c.domainCreateCmd)

	// service
	c.rootCmd.AddSubCommand(c.serviceCmd)
	c.serviceCmd.AddSubCommand(c.serviceLSCmd)
	c.serviceCmd.AddSubCommand(c.serviceCreateCmd)

	// port
	c.rootCmd.AddSubCommand(c.portCmd)
	c.portCmd.AddSubCommand(c.portLSCmd)

	// app
	c.rootCmd.AddSubCommand(c.appCmd)
	c.appCmd.AddSubCommand(c.appLSCmd)
	c.appCmd.AddSubCommand(c.appCreateCmd)

	// cmd
	c.rootCmd.AddSubCommand(c.entryPointCmd)
	c.entryPointCmd.AddSubCommand(c.entryPointLSCmd)
	c.entryPointCmd.AddSubCommand(c.entryPointCreateCmd)

	// infra
	c.rootCmd.AddSubCommand(c.infraCmd)
	c.infraCmd.AddSubCommand(c.infraLSCmd)
	c.infraCmd.AddSubCommand(c.infraCreateCmd)

	// pkg
	c.rootCmd.AddSubCommand(c.packageCmd)
	c.packageCmd.AddSubCommand(c.packageLSCmd)
	c.packageCmd.AddSubCommand(c.packageCreateCmd)

	// run
	c.rootCmd.AddSubCommand(c.runnerCmd)

	// doctor
	c.rootCmd.AddSubCommand(c.doctorCmd)

	// tree
	c.rootCmd.AddSubCommand(c.treeCmd)

	err := c.rootCmd.Command().ExecuteContext(ctx)
	if err != nil {
		return fmt.Errorf("rootCmd.Command().ExecuteContext: %w", err)
	}
	return nil
}
