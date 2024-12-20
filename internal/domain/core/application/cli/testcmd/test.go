package testcmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/internal/customerrors"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/ksckaan1/hexago/internal/port"
)

type ProjectService interface {
	GetModuleName(modulePath ...string) (string, error)
}

var _ port.Commander = (*TestCommand)(nil)

type TestCommand struct {
	cmd            *cobra.Command
	tuilog         *tuilog.TUILog
	projectService ProjectService

	// flags
	flagVerbose *bool

	// state
	moduleName string
}

func NewTestCommand(projectService ProjectService, tl *tuilog.TUILog) (*TestCommand, error) {
	return &TestCommand{
		cmd: &cobra.Command{
			Use:     "test",
			Example: "hexago test",
			Short:   "Run tests",
			Long:    `Run tests`,
		},
		projectService: projectService,
		tuilog:         tl,
	}, nil
}

func (c *TestCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *TestCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}

func (c *TestCommand) init() {
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := c.runner(cmd, args)
		if err != nil {
			return customerrors.ErrSuppressed
		}
		return nil
	}

	c.flagVerbose = c.cmd.Flags().BoolP("verbose", "v", false, "hexago test ./... -v")
}

type TestAction struct {
	Time    time.Time
	Action  string
	Package string
	Test    string
	Output  string
	Elapsed float64
}

func (c *TestCommand) runner(cmd *cobra.Command, args []string) error {
	moduleName, err := c.projectService.GetModuleName()
	if err != nil {
		c.tuilog.Error(err.Error())
		return err
	}

	c.moduleName = moduleName

	cmdArgs := []string{"test", "--json"}
	cmdArgs = append(cmdArgs, args...)
	subCmd := exec.CommandContext(cmd.Context(), "go", cmdArgs...)

	r, w := io.Pipe()

	subCmd.Stdout = w
	subCmd.Stderr = w

	sc := bufio.NewScanner(r)

	showErr := false

	go func() {

		for sc.Scan() {
			err := sc.Err()
			if err != nil {
				c.tuilog.Error(err.Error())
				break
			}

			var testAction TestAction
			err = json.Unmarshal(sc.Bytes(), &testAction)
			if err != nil {
				if !showErr {
					c.printError(sc.Text(), true)
				} else {
					c.printError(sc.Text(), false)
				}
				showErr = true
				continue
			}

			showErr = false

			c.printAction(&testAction)
		}
	}()

	err = subCmd.Run()
	if err != nil {
		return customerrors.ErrSuppressed
	}
	return nil
}

func (c *TestCommand) printAction(action *TestAction) {
	if action.Test == "" {
		return
	}

	action.Package = strings.TrimPrefix(action.Package, c.moduleName+"/")

	switch action.Action {
	case "run":
		if *c.flagVerbose {
			c.printRun(action)
		}
	case "pass":
		if *c.flagVerbose || c.isRootTest(action) {
			c.printResult(action, "#FFFFFF", "#008000", " PASS ")
		}
	case "fail":
		c.printResult(action, "#000000", "#FF0000", " FAIL ")
	case "skip":
		if *c.flagVerbose {
			c.printResult(action, "#FFFFFF", "#411F89", " SKIP ")
		}
	case "output":
		if *c.flagVerbose {
			c.printOutput(action)
		}
	}
}

func (c *TestCommand) isRootTest(action *TestAction) bool {
	return !strings.Contains(action.Test, "/")
}

func (c *TestCommand) printRun(action *TestAction) {
	badge := lipgloss.NewStyle().
		Background(lipgloss.Color("#87CEEB")).
		Foreground(lipgloss.Color("#000000")).
		Render(" RUN  ")

	pkg := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#808080")).
		Render(action.Package)
	test := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Render(action.Test)

	fmt.Printf("%s %s %s\n", badge, pkg, test)
}

func (c *TestCommand) printResult(action *TestAction, fgColor, bgColor, badgeText string) {
	badge := lipgloss.NewStyle().
		Background(lipgloss.Color(bgColor)).
		Foreground(lipgloss.Color(fgColor)).
		Render(badgeText)

	pkg := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#808080")).
		Render(action.Package)
	test := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Render(action.Test)

	duration := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#808080")).
		Render(fmt.Sprintf(" (%.2fs)", action.Elapsed))
	fmt.Printf("%s %s %s %s\n", badge, pkg, test, duration)
}

func (c *TestCommand) printOutput(action *TestAction) {
	if strings.HasPrefix(action.Output, "---") ||
		strings.HasPrefix(action.Output, "===") {
		return
	}
	fmt.Print(action.Output)
}

func (c *TestCommand) printError(text string, withBadge bool) {
	title := "      "
	if withBadge {
		title = " ERR  "
	}
	badge := lipgloss.NewStyle().
		Background(lipgloss.Color("#FF0000")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Render(title)
	fmt.Printf("%s %s\n", badge, text)
}
