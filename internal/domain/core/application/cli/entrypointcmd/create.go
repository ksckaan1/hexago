package entrypointcmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/internal/customerrors"
	"github.com/ksckaan1/hexago/internal/domain/core/model"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/ksckaan1/hexago/internal/port"
)

var _ port.Commander = (*EntryPointCreateCommand)(nil)

type EntryPointCreateCommand struct {
	cmd            *cobra.Command
	tuilog         *tuilog.TUILog
	projectService ProjectService
}

func NewEntryPointCreateCommand(projectService ProjectService, tl *tuilog.TUILog) (*EntryPointCreateCommand, error) {
	return &EntryPointCreateCommand{
		cmd: &cobra.Command{
			Use:     "new",
			Example: "hexago cmd new",
			Short:   "Create an entry point",
			Long:    `Create an entry point`,
		},
		tuilog:         tl,
		projectService: projectService,
	}, nil
}

func (c *EntryPointCreateCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *EntryPointCreateCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}

func (c *EntryPointCreateCommand) init() {
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := c.runner(cmd, args)
		if err != nil {
			return customerrors.ErrSuppressed
		}
		return nil
	}
}

func (c *EntryPointCreateCommand) runner(cmd *cobra.Command, args []string) error {
	var cmdName string

	if len(args) > 0 {
		cmdName = args[0]
	}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What’s entry point name?").
				Placeholder("entry-point-name").
				Validate(c.projectService.ValidateEntryPointName).
				Description("Entry point name must be kebab-case").
				Value(&cmdName),
		).WithShowHelp(true),
	).Run()
	if err != nil {
		return fmt.Errorf("input entry point name: %w", err)
	}

	epFile, err := c.projectService.CreateEntryPoint(
		cmd.Context(),
		model.CreateEntryPointParams{
			PackageName: cmdName,
		},
	)
	if err != nil {

		c.tuilog.Error(err.Error())

		return fmt.Errorf("projectService.CreateEntryPoint: %w", err)
	}

	c.tuilog.Success("Entry point created\n" + epFile)

	return nil
}
