package entrypointcmd

import (
	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/internal/port"
)

var _ port.Commander = (*EntryPointCommand)(nil)

type EntryPointCommand struct {
	cmd *cobra.Command
}

func NewEntryPointCommand() (*EntryPointCommand, error) {
	return &EntryPointCommand{
		cmd: &cobra.Command{
			Use:     "cmd",
			Example: "hexago cmd",
			Short:   "Entry Point processes",
			Long:    `Entry Point processes`,
		},
	}, nil
}

func (c *EntryPointCommand) Command() *cobra.Command {
	return c.cmd
}

func (c *EntryPointCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}
