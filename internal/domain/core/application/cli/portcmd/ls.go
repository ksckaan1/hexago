package portcmd

import (
	"fmt"
	"github.com/ksckaan1/hexago/internal/domain/core/dto"
	"github.com/ksckaan1/hexago/internal/port"
	"strings"

	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type PortLSCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
	tuilog   *tuilog.TUILog

	// flags
	flagLine *bool
}

func NewPortLSCommand(i *do.Injector) (*PortLSCommand, error) {
	return &PortLSCommand{
		cmd: &cobra.Command{
			Use:     "ls",
			Example: "hexago port ls -d <domainname>\nhexago port ls (select domain interactively)",
			Short:   "List ports",
			Long:    `List ports`,
		},
		injector: i,
		tuilog:   do.MustInvoke[*tuilog.TUILog](i),
		// flags
		flagLine: new(bool),
	}, nil
}

func (c *PortLSCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *PortLSCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *PortLSCommand) init() {
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := c.runner(cmd, args)
		if err != nil {
			return dto.ErrSuppressed
		}
		return nil
	}
	c.flagLine = c.cmd.Flags().BoolP("line", "l", false, "hexago port ls -l")
}

func (c *PortLSCommand) runner(cmd *cobra.Command, _ []string) error {
	projectService, err := do.Invoke[port.ProjectService](c.injector)
	if err != nil {
		return fmt.Errorf("invoke project service: %w", err)
	}

	allPorts, err := projectService.GetAllPorts(cmd.Context())
	if err != nil {
		fmt.Println("")
		c.tuilog.Error(err.Error())
		fmt.Println("")
		return fmt.Errorf("project service: get all ports: %w", err)
	}

	seperator := lo.Ternary(*c.flagLine, "\n", " ")

	fmt.Println(strings.Join(allPorts, seperator))

	return nil
}
