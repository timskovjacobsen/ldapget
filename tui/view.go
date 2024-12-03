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

	// doc := strings.Builder{}

	// var renderedTabs []string

	// for i, t := range m.Tabs {
	// 	var style lipgloss.Style
	// 	isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.ActiveTab
	// 	if isActive {
	// 		style = activeTabStyle
	// 	} else {
	// 		style = inactiveTabStyle
	// 	}
	// 	border, _, _, _, _ := style.GetBorder()
	// 	if isFirst && isActive {
	// 		border.BottomLeft = "│"
	// 	} else if isFirst && !isActive {
	// 		border.BottomLeft = "├"
	// 	} else if isLast && isActive {
	// 		border.BottomRight = "│"
	// 	} else if isLast && !isActive {
	// 		border.BottomRight = "┤"
	// 	}
	// 	style = style.Border(border)
	// 	renderedTabs = append(renderedTabs, style.Render(t))
	// }

	// row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	// contentWidth := w - 4
	// contentStyle := lipgloss.NewStyle().
	// 	BorderStyle(lipgloss.NormalBorder()).
	// 	BorderForeground(lipgloss.Color("240")).
	// 	Padding(1, 2)
	// renderedContent := contentStyle.Width(contentWidth).Render()

	// doc.WriteString(row)
	// doc.WriteString("\n")

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
			itemStyle := lipgloss.NewStyle()

			if i == m.Cursor {
				itemStyle = style.Highlight
			}
			content.WriteString(itemStyle.Render(FormatGroup(group)))
			content.WriteString("\n")
		}
		content.WriteString("  " + m.Paginator.View() + "\n")
		content.WriteString("\n  h/l ←/→: change page • q: quit\n")

		// Render content in content area
		b.WriteString(style.Content.Width(width).Height(height).Render(content.String()))
	}

	return style.Window.Width(w - 10).Height(h - 545).Render(b.String())
}
