package util

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type UILogType uint

const (
	Info UILogType = iota
	Success
	Warning
	Error
)

func UILog(lt UILogType, msg string, title ...string) {
	var (
		logColor string
		logTitle string
		logIcon  string
	)
	switch lt {
	case Success:
		logColor = "#32CD32"
		logTitle = "Success"
		logIcon = "✓"
	case Warning:
		logColor = "#FFA500"
		logTitle = "Warning"
		logIcon = "�"
	case Error:
		logColor = "#FF0000"
		logTitle = "Error"
		logIcon = "✗"
	default:
		logColor = "#87CEEB"
		logTitle = "Info"
		logIcon = "ⓘ"
	}

	if len(title) > 0 {
		logTitle = title[0]
	}

	op := lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(logColor)).
		BorderLeft(true).
		PaddingLeft(1).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(logColor)).
				Render(logIcon+" "+logTitle),
			msg,
		))

	fmt.Println(op)
}
