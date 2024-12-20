package tuilog

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type logType uint

const (
	infoType logType = iota
	successType
	warningType
	errorType
)

func (*TUILog) log(lt logType, msg string, title ...string) {
	var (
		logColor string
		logTitle string
		logIcon  string
	)
	switch lt {
	case successType:
		logColor = "#32CD32"
		logTitle = "Success"
		logIcon = "✓"
	case warningType:
		logColor = "#FFA500"
		logTitle = "Warning"
		logIcon = "�"
	case errorType:
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

func (*TUILog) logWithoutTitle(lt logType, msg string) {
	var (
		logColor string
	)
	switch lt {
	case successType:
		logColor = "#32CD32"
	case warningType:
		logColor = "#FFA500"
	case errorType:
		logColor = "#FF0000"
	default:
		logColor = "#87CEEB"
	}

	op := lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(logColor)).
		BorderLeft(true).
		PaddingLeft(1).
		Render(msg)

	fmt.Println(op)
}
