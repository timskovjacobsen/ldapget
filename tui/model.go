package tui

import (
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/timskovjacobsen/ldapget/client"
	"github.com/timskovjacobsen/ldapget/config"
)

type Model struct {
	Tabs       []string
	TabContent [][]string
	ActiveTab  int
	Groups     []client.GroupInfo
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
	highlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("137")).
			Bold(true)
)

func NewModel(cfg *config.Config) *Model {

	tabs := []string{"Groups", "Users"}
	groups := client.Groups(cfg)
	var groupsContent []string
	for _, group := range groups {
		groupsContent = append(groupsContent, FormatGroup(group))
	}
	tabContent := [][]string{groupsContent, {"Users"}}

	p := paginator.New()
	p.Type = paginator.Arabic
	p.PerPage = 5
	p.SetTotalPages(len(groups))

	return &Model{
		Tabs:       tabs,
		TabContent: tabContent,
		ActiveTab:  0,
		Groups:     groups,
		Paginator:  p,
	}
}
