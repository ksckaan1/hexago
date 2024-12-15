package servicecmd

import (
	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/internal/port"
)

var _ port.Commander = (*ServiceCommand)(nil)

type ServiceCommand struct {
	cmd *cobra.Command
}

func NewServiceCommand() (*ServiceCommand, error) {
	return &ServiceCommand{
		cmd: &cobra.Command{
			Use:     "service",
			Example: "hexago service",
			Short:   "Service processes",
			Long:    `Service processes`,
		},
	}, nil
}

func (c *ServiceCommand) Command() *cobra.Command {
	return c.cmd
}

func (c *ServiceCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}
