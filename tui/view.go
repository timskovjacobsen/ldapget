package tui

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/timskovjacobsen/ldapget/style"
	"golang.org/x/term"
)

var w, h, _ = term.GetSize(int(os.Stderr.Fd()))

func (m *Model) View() string {
	var b strings.Builder
	height := m.WindowSize.Height - 6
	width := m.WindowSize.Width - 4

	var tabs []string
	for i, tab := range m.Tabs {
		tabStyle := style.InactiveTab
		if i == m.ActiveTab {
			tabStyle = style.ActiveTab
		}
		tabs = append(tabs, tabStyle.Render(tab))
	}
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, tabs...))
	b.WriteString("\n")

	if m.ActiveTab == 0 {

		// Get the entries for the current page
		start, end := m.Paginator.GetSliceBounds(len(m.Groups))
		visibleGroups := m.Groups[start:end]

		// Render entries
		var content strings.Builder
		for i, group := range visibleGroups {
			var itemStyle = style.InactiveItem

			if i == m.Cursor {
				itemStyle = style.ActiveItem
			}
			content.WriteString(itemStyle.Render(FormatGroup(group, m.WindowSize.Width)))
			content.WriteString(Hrule("#555555", m.WindowSize.Width-16))
			content.WriteString("\n")
		}
		content.WriteString("  " + m.Paginator.View() + "\n")
		content.WriteString("\n  h/l ←/→: change page • q: quit\n")

		// Render content in content area
		b.WriteString(style.Content.Width(width).Height(height).UnsetAlign().Align(lipgloss.Left).Render(content.String()))
	}

	return style.Window.Width(w).Height(h).Render(b.String())
}
