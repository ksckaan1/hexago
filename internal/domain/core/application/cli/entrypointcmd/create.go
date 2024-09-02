package entrypointcmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/ksckaan1/hexago/internal/domain/core/dto"
	"github.com/ksckaan1/hexago/internal/domain/core/port"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type EntryPointCreateCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
	tuilog   *tuilog.TUILog
}

func NewEntryPointCreateCommand(i *do.Injector) (*EntryPointCreateCommand, error) {
	return &EntryPointCreateCommand{
		cmd: &cobra.Command{
			Use:     "new",
			Example: "hexago cmd new",
			Short:   "Create an entry point",
			Long:    `Create an entry point`,
		},
		injector: i,
		tuilog:   do.MustInvoke[*tuilog.TUILog](i),
	}, nil
}

func (c *EntryPointCreateCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *EntryPointCreateCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *EntryPointCreateCommand) init() {
	c.cmd.RunE = c.runner
}

func (c *EntryPointCreateCommand) runner(cmd *cobra.Command, args []string) error {
	projectService, err := do.Invoke[port.ProjectService](c.injector)
	if err != nil {
		return fmt.Errorf("invoke project service: %w", err)
	}

	var cmdName string
	if len(args) > 0 {
		cmdName = args[0]
	}

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What’s entry point name?").
				Placeholder("entry-point-name").
				Validate(projectService.ValidateEntryPointName).
				Description("Entry point name must be kebab-case").
				Value(&cmdName),
		).WithShowHelp(true),
	).Run()
	if err != nil {
		return fmt.Errorf("input entry point name: %w", err)
	}

	epFile, err := projectService.CreateEntryPoint(
		cmd.Context(),
		dto.CreateEntryPointParams{
			PackageName: cmdName,
		},
	)
	if err != nil {
		fmt.Println("")
		c.tuilog.Error(err.Error())
		fmt.Println("")
		return fmt.Errorf("project service: create entry point: %w", err)
	}

	fmt.Println("")
	c.tuilog.Success("Entry point created\n" + epFile)
	fmt.Println("")

	return nil
}
