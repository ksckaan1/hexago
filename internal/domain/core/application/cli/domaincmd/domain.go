package domaincmd

import (
	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/internal/port"
)

var _ port.Commander = (*DomainCommand)(nil)

type DomainCommand struct {
	cmd *cobra.Command
}

const domainLong = `domain command includes sub-commands for listing and creating domains.`

func NewDomainCommand() (*DomainCommand, error) {
	return &DomainCommand{
		cmd: &cobra.Command{
			Use:     "domain",
			Example: "hexago domain",
			Short:   "Domain processes",
			Long:    domainLong,
		},
	}, nil
}

func (c *DomainCommand) Command() *cobra.Command {
	return c.cmd
}

func (c *DomainCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}
