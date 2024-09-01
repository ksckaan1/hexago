package doctorcmd

import (
	"fmt"

	"github.com/ksckaan1/hexago/internal/domain/core/port"
	"github.com/ksckaan1/hexago/internal/util"
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
	c.cmd.RunE = c.runner
}

func (c *DoctorCommand) runner(cmd *cobra.Command, args []string) error {
	projectService, err := do.Invoke[port.ProjectService](c.injector)
	if err != nil {
		return fmt.Errorf("invoke project service: %w", err)
	}

	result, err := projectService.Doctor(cmd.Context())
	if err != nil {
		return fmt.Errorf("project service: doctor: %w", err)
	}

	fmt.Println("")
	util.UILog(util.Info, result.OSResult)
	fmt.Println("")

	if result.GoResult != "" {
		util.UILog(util.Success, result.GoResult, "go")
	} else {
		util.UILog(util.Error, "not found", "go")
	}
	fmt.Println("")

	if result.ImplResult != "" {
		util.UILog(util.Success, "installed", "impl")
	} else {
		util.UILog(util.Error, "not found", "impl")
	}
	fmt.Println("")

	return nil
}
