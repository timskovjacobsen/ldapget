package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"golang.org/x/term"
)

func Arabic(_ list.Items, i int) string {
	return fmt.Sprintf("  %d.", i+1)
}

func Hrule(color string) string {
	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color))
	width, _, _ := term.GetSize(int(os.Stderr.Fd()))
	return separator.Render(strings.Repeat("─", width-50))
}
