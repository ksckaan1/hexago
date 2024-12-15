package infracmd

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/internal/customerrors"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/ksckaan1/hexago/internal/port"
)

var _ port.Commander = (*InfraLSCommand)(nil)

type InfraLSCommand struct {
	cmd            *cobra.Command
	tuilog         *tuilog.TUILog
	projectService ProjectService

	// flags
	flagLine *bool
}

func NewInfraLSCommand(projectService ProjectService, tl *tuilog.TUILog) (*InfraLSCommand, error) {
	return &InfraLSCommand{
		cmd: &cobra.Command{
			Use:     "ls",
			Example: "hexago infra ls",
			Short:   "List Infrastructures",
			Long:    `List Infrastructures`,
		},
		projectService: projectService,
		tuilog:         tl,
	}, nil
}

func (c *InfraLSCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *InfraLSCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}

func (c *InfraLSCommand) init() {
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := c.runner(cmd, args)
		if err != nil {
			return customerrors.ErrSuppressed
		}
		return nil
	}
	c.flagLine = c.cmd.Flags().BoolP("line", "l", false, "hexago infra ls -l")
}

func (c *InfraLSCommand) runner(cmd *cobra.Command, _ []string) error {
	infras, err := c.projectService.GetAllInfrastructures(cmd.Context())
	if err != nil {

		c.tuilog.Error(err.Error())

		return fmt.Errorf("projectService.GetAllInfrastructures: %w", err)
	}

	separator := lo.Ternary(*c.flagLine, "\n", " ")

	fmt.Println(strings.Join(infras, separator))

	return nil
}
