package doctorcmd

import (
	"fmt"
	"github.com/ksckaan1/hexago/internal/domain/core/dto"
	"github.com/ksckaan1/hexago/internal/port"

	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type Commander interface {
	Command() *cobra.Command
}

type DoctorCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
	tuilog   *tuilog.TUILog
}

const doctorLong = `doctor command check dependencies.`

func NewDoctorCommand(i *do.Injector) (*DoctorCommand, error) {
	return &DoctorCommand{
		cmd: &cobra.Command{
			Use:     "doctor",
			Example: "hexago doctor",
			Short:   "Check hexago command stability",
			Long:    doctorLong,
		},
		injector: i,
		tuilog:   do.MustInvoke[*tuilog.TUILog](i),
	}, nil
}

func (c *DoctorCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *DoctorCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *DoctorCommand) init() {
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := c.runner(cmd, args)
		if err != nil {
			return dto.ErrSuppressed
		}
		return nil
	}
}

func (c *DoctorCommand) runner(cmd *cobra.Command, _ []string) error {
	projectService, err := do.Invoke[port.ProjectService](c.injector)
	if err != nil {
		return fmt.Errorf("invoke project service: %w", err)
	}

	result, err := projectService.Doctor(cmd.Context())
	if err != nil {
		fmt.Println("")
		c.tuilog.Error(err.Error())
		fmt.Println("")
		return fmt.Errorf("project service: doctor: %w", err)
	}

	fmt.Println("")
	c.tuilog.Info(result.OSResult, "os/arch")
	fmt.Println("")

	if result.GoResult.IsInstalled {
		c.tuilog.Success(result.GoResult.Output, "go")
	} else {
		c.tuilog.Error(result.GoResult.Output, "go")
	}
	fmt.Println("")

	if result.ImplResult.IsInstalled {
		c.tuilog.Success("installed", "impl")
	} else {
		c.tuilog.Error(result.ImplResult.Output, "impl")
	}
	fmt.Println("")

	return nil
}
