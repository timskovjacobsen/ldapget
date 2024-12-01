package tui

import (
	"strings"
)

func (m *Model) View() string {
	var b strings.Builder

	// Get the entries for the current page
	start, end := m.Paginator.GetSliceBounds(len(m.Groups))
	visibleGroups := m.Groups[start:end]

	// Render entries
	for i, group := range visibleGroups {
		currentIdx := i + m.Viewport.Top
		// group := FormatGroup(group)

		if currentIdx == m.Cursor {
			group = highlightStyle.Render(group)
		}

		b.WriteString(group)
		b.WriteString("\n")
	}

	// Page indicator and help
	b.WriteString("  " + m.Paginator.View() + "\n")
	b.WriteString("\n  h/l ←/→: change page • q: quit\n")

	return b.String()
}
