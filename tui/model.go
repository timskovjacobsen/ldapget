package tui

import (
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Groups     []string
	Viewport   viewport
	Paginator  paginator.Model
	Cursor     int
	WindowSize tea.WindowSizeMsg
	SearchMode bool
	SearchTerm string
	StatusMsg  string
}

type viewport struct {
	Top    int
	Height int
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("87"))

	// separatorStyle = lipgloss.NewStyle().
	// Foreground(lipgloss.Color("240"))

	highlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("137")).
			Bold(true)
)

func NewModel(groups []string) *Model {
	p := paginator.New()
	p.Type = paginator.Arabic // Using numbers instead of dots since we have more content
	p.PerPage = 6             // Show 5 LDAP entries per page
	p.SetTotalPages(len(groups))
	return &Model{
		Groups:    groups,
		Paginator: p,
	}
}
