package portcmd

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/internal/customerrors"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/ksckaan1/hexago/internal/port"
)

var _ port.Commander = (*PortLSCommand)(nil)

type PortLSCommand struct {
	cmd            *cobra.Command
	tuilog         *tuilog.TUILog
	projectService ProjectService

	// flags
	flagLine *bool
}

func NewPortLSCommand(prokectService ProjectService, tl *tuilog.TUILog) (*PortLSCommand, error) {
	return &PortLSCommand{
		cmd: &cobra.Command{
			Use:     "ls",
			Example: "hexago port ls -d <domainname>\nhexago port ls (select domain interactively)",
			Short:   "List ports",
			Long:    `List ports`,
		},
		projectService: prokectService,
		tuilog:         tl,
	}, nil
}

func (c *PortLSCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *PortLSCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}

func (c *PortLSCommand) init() {
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := c.runner(cmd, args)
		if err != nil {
			return customerrors.ErrSuppressed
		}
		return nil
	}
	c.flagLine = c.cmd.Flags().BoolP("line", "l", false, "hexago port ls -l")
}

func (c *PortLSCommand) runner(cmd *cobra.Command, _ []string) error {
	allPorts, err := c.projectService.GetAllPorts(cmd.Context())
	if err != nil {

		c.tuilog.Error(err.Error())

		return fmt.Errorf("projectService.GetAllPorts: %w", err)
	}

	separator := lo.Ternary(*c.flagLine, "\n", " ")

	fmt.Println(strings.Join(allPorts, separator))

	return nil
}
