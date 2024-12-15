package appcmd

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/ksckaan1/hexago/internal/customerrors"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/ksckaan1/hexago/internal/port"
)

var _ port.Commander = (*AppLSCommand)(nil)

type AppLSCommand struct {
	cmd            *cobra.Command
	tuilog         *tuilog.TUILog
	projectService ProjectService

	// flags
	flagLine   *bool
	flagDomain *string
}

func NewAppLSCommand(projectService ProjectService, tl *tuilog.TUILog) (*AppLSCommand, error) {
	return &AppLSCommand{
		cmd: &cobra.Command{
			Use:     "ls",
			Example: "hexago app ls -d core\nhexago app ls (select domain interatively)",
			Short:   "List Applications",
			Long:    `List Applications`,
		},
		projectService: projectService,
		tuilog:         tl,
	}, nil
}

func (c *AppLSCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *AppLSCommand) AddSubCommand(cmd port.Commander) {
	c.cmd.AddCommand(cmd.Command())
}

func (c *AppLSCommand) init() {
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := c.runner(cmd, args)
		if err != nil {
			return customerrors.ErrSuppressed
		}
		return nil
	}
	c.flagLine = c.cmd.Flags().BoolP("line", "l", false, "hexago app ls -l")
	c.flagDomain = c.cmd.Flags().StringP("domain", "d", "", "hexago app ls -d <domainname>")
}

func (c *AppLSCommand) runner(cmd *cobra.Command, _ []string) error {
	domains, err := c.projectService.GetAllDomains(cmd.Context())
	if err != nil {

		c.tuilog.Error(err.Error())

		return fmt.Errorf("projectService.GetAllDomains: %w", err)
	}

	if len(domains) == 0 {

		c.tuilog.Error("No domains found.\nA domain needs to be created first")

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

				c.tuilog.Error("Select a domain: " + err2.Error())

				return fmt.Errorf("select a domain: %w", err2)
			}
		}
	} else if !slices.Contains(domains, *c.flagDomain) {

		c.tuilog.Error("Domain not found: " + *c.flagDomain)

		return fmt.Errorf("domain not found: %s", *c.flagDomain)
	}

	allApps := make([]string, 0)

	for i := range domains {
		if !(*c.flagDomain == "*" || domains[i] == *c.flagDomain) {
			continue
		}

		apps, err2 := c.projectService.GetAllApplications(cmd.Context(), domains[i])
		if err2 != nil {

			c.tuilog.Error(err2.Error())

			return fmt.Errorf("projectService.GetAllApplications: %w", err2)
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
