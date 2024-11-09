package infracmd

import (
	"fmt"
	"github.com/ksckaan1/hexago/internal/port"
	"strings"

	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type InfraLSCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
	tuilog   *tuilog.TUILog

	// flags
	flagLine *bool
}

func NewInfraLSCommand(i *do.Injector) (*InfraLSCommand, error) {
	return &InfraLSCommand{
		cmd: &cobra.Command{
			Use:     "ls",
			Example: "hexago infra ls",
			Short:   "List Infrastructures",
			Long:    `List Infrastructures`,
		},
		injector: i,
		tuilog:   do.MustInvoke[*tuilog.TUILog](i),
	}, nil
}

func (c *InfraLSCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *InfraLSCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *InfraLSCommand) init() {
	c.cmd.RunE = c.runner
	c.flagLine = c.cmd.Flags().BoolP("line", "l", false, "hexago infra ls -l")
}

func (c *InfraLSCommand) runner(cmd *cobra.Command, _ []string) error {
	projectService, err := do.Invoke[port.ProjectService](c.injector)
	if err != nil {
		return fmt.Errorf("invoke project service: %w", err)
	}

	infras, err := projectService.GetAllInfrastructes(cmd.Context())
	if err != nil {
		fmt.Println("")
		c.tuilog.Error(err.Error())
		fmt.Println("")
		return fmt.Errorf("project service: get all infrastructures: %w", err)
	}

	seperator := lo.Ternary(*c.flagLine, "\n", " ")

	fmt.Println(strings.Join(infras, seperator))

	return nil
}
