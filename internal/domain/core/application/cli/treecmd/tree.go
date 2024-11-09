package treecmd

import (
	"fmt"
	"github.com/ksckaan1/hexago/internal/port"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type Commander interface {
	Command() *cobra.Command
}

type TreeCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
	tuilog   *tuilog.TUILog
}

func NewTreeCommand(i *do.Injector) (*TreeCommand, error) {
	return &TreeCommand{
		cmd: &cobra.Command{
			Use:     "tree",
			Example: "hexago tree",
			Short:   "Project structure tree",
			Long:    `Project structure tree`,
		},
		injector: i,
	}, nil
}

func (c *TreeCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *TreeCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *TreeCommand) init() {
	c.cmd.RunE = c.runner
}

func (c *TreeCommand) runner(cmd *cobra.Command, _ []string) error {
	projectService, err := do.Invoke[port.ProjectService](c.injector)
	if err != nil {
		return fmt.Errorf("invoke project service: %w", err)
	}

	moduleName, err := projectService.GetModuleName("go.mod")
	if err != nil {
		return fmt.Errorf("project service: get module name: %w", err)
	}

	entryPoints, err := projectService.GetAllEntryPoints(cmd.Context())
	if err != nil {
		return fmt.Errorf("project service: get all entry points: %w", err)
	}

	domains, err := projectService.GetAllDomains(cmd.Context())
	if err != nil {
		return fmt.Errorf("project service: get all domains: %w", err)
	}

	domainTree := make([]any, len(domains))

	for i := range domains {
		services, err := projectService.GetAllServices(cmd.Context(), domains[i])
		if err != nil {
			return fmt.Errorf("project service: get all services: %w", err)
		}

		apps, err := projectService.GetAllApplications(cmd.Context(), domains[i])
		if err != nil {
			return fmt.Errorf("project service: get all applications: %w", err)
		}

		domainTree = append(domainTree, tree.Root(domains[i]).Child(
			tree.Root(c.title("Services", len(services))).
				Child(c.colorizeElements(services)),
			tree.Root(c.title("Applications", len(apps))).
				Child(c.colorizeElements(apps)),
		))
	}

	globalPackages, err := projectService.GetAllPackages(cmd.Context(), true)
	if err != nil {
		return fmt.Errorf("project service: get all global packages: %w", err)
	}

	internalPackages, err := projectService.GetAllPackages(cmd.Context(), false)
	if err != nil {
		return fmt.Errorf("project service: get all internal packages: %w", err)
	}

	infras, err := projectService.GetAllInfrastructures(cmd.Context())
	if err != nil {
		return fmt.Errorf("project service: get all infrastructures: %w", err)
	}

	ports, err := projectService.GetAllPorts(cmd.Context())
	if err != nil {
		return fmt.Errorf("project service: get all ports: %w", err)
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
