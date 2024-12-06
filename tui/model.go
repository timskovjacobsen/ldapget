package tui

import (
	"os"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/timskovjacobsen/ldapget/client"
	"github.com/timskovjacobsen/ldapget/config"
	"golang.org/x/term"
)

// Represent viewing states of the TUI (mutually exclusive)
type ViewState int

const (
	ViewingGroups ViewState = iota
	ViewingGroupMembers
	SearchingGroups
	ViewingUsers
	ViewingUserGroups
	SearchingUsers
)

type Model struct {
	Config         *config.Config
	Tabs           []string
	ActiveTab      int
	Groups         []client.GroupInfo
	FilteredGroups []client.GroupInfo
	SelectedGroup  *client.GroupInfo
	Paginator      paginator.Model
	Cursor         int
	WindowSize     tea.WindowSizeMsg
	TUIState       ViewState
	SearchInput    string
	// StatusMsg      string
}

// Return true if the current TUI state is a type of searching state
func (m *Model) IsSearching() bool {
	if m.TUIState == SearchingGroups {
		return true
	}
	if m.TUIState == SearchingUsers {
		return true
	}
	return false
}

func NewModel(cfg *config.Config) *Model {

	w, h, _ := term.GetSize(int(os.Stderr.Fd()))
	tabs := []string{"Groups", "Users"}
	groups := client.Groups(cfg)
	var groupsContent []string
	for _, group := range groups {
		groupsContent = append(groupsContent, FormatGroup(group, w))
	}

	p := paginator.New()
	p.Type = paginator.Arabic
	p.PerPage = 5
	p.SetTotalPages(len(groups))

	return &Model{
		Config:     cfg,
		Tabs:       tabs,
		ActiveTab:  0,
		Groups:     groups,
		Paginator:  p,
		WindowSize: tea.WindowSizeMsg{Width: w, Height: h},
		TUIState:   ViewingGroups,
	}
}
