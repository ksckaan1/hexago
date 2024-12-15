package packagecmd

import (
	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/internal/port"
)

var _ port.Commander = (*PackageCommand)(nil)

type PackageCommand struct {
	cmd *cobra.Command
}

func NewPackageCommand() (*PackageCommand, error) {
	return &PackageCommand{
		cmd: &cobra.Command{
			Use:     "pkg",
			Example: "hexago pkg",
			Short:   "Package processes",
			Long:    `Package processes`,
		},
	}, nil
}

func (c *PackageCommand) Command() *cobra.Command {
	return c.cmd
}

func (c *PackageCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}
