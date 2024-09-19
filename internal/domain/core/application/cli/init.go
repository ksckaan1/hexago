package cli

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/ksckaan1/hexago/internal/domain/core/dto"
	"github.com/ksckaan1/hexago/internal/domain/core/port"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var _ Commander = (*InitCommand)(nil)

type InitCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
	tuilog   *tuilog.TUILog
}

const initLongDescription = `init command initialize a hexagonal Go project.
This command creates a folder called .hexago. This folder includes a config.yaml file that setting up general hexago options.
In the other side, it initialize traditional Go hexagonal folder structure.

Requires empty folder to init new project.`

const initExamples = `hexago init <project-name> (prompts module name interactively)
`

func NewInitCommand(i *do.Injector) (*InitCommand, error) {
	return &InitCommand{
		cmd: &cobra.Command{
			Use:     "init",
			Example: initExamples,
			Short:   "Initialize a hexagonal Go project",
			Long:    initLongDescription,
			Args:    cobra.ExactArgs(1),
		},
		injector: i,
		tuilog:   do.MustInvoke[*tuilog.TUILog](i),
		// flags
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
}

func (c *InitCommand) runner(cmd *cobra.Command, args []string) error {
	projectService, err := do.Invoke[port.ProjectService](c.injector)
	if err != nil {
		return fmt.Errorf("invoke project service: %w", err)
	}

	createModule := true

	existingModuleName, err := projectService.GetModuleName(filepath.Join(args[0], "go.mod"))
	if err == nil {
		err = huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Do you want to overwrite the existing module?").
					Description("existing -> " + existingModuleName).
					Affirmative("Yes").
					Negative("No").
					Value(&createModule),
			).WithShowHelp(true),
		).Run()
		if err != nil {
			return fmt.Errorf("confirm module create: %w", err)
		}
	}

	var moduleName string
	if createModule {
		defaultModuleName := filepath.Base(args[0])

		if defaultModuleName == "." {
			abs, err := filepath.Abs(defaultModuleName)
			if err != nil {
				return fmt.Errorf("filepath: abs: %w", err)
			}

			defaultModuleName = filepath.Base(abs)
		}

		err = huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("What’s module name?").
					Placeholder(defaultModuleName).
					Description("If you leave it blank, it will have the same name\nas the project directory").
					Value(&moduleName),
			).WithShowHelp(true),
		).Run()
		if err != nil {
			return fmt.Errorf("input module name: %w", err)
		}

		if strings.TrimSpace(moduleName) == "" {
			moduleName = defaultModuleName
		}
	}

	err = projectService.InitNewProject(cmd.Context(), dto.InitNewProjectParams{
		ProjectDirectory: args[0],
		ModuleName:       moduleName,
		CreateModule:     createModule,
	})

	if err != nil {
		fmt.Println("")
		if errors.Is(err, dto.ErrDirMustBeFolder) {
			c.tuilog.Error("project dir must be folder")
		} else if err2, ok := lo.ErrorsAs[dto.ErrInitGoModule](err); ok {
			c.tuilog.Error(err2.Message)
		} else {
			c.tuilog.Error(err.Error())
		}
		fmt.Println("")
		return fmt.Errorf("project service: init new project: %w", err)
	}

	msg := "Ready to go!"
	if args[0] != "." {
		msg = "cd " + args[0]
	}

	fmt.Println("")
	c.tuilog.Success(msg, "Project Initialized")
	fmt.Println("")

	return nil
}
