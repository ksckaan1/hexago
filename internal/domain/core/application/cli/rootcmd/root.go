package rootcmd

import (
	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/internal/port"
)

var _ port.Commander = (*RootCommand)(nil)

type RootCommand struct {
	cmd *cobra.Command
}

func NewRootCommand() (*RootCommand, error) {
	return &RootCommand{
		cmd: &cobra.Command{
			Use:           "hexago",
			Short:         "short description",
			Long:          header + "\nhexago is a cli tool for initializing and managing hexagonal Go projects.",
			Version:       "v0.5.0",
			SilenceUsage:  true,
			SilenceErrors: true,
		},
	}, nil
}

func (c *RootCommand) Command() *cobra.Command {
	return c.cmd
}

func (c *RootCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}
