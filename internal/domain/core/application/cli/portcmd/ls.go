package portcmd

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

type PortLSCommand struct {
	cmd      *cobra.Command
	injector *do.Injector
	tuilog   *tuilog.TUILog

	// flags
	flagLine   *bool
	flagDomain *string
}

func NewPortLSCommand(i *do.Injector) (*PortLSCommand, error) {
	return &PortLSCommand{
		cmd: &cobra.Command{
			Use:     "ls",
			Example: "hexago port ls -d <domainname>\nhexago port ls (select domain interactively)",
			Short:   "List ports",
			Long:    `List ports`,
		},
		injector: i,
		tuilog:   do.MustInvoke[*tuilog.TUILog](i),
		// flags
		flagLine:   new(bool),
		flagDomain: new(string),
	}, nil
}

func (c *PortLSCommand) Command() *cobra.Command {
	c.init()
	return c.cmd
}

func (c *PortLSCommand) AddCommand(cmds ...Commander) {
	c.cmd.AddCommand(lo.Map(cmds, func(cmd Commander, _ int) *cobra.Command {
		return cmd.Command()
	})...)
}

func (c *PortLSCommand) init() {
	c.cmd.RunE = c.runner
	c.flagLine = c.cmd.Flags().BoolP("line", "l", false, "hexago domain ls -l")
	c.flagDomain = c.cmd.Flags().StringP("domain", "d", "", "hexago service ls -d <domainname>")
}

func (c *PortLSCommand) runner(cmd *cobra.Command, _ []string) error {
	projectService, err := do.Invoke[port.ProjectService](c.injector)
	if err != nil {
		return fmt.Errorf("invoke project service: %w", err)
	}

	domains, err := projectService.GetAllDomains(cmd.Context())
	if err != nil {
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
				c.tuilog.Error("Select a domain: ", err2.Error())
				fmt.Println("")
				return fmt.Errorf("select a domain: %w", err2)
			}
		}
	} else if !slices.Contains(domains, *c.flagDomain) {
		fmt.Println("")
		c.tuilog.Error("Domain not found: ", *c.flagDomain)
		fmt.Println("")
		return fmt.Errorf("domain not found: %s", *c.flagDomain)
	}

	allPorts := make([]string, 0)

	for i := range domains {
		if !(*c.flagDomain == "*" || domains[i] == *c.flagDomain) {
			continue
		}

		ports, err2 := projectService.GetAllPorts(cmd.Context(), domains[i])
		if err2 != nil {
			fmt.Println("")
			c.tuilog.Error(err2.Error())
			fmt.Println("")
			return fmt.Errorf("project service: get all ports: %w", err2)
		}

		if *c.flagDomain == "*" {
			for j := range ports {
				ports[j] = domains[i] + ":" + ports[j]
			}
		}

		allPorts = append(allPorts, ports...)
	}

	seperator := lo.Ternary(*c.flagLine, "\n", " ")

	fmt.Println(strings.Join(allPorts, seperator))

	return nil
}
