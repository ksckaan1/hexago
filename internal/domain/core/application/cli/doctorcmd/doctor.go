package doctorcmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/internal/customerrors"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/ksckaan1/hexago/internal/port"
)

var _ port.Commander = (*DoctorCommand)(nil)

type DoctorCommand struct {
	cmd            *cobra.Command
	tuilog         *tuilog.TUILog
	projectService ProjectService
}

const doctorLong = `doctor command check dependencies.`

func NewDoctorCommand(projectService ProjectService, tl *tuilog.TUILog) (*DoctorCommand, error) {
	return &DoctorCommand{
		cmd: &cobra.Command{
			Use:     "doctor",
			Example: "hexago doctor",
			Short:   "Check hexago command stability",
			Long:    doctorLong,
		},
		projectService: projectService,
		tuilog:         tl,
	}, nil
}

func (c *DoctorCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *DoctorCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}

func (c *DoctorCommand) init() {
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := c.runner(cmd, args)
		if err != nil {
			return customerrors.ErrSuppressed
		}
		return nil
	}
}

func (c *DoctorCommand) runner(cmd *cobra.Command, _ []string) error {
	result, err := c.projectService.Doctor(cmd.Context())
	if err != nil {

		c.tuilog.Error(err.Error())

		return fmt.Errorf("projectService.Doctor: %w", err)
	}

	c.tuilog.Info(result.OSResult, "os/arch")

	if result.GoResult.IsInstalled {
		c.tuilog.Success(result.GoResult.Output, "go")
	} else {
		c.tuilog.Error(result.GoResult.Output, "go")
	}

	if result.ImplResult.IsInstalled {
		c.tuilog.Success("installed", "impl")
	} else {
		c.tuilog.Error(result.ImplResult.Output, "impl")
	}

	return nil
}
