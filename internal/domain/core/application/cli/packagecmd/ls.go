package packagecmd

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

type PackageLSCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
	tuilog   *tuilog.TUILog

	// flags
	flagLine   *bool
	flagGlobal *bool
	flagAll    *bool
}

func NewPackageLSCommand(i *do.Injector) (*PackageLSCommand, error) {
	return &PackageLSCommand{
		cmd: &cobra.Command{
			Use:     "ls",
			Example: "hexago pkg ls",
			Short:   "List Packages",
			Long:    `List Packages`,
		},
		injector: i,
		tuilog:   do.MustInvoke[*tuilog.TUILog](i),
	}, nil
}

func (c *PackageLSCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *PackageLSCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *PackageLSCommand) init() {
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := c.runner(cmd, args)
		if err != nil {
			return dto.ErrSuppressed
		}
		return nil
	}
	c.flagLine = c.cmd.Flags().BoolP("line", "l", false, "hexago pkg ls -l")
	c.flagGlobal = c.cmd.Flags().BoolP("global", "g", false, "hexago pkg ls -g")
	c.flagAll = c.cmd.Flags().BoolP("all", "a", false, "hexago pkg ls -a")
}

func (c *PackageLSCommand) runner(cmd *cobra.Command, _ []string) error {
	projectService, err := do.Invoke[port.ProjectService](c.injector)
	if err != nil {
		return fmt.Errorf("invoke project service: %w", err)
	}

	allPackages := make([]string, 0)

	if *c.flagAll {
		globalPackages, err := projectService.GetAllPackages(cmd.Context(), true)
		if err != nil {
			fmt.Println("")
			c.tuilog.Error(err.Error())
			fmt.Println("")
			return fmt.Errorf("project service: get all packages: %w", err)
		}

		globalPackages = lo.Map(globalPackages, func(p string, _ int) string {
			return "global:" + p
		})

		allPackages = append(allPackages, globalPackages...)

		packages, err := projectService.GetAllPackages(cmd.Context(), false)
		if err != nil {
			fmt.Println("")
			c.tuilog.Error(err.Error())
			fmt.Println("")
			return fmt.Errorf("project service: get all packages: %w", err)
		}

		packages = lo.Map(packages, func(p string, _ int) string {
			return "internal:" + p
		})

		allPackages = append(allPackages, packages...)

	} else {
		packages, err := projectService.GetAllPackages(cmd.Context(), *c.flagGlobal)
		if err != nil {
			fmt.Println("")
			c.tuilog.Error(err.Error())
			fmt.Println("")
			return fmt.Errorf("project service: get all packages: %w", err)
		}

		allPackages = append(allPackages, packages...)
	}

	seperator := lo.Ternary(*c.flagLine, "\n", " ")

	fmt.Println(strings.Join(allPackages, seperator))

	return nil
}
