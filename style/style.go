package style

import (
	"github.com/charmbracelet/lipgloss"
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	ItemTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#cbba82"))
	Highlight = lipgloss.NewStyle().
			Foreground(lipgloss.Color("137")).
			Bold(true)
	NotSet = lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#3C3C3C"))
	HighlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	InactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	ActiveTabBorder   = tabBorderWithBottom("┘", " ", "└")
	InactiveTab       = lipgloss.NewStyle().Border(InactiveTabBorder, true).BorderForeground(HighlightColor).Padding(0, 1)
	ActiveTab         = InactiveTab.Border(ActiveTabBorder, true)
	Content           = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240")).
				Padding(1, 2)
	Window = lipgloss.NewStyle().
		BorderForeground(HighlightColor).
		Padding(2, 0).
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder()).
		UnsetBorderTop()
)
