package appcmd

import (
	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/internal/port"
)

var _ port.Commander = (*AppCommand)(nil)

type AppCommand struct {
	cmd *cobra.Command
}

func NewAppCommand() (*AppCommand, error) {
	return &AppCommand{
		cmd: &cobra.Command{
			Use:     "app",
			Example: "hexago app",
			Aliases: []string{"a"},
			Short:   "Application processes",
			Long:    `Application processes`,
		},
	}, nil
}

func (c *AppCommand) Command() *cobra.Command {
	return c.cmd
}

func (c *AppCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}
