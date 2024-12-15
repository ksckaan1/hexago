package infracmd

import (
	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/internal/port"
)

var _ port.Commander = (*InfraCommand)(nil)

type InfraCommand struct {
	cmd *cobra.Command
}

func NewInfraCommand() (*InfraCommand, error) {
	return &InfraCommand{
		cmd: &cobra.Command{
			Use:     "infra",
			Example: "hexago infra",
			Short:   "Infrastructure processes",
			Long:    `Infrastructure processes`,
		},
	}, nil
}

func (c *InfraCommand) Command() *cobra.Command {
	return c.cmd
}

func (c *InfraCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}
