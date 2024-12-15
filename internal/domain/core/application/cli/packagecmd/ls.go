package packagecmd

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/internal/customerrors"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/ksckaan1/hexago/internal/port"
)

var _ port.Commander = (*PackageLSCommand)(nil)

type PackageLSCommand struct {
	cmd            *cobra.Command
	tuilog         *tuilog.TUILog
	projectService ProjectService

	// flags
	flagLine   *bool
	flagGlobal *bool
	flagAll    *bool
}

func NewPackageLSCommand(projectService ProjectService, tl *tuilog.TUILog) (*PackageLSCommand, error) {
	return &PackageLSCommand{
		cmd: &cobra.Command{
			Use:     "ls",
			Example: "hexago pkg ls",
			Short:   "List Packages",
			Long:    `List Packages`,
		},
		tuilog:         tl,
		projectService: projectService,
	}, nil
}

func (c *PackageLSCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *PackageLSCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}

func (c *PackageLSCommand) init() {
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := c.runner(cmd, args)
		if err != nil {
			return customerrors.ErrSuppressed
		}
		return nil
	}
	c.flagLine = c.cmd.Flags().BoolP("line", "l", false, "hexago pkg ls -l")
	c.flagGlobal = c.cmd.Flags().BoolP("global", "g", false, "hexago pkg ls -g")
	c.flagAll = c.cmd.Flags().BoolP("all", "a", false, "hexago pkg ls -a")
}

func (c *PackageLSCommand) runner(cmd *cobra.Command, _ []string) error {
	allPackages := make([]string, 0)

	if *c.flagAll {
		globalPackages, err := c.projectService.GetAllPackages(cmd.Context(), true)
		if err != nil {

			c.tuilog.Error(err.Error())

			return fmt.Errorf("projectService.GetAllPackages: %w", err)
		}

		globalPackages = lo.Map(globalPackages, func(p string, _ int) string {
			return "global:" + p
		})

		allPackages = append(allPackages, globalPackages...)

		packages, err := c.projectService.GetAllPackages(cmd.Context(), false)
		if err != nil {

			c.tuilog.Error(err.Error())

			return fmt.Errorf("projectService.GetAllPackages: %w", err)
		}

		packages = lo.Map(packages, func(p string, _ int) string {
			return "internal:" + p
		})

		allPackages = append(allPackages, packages...)

	} else {
		packages, err := c.projectService.GetAllPackages(cmd.Context(), *c.flagGlobal)
		if err != nil {

			c.tuilog.Error(err.Error())

			return fmt.Errorf("projectService.GetAllPackages: %w", err)
		}

		allPackages = append(allPackages, packages...)
	}

	separator := lo.Ternary(*c.flagLine, "\n", " ")

	fmt.Println(strings.Join(allPackages, separator))

	return nil
}
