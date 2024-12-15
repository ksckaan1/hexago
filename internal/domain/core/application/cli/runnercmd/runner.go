package runnercmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/config"
	"github.com/ksckaan1/hexago/internal/customerrors"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/ksckaan1/hexago/internal/port"
)

var _ port.Commander = (*RunnerCommand)(nil)

type RunnerCommand struct {
	cmd            *cobra.Command
	tuilog         *tuilog.TUILog
	projectService ProjectService
	cfg            *config.Config

	// flags
	flagEnvVars *[]string
	flagVerbose *bool
}

func NewRunnerCommand(projectService ProjectService, cfg *config.Config, tl *tuilog.TUILog) (*RunnerCommand, error) {
	return &RunnerCommand{
		cmd: &cobra.Command{
			Use:     "run",
			Example: "hexago run",
			Short:   "Runner processes",
			Long:    ``,
		},
		projectService: projectService,
		tuilog:         tl,
		cfg:            cfg,
	}, nil
}

func (c *RunnerCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *RunnerCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}

func (c *RunnerCommand) init() {
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := c.runner(cmd, args)
		if err != nil {
			return customerrors.ErrSuppressed
		}
		return nil
	}
	c.flagEnvVars = c.cmd.Flags().StringSliceP("env", "e", nil, "hexago run <runner> -e <KEY1>=<VALUE1> -e <KEY1>=<VALUE1>")
	c.flagVerbose = c.cmd.Flags().BoolP("verbose", "v", false, "hexago run <runner> -v")
}

func (c *RunnerCommand) runner(cmd *cobra.Command, args []string) error {
	err := c.cfg.Load()
	if err != nil {

		c.tuilog.Error(err.Error())

		return fmt.Errorf("cfg.Load: %w", err)
	}

	err = c.projectService.Run(cmd.Context(), args[0], *c.flagEnvVars, *c.flagVerbose)
	if err != nil {

		c.tuilog.Error(err.Error())

		return fmt.Errorf("projectService.Run: %w", err)
	}

	return nil
}
