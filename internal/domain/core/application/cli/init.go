package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/ksckaan1/hexago/internal/domain/core/dto"
	"github.com/ksckaan1/hexago/internal/domain/core/port"
	"github.com/ksckaan1/hexago/internal/util"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var _ Commander = (*InitCommand)(nil)

type InitCommand struct {
	cmd      *cobra.Command
	injector *do.Injector

	// Flags
	flagModuleName *string
}

func NewInitCommand(i *do.Injector) (*InitCommand, error) {
	return &InitCommand{
		cmd: &cobra.Command{
			Use:     "init",
			Example: "hexago init <project-name>",
			Short:   "Initialize a hexagonal Go project",
			Long:    `Initialize a hexagonal Go project`,
			Args:    cobra.ExactArgs(1),
		},
		injector: i,
		// flags
		flagModuleName: new(string),
	}, nil
}

func (c *InitCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *InitCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *InitCommand) init() {
	c.cmd.RunE = c.runner
	c.flagModuleName = c.cmd.Flags().StringP("module-name", "m", "", "hexago init <project-name> -m <module-name>")
}

func (c *InitCommand) runner(cmd *cobra.Command, args []string) error {
	if *c.flagModuleName == "" {
		defaultName := filepath.Base(args[0])

		if defaultName == "." {
			abs, err := filepath.Abs(defaultName)
			if err != nil {
				return fmt.Errorf("filepath: abs: %w", err)
			}

			defaultName = filepath.Base(abs)
		}

		err := huh.NewInput().
			Title("What’s module name?").
			Placeholder(defaultName).
			Description("If you leave it blank, it will have the same name as the project directory").
			Value(c.flagModuleName).
			Run()
		if err != nil {
			return fmt.Errorf("input module name: %w", err)
		}

		if strings.TrimSpace(*c.flagModuleName) == "" {
			*c.flagModuleName = defaultName
		}
	}

	projectService, err := do.Invoke[port.ProjectService](c.injector)
	if err != nil {
		return fmt.Errorf("invoke project service: %w", err)
	}

	err = projectService.InitNewProject(cmd.Context(), dto.InitNewProjectParams{
		ProjectDirectory: args[0],
		ModuleName:       *c.flagModuleName,
	})
	if err != nil {
		return fmt.Errorf("project service: init new project: %w", err)
	}

	msg := "Ready to go!"
	if args[0] != "." {
		msg = "cd " + args[0]
	}

	fmt.Println("")
	util.UILog(util.Success, msg, "Project Initialized")
	fmt.Println("")

	return nil
}
