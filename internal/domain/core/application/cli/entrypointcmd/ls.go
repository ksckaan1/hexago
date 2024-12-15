package entrypointcmd

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/internal/customerrors"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/ksckaan1/hexago/internal/port"
)

var _ port.Commander = (*EntryPointLSCommand)(nil)

type EntryPointLSCommand struct {
	cmd            *cobra.Command
	tuilog         *tuilog.TUILog
	projectService ProjectService

	// flags
	flagLine *bool
}

func NewEntryPointLSCommand(projectService ProjectService, tl *tuilog.TUILog) (*EntryPointLSCommand, error) {
	return &EntryPointLSCommand{
		cmd: &cobra.Command{
			Use:     "ls",
			Example: "hexago cmd ls",
			Short:   "List Entry Points",
			Long:    `List Entry Points`,
		},
		projectService: projectService,
		tuilog:         tl,
	}, nil
}

func (c *EntryPointLSCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *EntryPointLSCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}

func (c *EntryPointLSCommand) init() {
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := c.runner(cmd, args)
		if err != nil {
			return customerrors.ErrSuppressed
		}
		return nil
	}
	c.flagLine = c.cmd.Flags().BoolP("line", "l", false, "hexago service ls -l")
}

func (c *EntryPointLSCommand) runner(cmd *cobra.Command, _ []string) error {
	entryPoints, err := c.projectService.GetAllEntryPoints(cmd.Context())
	if err != nil {

		c.tuilog.Error(err.Error())

		return fmt.Errorf("projectService.GetAllEntryPoints: %w", err)
	}

	seperator := lo.Ternary(*c.flagLine, "\n", " ")

	fmt.Println(strings.Join(entryPoints, seperator))

	return nil
}
