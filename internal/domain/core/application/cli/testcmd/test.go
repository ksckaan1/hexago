package testcmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
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
	showError  bool
	startTime  time.Time

	// total
	totalMut  *sync.Mutex
	totalRun  int
	totalPass int
	totalFail int
	totalSkip int
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
		totalMut:       new(sync.Mutex),
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
		return fmt.Errorf("projectService.GetModuleName: %w", err)
	}

	c.moduleName = moduleName

	cmdArgs := []string{"test", "--json"}
	cmdArgs = append(cmdArgs, args...)
	subCmd := exec.CommandContext(cmd.Context(), "go", cmdArgs...)

	r, w := io.Pipe()
	subCmd.Stdout = w
	subCmd.Stderr = w

	sc := bufio.NewScanner(r)

	err = subCmd.Start()
	if err != nil {
		return customerrors.ErrSuppressed
	}

	c.startTime = time.Now()

	go func() {
		for sc.Scan() {
			c.parseLine(sc.Text())
		}
	}()

	err = subCmd.Wait()
	if err != nil {
		c.printSummary(subCmd.ProcessState.ExitCode() == 0)
		return fmt.Errorf("subCmd.Wait: %w", err)
	}

	c.printSummary(true)

	return nil
}

func (c *TestCommand) parseLine(line string) {
	if strings.HasPrefix(line, "no Go files in") {
		c.showError = false
		c.tuilog.Error(line)
		return
	}

	var testAction TestAction
	err := json.Unmarshal([]byte(line), &testAction)
	if err == nil {
		c.showError = false
		c.printAction(&testAction)
		return
	}

	if !c.showError {
		c.showError = true
		c.printError(line, true)
	} else {
		c.printError(line, false)
	}
}

func (c *TestCommand) printAction(action *TestAction) {
	if action.Test == "" {
		return
	}

	action.Package = strings.TrimPrefix(action.Package, c.moduleName+"/")
	c.countAction(action.Action)
	switch action.Action {
	case "run":
		if *c.flagVerbose {
			c.printRun(action)
		}
	case "pass":
		if *c.flagVerbose || c.isRootTest(action) {
			c.printResult(action, ColorWhite, ColorGreen, " PASS ")
		}
	case "fail":
		c.printResult(action, ColorWhite, ColorRed, " FAIL ")
	case "skip":
		c.printResult(action, ColorWhite, ColorPurple, " SKIP ")
	case "output":
		if *c.flagVerbose {
			c.printOutput(action)
		}
	}
}

func (c *TestCommand) countAction(action string) {
	c.totalMut.Lock()
	defer c.totalMut.Unlock()
	switch action {
	case "run":
		c.totalRun++
	case "pass":
		c.totalPass++
	case "fail":
		c.totalFail++
	case "skip":
		c.totalSkip++
	}
}

func (c *TestCommand) isRootTest(action *TestAction) bool {
	return !strings.Contains(action.Test, "/")
}

func (c *TestCommand) printRun(action *TestAction) {
	badge := lipgloss.NewStyle().
		Background(lipgloss.Color(ColorBlue)).
		Foreground(lipgloss.Color(ColorBlack)).
		Render(" RUN  ")

	pkg := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorGray)).
		Render(action.Package)
	test := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorWhite)).
		Render(action.Test)

	fmt.Printf("%s %s %s\n", badge, pkg, test)
}

func (c *TestCommand) printResult(action *TestAction, fgColor, bgColor, badgeText string) {
	badge := c.createBadge(badgeText, bgColor, fgColor)
	pkg := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorGray)).
		Render(action.Package)
	test := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorWhite)).
		Render(action.Test)
	duration := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorGray)).
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
	if !withBadge {
		fmt.Printf("       %s\n", text)
		return
	}
	badge := c.createBadge(" ERR  ", ColorRed, ColorWhite)
	fmt.Printf("%s %s\n", badge, text)
}

func (c *TestCommand) printSummary(success bool) {
	result := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorRed)).Render("TESTS FAILED")
	if success {
		result = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorGreen)).Render("TESTS PASSED")
	}

	rows := [][]string{
		{lipgloss.NewStyle().Foreground(lipgloss.Color(ColorBlue)).Render("Total Run     "), strconv.Itoa(c.totalRun)},
		{lipgloss.NewStyle().Foreground(lipgloss.Color(ColorGreen)).Render("Total Pass"), strconv.Itoa(c.totalPass)},
		{lipgloss.NewStyle().Foreground(lipgloss.Color(ColorRed)).Render("Total Fail"), strconv.Itoa(c.totalFail)},
		{lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPurple)).Render("Total Skip"), strconv.Itoa(c.totalSkip)},
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color(ColorGray))).
		Headers("ACTION", "COUNT").
		Rows(rows...)

	fmt.Printf("\n%s\n⌛ %.2fs\n\nTEST SUMMARY\n%s\n",
		result,
		time.Since(c.startTime).Seconds(),
		t.Render(),
	)
}

func (c *TestCommand) createBadge(text, bgColor, fgColor string) string {
	badgeText := lipgloss.NewStyle().
		Background(lipgloss.Color(bgColor)).
		Foreground(lipgloss.Color(fgColor)).
		Render(text)
	badgeArrow := lipgloss.NewStyle().
		Foreground(lipgloss.Color(bgColor)).
		Render("\uE0B0")
	return badgeText + badgeArrow
}
