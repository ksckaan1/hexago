package appcmd

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/ksckaan1/hexago/internal/domain/core/port"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type AppLSCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
	tuilog   *tuilog.TUILog

	// flags
	flagLine   *bool
	flagDomain *string
}

func NewAppLSCommand(i *do.Injector) (*AppLSCommand, error) {
	return &AppLSCommand{
		cmd: &cobra.Command{
			Use:     "ls",
			Example: "hexago app ls -d core\nhexago app ls (select domain interatively)",
			Short:   "List Applications",
			Long:    `List Applications`,
		},
		injector: i,
		tuilog:   do.MustInvoke[*tuilog.TUILog](i),
	}, nil
}

func (c *AppLSCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *AppLSCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *AppLSCommand) init() {
	c.cmd.RunE = c.runner
	c.flagLine = c.cmd.Flags().BoolP("line", "l", false, "hexago app ls -l")
	c.flagDomain = c.cmd.Flags().StringP("domain", "d", "", "hexago app ls -d <domainname>")
}

func (c *AppLSCommand) runner(cmd *cobra.Command, _ []string) error {
	projectService, err := do.Invoke[port.ProjectService](c.injector)
	if err != nil {
		return fmt.Errorf("invoke project service: %w", err)
	}

	domains, err := projectService.GetAllDomains(cmd.Context())
	if err != nil {
		fmt.Println("")
		c.tuilog.Error(err.Error())
		fmt.Println("")
		return fmt.Errorf("project service: get all domains: %w", err)
	}

	if len(domains) == 0 {
		fmt.Println("")
		c.tuilog.Error("No domains found.\nA domain needs to be created first")
		fmt.Println("")
		return fmt.Errorf("No domains found.\nA domain needs to be created first")
	}

	if *c.flagDomain == "" {
		if len(domains) == 1 {
			*c.flagDomain = domains[0]
		} else {

			selectList := []huh.Option[string]{
				huh.NewOption("* (All Domains)", "*"),
			}

			selectList = append(selectList, lo.Map(domains, func(d string, _ int) huh.Option[string] {
				return huh.NewOption(d, d)
			})...)

			err2 := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Select a domain.").
						Options(
							selectList...,
						).
						Value(c.flagDomain),
				).WithShowHelp(true),
			).Run()
			if err2 != nil {
				fmt.Println("")
				c.tuilog.Error("Select a domain: " + err2.Error())
				fmt.Println("")
				return fmt.Errorf("select a domain: %w", err2)
			}
		}
	} else if !slices.Contains(domains, *c.flagDomain) {
		fmt.Println("")
		c.tuilog.Error("Domain not found: " + *c.flagDomain)
		fmt.Println("")
		return fmt.Errorf("domain not found: %s", *c.flagDomain)
	}

	allApps := make([]string, 0)

	for i := range domains {
		if !(*c.flagDomain == "*" || domains[i] == *c.flagDomain) {
			continue
		}

		apps, err2 := projectService.GetAllApplications(cmd.Context(), domains[i])
		if err2 != nil {
			fmt.Println("")
			c.tuilog.Error(err.Error())
			fmt.Println("")
			return fmt.Errorf("project service: get all apps: %w", err2)
		}

		if *c.flagDomain == "*" {
			for j := range apps {
				apps[j] = domains[i] + ":" + apps[j]
			}
		}

		allApps = append(allApps, apps...)
	}

	seperator := lo.Ternary(*c.flagLine, "\n", " ")

	fmt.Println(strings.Join(allApps, seperator))

	return nil
}
