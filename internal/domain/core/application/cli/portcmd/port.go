package portcmd

import (
	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/internal/port"
)

var _ port.Commander = (*PortCommand)(nil)

type PortCommand struct {
	cmd *cobra.Command
}

func NewPortCommand() (*PortCommand, error) {
	return &PortCommand{
		cmd: &cobra.Command{
			Use:     "port",
			Example: "hexago port",
			Short:   "Port processes",
			Long:    `Port processes`,
		},
	}, nil
}

func (c *PortCommand) Command() *cobra.Command {
	return c.cmd
}

func (c *PortCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}
