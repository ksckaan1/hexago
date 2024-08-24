package entrypointcmd

import (
	"fmt"

	"github.com/ksckaan1/hexago/internal/domain/core/dto"
	"github.com/ksckaan1/hexago/internal/domain/core/port"
	"github.com/ksckaan1/hexago/internal/util"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type EntryPointCreateCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
}

func NewEntryPointCreateCommand(i *do.Injector) (*EntryPointCreateCommand, error) {
	return &EntryPointCreateCommand{
		cmd: &cobra.Command{
			Use:     "new",
			Example: "hexago cmd new <entrypoint-name>",
			Short:   "Create an entry point",
			Long:    `Create an entry point`,
			Args:    cobra.ExactArgs(1),
		},
		injector: i,
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

	epFile, err := projectService.CreateEntryPoint(
		cmd.Context(),
		dto.CreateEntryPointParams{
			PackageName: args[0],
		},
	)
	if err != nil {
		return fmt.Errorf("project service: create entry point: %w", err)
	}

	fmt.Println("")
	util.UILog(util.Success, "entry point created\n"+epFile)
	fmt.Println("")

	return nil
}
