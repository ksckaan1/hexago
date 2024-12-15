package initcmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/internal/customerrors"
	"github.com/ksckaan1/hexago/internal/domain/core/model"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/ksckaan1/hexago/internal/port"
)

var _ port.Commander = (*InitCommand)(nil)

type InitCommand struct {
	cmd            *cobra.Command
	tuilog         *tuilog.TUILog
	projectService ProjectService
}

const initLongDescription = `init command initialize a hexagonal Go project.
This command creates a folder called .hexago. This folder includes a config.yaml file that setting up general hexago options.
In the other side, it initialize traditional Go hexagonal folder structure.

Requires empty folder to init new project.`

const initExamples = `hexago init <project-name> (prompts module name interactively)
`

func NewInitCommand(projectService ProjectService, tl *tuilog.TUILog) (*InitCommand, error) {
	return &InitCommand{
		cmd: &cobra.Command{
			Use:     "init",
			Example: initExamples,
			Short:   "Initialize a hexagonal Go project",
			Long:    initLongDescription,
			Args:    cobra.ExactArgs(1),
		},
		tuilog:         tl,
		projectService: projectService,
	}, nil
}

func (c *InitCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *InitCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}

func (c *InitCommand) init() {
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := c.runner(cmd, args)
		if err != nil {
			return customerrors.ErrSuppressed
		}
		return nil
	}
}

func (c *InitCommand) runner(cmd *cobra.Command, args []string) error {
	createModule := true

	existingModuleName, err := c.projectService.GetModuleName(filepath.Join(args[0], "go.mod"))
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
			abs, err2 := filepath.Abs(defaultModuleName)
			if err2 != nil {
				return fmt.Errorf("filepath: abs: %w", err2)
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

	err = c.projectService.InitNewProject(cmd.Context(), model.InitNewProjectParams{
		ProjectDirectory: args[0],
		ModuleName:       moduleName,
		CreateModule:     createModule,
	})
	if err != nil {

		if errors.Is(err, customerrors.ErrDirMustBeFolder) {
			c.tuilog.Error("project dir must be folder")
		} else if err2, ok := lo.ErrorsAs[customerrors.ErrInitGoModule](err); ok {
			c.tuilog.Error(err2.Message)
		} else {
			c.tuilog.Error(err.Error())
		}

		return fmt.Errorf("projectService.InitNewProject: %w", err)
	}

	msg := "Ready to go!"
	if args[0] != "." {
		msg = "cd " + args[0]
	}

	c.tuilog.Success(msg, "Project Initialized")

	return nil
}
