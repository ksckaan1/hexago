package cli

import (
	_ "embed"

	"github.com/charmbracelet/lipgloss"
)

//go:embed assets/figlet.txt
var figletFile string

//go:embed assets/logo.txt
var logoFile string

var header = lipgloss.JoinHorizontal(lipgloss.Top, logoFile, " ", figletFile)
