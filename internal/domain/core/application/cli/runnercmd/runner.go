package runnercmd

import (
	"fmt"
	"github.com/ksckaan1/hexago/internal/port"

	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type Commander interface {
	Command() *cobra.Command
}

type RunnerCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
	tuilog   *tuilog.TUILog

	// flags
	flagEnvVars *[]string
	flagVerbose *bool
}

func NewRunnerCommand(i *do.Injector) (*RunnerCommand, error) {
	return &RunnerCommand{
		cmd: &cobra.Command{
			Use:     "run",
			Example: "hexago run",
			Short:   "Runner processes",
			Long:    ``,
		},
		injector: i,
		tuilog:   do.MustInvoke[*tuilog.TUILog](i),
	}, nil
}

func (c *RunnerCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *RunnerCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *RunnerCommand) init() {
	c.cmd.RunE = c.runner
	c.flagEnvVars = c.cmd.Flags().StringSliceP("env", "e", nil, "hexago run <runner> -e <KEY1>=<VALUE1> -e <KEY1>=<VALUE1>")
	c.flagVerbose = c.cmd.Flags().BoolP("verbose", "v", false, "hexago run <runner> -v")
}

func (c *RunnerCommand) runner(cmd *cobra.Command, args []string) error {
	projectService, err := do.Invoke[port.ProjectService](c.injector)
	if err != nil {
		return fmt.Errorf("invoke project service: %w", err)
	}

	cfg, err := do.Invoke[port.ConfigService](c.injector)
	if err != nil {
		return fmt.Errorf("invoke config service: %w", err)
	}

	err = cfg.Load(".hexago/config.yaml")
	if err != nil {
		fmt.Println("")
		c.tuilog.Error(err.Error())
		fmt.Println("")
		return fmt.Errorf("load config: %w", err)
	}

	err = projectService.Run(cmd.Context(), args[0], *c.flagEnvVars, *c.flagVerbose)
	if err != nil {
		fmt.Println("")
		c.tuilog.Error(err.Error())
		fmt.Println("")
		return fmt.Errorf("project service: run: %w", err)
	}

	return nil
}
