package entrypointcmd

import (
	"fmt"
	"strings"

	"github.com/ksckaan1/hexago/internal/domain/core/port"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type EntryPointLSCommand struct {
	cmd      *cobra.Command
	injector *do.Injector

	// flags
	flagLine *bool
}

func NewEntryPointLSCommand(i *do.Injector) (*EntryPointLSCommand, error) {
	return &EntryPointLSCommand{
		cmd: &cobra.Command{
			Use:     "ls",
			Example: "hexago cmd ls",
			Short:   "List Entry Points",
			Long:    `List Entry Points`,
		},
		injector: i,
		// flags
		flagLine: new(bool),
	}, nil
}

func (c *EntryPointLSCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *EntryPointLSCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *EntryPointLSCommand) init() {
	c.cmd.RunE = c.runner
	c.flagLine = c.cmd.Flags().BoolP("line", "l", false, "hexago service ls -l")
}

func (c *EntryPointLSCommand) runner(cmd *cobra.Command, _ []string) error {
	projectService, err := do.Invoke[port.ProjectService](c.injector)
	if err != nil {
		return fmt.Errorf("invoke project service: %w", err)
	}

	entryPoints, err := projectService.GetAllEntryPoints(cmd.Context())
	if err != nil {
		return fmt.Errorf("project service: get all entry points: %w", err)
	}

	seperator := lo.Ternary(*c.flagLine, "\n", " ")

	fmt.Println(strings.Join(entryPoints, seperator))

	return nil
}
