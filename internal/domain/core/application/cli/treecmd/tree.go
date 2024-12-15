package treecmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/internal/customerrors"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/ksckaan1/hexago/internal/port"
)

var _ port.Commander = (*TreeCommand)(nil)

type TreeCommand struct {
	cmd            *cobra.Command
	projectService ProjectService
	tuilog         *tuilog.TUILog
}

func NewTreeCommand(projectService ProjectService, tl *tuilog.TUILog) (*TreeCommand, error) {
	return &TreeCommand{
		cmd: &cobra.Command{
			Use:     "tree",
			Example: "hexago tree",
			Short:   "Project structure tree",
			Long:    `Project structure tree`,
		},
		projectService: projectService,
		tuilog:         tl,
	}, nil
}

func (c *TreeCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *TreeCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}

func (c *TreeCommand) init() {
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := c.runner(cmd, args)
		if err != nil {
			return customerrors.ErrSuppressed
		}
		return nil
	}
}

func (c *TreeCommand) runner(cmd *cobra.Command, _ []string) error {
	moduleName, err := c.projectService.GetModuleName("go.mod")
	if err != nil {
		return fmt.Errorf("projectService.GetModuleName: %w", err)
	}

	entryPoints, err := c.projectService.GetAllEntryPoints(cmd.Context())
	if err != nil {
		return fmt.Errorf("projectService.GetAllEntryPoints: %w", err)
	}

	domains, err := c.projectService.GetAllDomains(cmd.Context())
	if err != nil {
		return fmt.Errorf("projectService.GetAllDomains: %w", err)
	}

	domainTree := make([]any, len(domains))

	for i := range domains {
		services, err2 := c.projectService.GetAllServices(cmd.Context(), domains[i])
		if err2 != nil {
			return fmt.Errorf("projectService.GetAllServices: %w", err2)
		}

		apps, err2 := c.projectService.GetAllApplications(cmd.Context(), domains[i])
		if err2 != nil {
			return fmt.Errorf("projectService.GetAllApplications: %w", err2)
		}

		domainTree = append(domainTree, tree.Root(domains[i]).Child(
			tree.Root(c.title("Services", len(services))).
				Child(c.colorizeElements(services)),
			tree.Root(c.title("Applications", len(apps))).
				Child(c.colorizeElements(apps)),
		))
	}

	globalPackages, err := c.projectService.GetAllPackages(cmd.Context(), true)
	if err != nil {
		return fmt.Errorf("projectService.GetAllPackages: %w", err)
	}

	internalPackages, err := c.projectService.GetAllPackages(cmd.Context(), false)
	if err != nil {
		return fmt.Errorf("projectService.GetAllPackages: %w", err)
	}

	infras, err := c.projectService.GetAllInfrastructures(cmd.Context())
	if err != nil {
		return fmt.Errorf("projectService.GetAllInfrastructures: %w", err)
	}

	ports, err := c.projectService.GetAllPorts(cmd.Context())
	if err != nil {
		return fmt.Errorf("projectService.GetAllPorts: %w", err)
	}

	treePresentation := tree.Root(fmt.Sprintf("Project (%s)", moduleName)).Child(
		tree.Root(c.title("Entry Points", len(entryPoints))).
			Child(c.colorizeElements(entryPoints)),
		tree.Root(c.title("Domains", len(domains))).
			Child(domainTree...),
		tree.Root(c.title("Infrastructures", len(infras))).
			Child(c.colorizeElements(infras)),
		tree.Root(c.title("Packages", len(globalPackages)+len(internalPackages))).
			Child(
				tree.Root(c.title("Global", len(globalPackages))).
					Child(c.colorizeElements(globalPackages)),
				tree.Root(c.title("Internal", len(internalPackages))).
					Child(c.colorizeElements(internalPackages)),
			),
		tree.Root(c.title("Ports", len(ports))).
			Child(c.colorizeElements(ports)),
	).String()

	fmt.Println(treePresentation)

	return nil
}

func (c *TreeCommand) colorizeElements(elems []string) []string {
	renderer := lipgloss.NewStyle().Foreground(lipgloss.Color("#7571F9"))
	return lo.Map(elems, func(item string, index int) string {
		return renderer.Render(item)
	})
}

const title = "%s %s"

func (c *TreeCommand) title(label string, count int) string {
	return fmt.Sprintf(title, label, lipgloss.NewStyle().Foreground(lipgloss.Color("#485058")).Render(fmt.Sprintf("(%d)", count)))
}
