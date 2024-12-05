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

const (
	GREY = "#545454"
)

var (
	ItemTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#cbba82"))
	ActiveItem = lipgloss.NewStyle().
			Border(lipgloss.Border{Left: "│"}).
			BorderForeground(lipgloss.Color("99")).
			Align(lipgloss.Left).
			Bold(true)
	InactiveItem = lipgloss.NewStyle().
			Border(lipgloss.Border{Left: " "}).
			BorderForeground(lipgloss.Color("1")). // there must be a color here in order for left-alignment to work, no idea why
			Align(lipgloss.Left)

	NotSet = lipgloss.NewStyle(). // for entries with no value, e.g. group description
		Italic(true).
		Foreground(lipgloss.Color(GREY))
	HighlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	InactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	ActiveTabBorder   = tabBorderWithBottom("┘", " ", "└")
	InactiveTab       = lipgloss.NewStyle().
				Border(InactiveTabBorder, true).
				BorderForeground(HighlightColor).
				Padding(0, 1)
	ActiveTab = InactiveTab.Border(ActiveTabBorder, true)
	Content   = lipgloss.NewStyle().
			Align(lipgloss.Left).
			UnsetAlign().
			Padding(1, 2)
	Window = lipgloss.NewStyle().
		Padding(0, 0).
		Align(lipgloss.Left).
		UnsetBorderTop()
	SecondaryText = lipgloss.NewStyle().
			Foreground(lipgloss.Color(GREY))
	Controls = lipgloss.NewStyle().
			Foreground(lipgloss.Color(GREY))
	Enumerate = lipgloss.NewStyle().
			Foreground(lipgloss.Color(GREY))
)
