package tui

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	w, h, _           = term.GetSize(int(os.Stderr.Fd()))
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	contentStyle      = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240")).
				Padding(1, 2)
	windowStyle = lipgloss.NewStyle().
			BorderForeground(highlightColor).
			Padding(2, 0).
			Align(lipgloss.Center).
			Border(lipgloss.NormalBorder()).
			UnsetBorderTop()
)

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
		style := inactiveTabStyle
		if i == m.ActiveTab {
			style = activeTabStyle
		}
		tabs = append(tabs, style.Render(tab))
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
			style := lipgloss.NewStyle()

			if i == m.Cursor {
				style = highlightStyle
			}
			content.WriteString(style.Render(FormatGroup(group)))
			content.WriteString("\n")
		}
		content.WriteString("  " + m.Paginator.View() + "\n")
		content.WriteString("\n  h/l ←/→: change page • q: quit\n")

		// Render content in content area
		b.WriteString(contentStyle.Width(width).Height(height).Render(content.String()))
	}

	return windowStyle.Width(w - 10).Height(h - 545).Render(b.String())
}
