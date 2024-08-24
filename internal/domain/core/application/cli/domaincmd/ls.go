package domaincmd

import (
	"fmt"
	"strings"

	"github.com/ksckaan1/hexago/internal/domain/core/port"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type DomainLSCommand struct {
	cmd      *cobra.Command
	injector *do.Injector

	// flags
	flagLine *bool
}

const domainLSLong = `ls command lists domains in project.

Domains are located under the "internal/domain/" directory.`

const domainLSExamples = `hexago domain ls
hexago domain ls -l`

func NewDomainLSCommand(i *do.Injector) (*DomainLSCommand, error) {
	return &DomainLSCommand{
		cmd: &cobra.Command{
			Use:     "ls",
			Example: "hexago domain ls",
			Short:   "List domains",
			Long:    domainLSLong,
		},
		injector: i,
		// flags
		flagLine: new(bool),
	}, nil
}

func (c *DomainLSCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *DomainLSCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *DomainLSCommand) init() {
	c.cmd.RunE = c.runner
	c.flagLine = c.cmd.Flags().BoolP("line", "l", false, "hexago domain ls -l")
}

func (c *DomainLSCommand) runner(cmd *cobra.Command, _ []string) error {
	projectService, err := do.Invoke[port.ProjectService](c.injector)
	if err != nil {
		return fmt.Errorf("invoke project service: %w", err)
	}

	domains, err := projectService.GetAllDomains(cmd.Context())
	if err != nil {
		return fmt.Errorf("project service: get all domains: %w", err)
	}

	seperator := lo.Ternary(*c.flagLine, "\n", " ")

	fmt.Println(strings.Join(domains, seperator))

	return nil
}
