package domaincmd

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/internal/customerrors"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/ksckaan1/hexago/internal/port"
)

var _ port.Commander = (*DomainLSCommand)(nil)

type DomainLSCommand struct {
	cmd            *cobra.Command
	tuilog         *tuilog.TUILog
	projectService ProjectService

	// flags
	flagLine *bool
}

const domainLSLong = `ls command lists domains in project.

Domains are located under the "internal/domain/" directory.`

func NewDomainLSCommand(projectService ProjectService, tl *tuilog.TUILog) (*DomainLSCommand, error) {
	return &DomainLSCommand{
		cmd: &cobra.Command{
			Use:     "ls",
			Example: "hexago domain ls",
			Short:   "List domains",
			Long:    domainLSLong,
		},
		projectService: projectService,
		tuilog:         tl,
	}, nil
}

func (c *DomainLSCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *DomainLSCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}

func (c *DomainLSCommand) init() {
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := c.runner(cmd, args)
		if err != nil {
			return customerrors.ErrSuppressed
		}
		return nil
	}
	c.flagLine = c.cmd.Flags().BoolP("line", "l", false, "hexago domain ls -l")
}

func (c *DomainLSCommand) runner(cmd *cobra.Command, _ []string) error {
	domains, err := c.projectService.GetAllDomains(cmd.Context())
	if err != nil {

		c.tuilog.Error(err.Error())

		return fmt.Errorf("projectService.GetAllDomains: %w", err)
	}

	separator := lo.Ternary(*c.flagLine, "\n", " ")

	fmt.Println(strings.Join(domains, separator))

	return nil
}
