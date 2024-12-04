package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
)

func Arabic(_ list.Items, i int) string {
	return fmt.Sprintf("  %d.", i+1)
}

func Hrule(color string, width int) string {
	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color))
	return separator.Render(strings.Repeat("─", width))
}
