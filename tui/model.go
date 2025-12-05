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
	Users          []client.UserInfo
	FilteredUsers  []client.UserInfo
	SelectedUser   *client.UserInfo
	Paginator      paginator.Model
	Cursor         int
	WindowSize     tea.WindowSizeMsg
	TUIState       ViewState
	SearchInput    string
	ErrorMsg       string
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

func itemsPerPage(itemHeight int) int {
	// Assmuing everything other than items, i.e. header, footer, etc.
	// takes up this many lines
	linesOther := 23
	_, h, _ := term.GetSize(int(os.Stderr.Fd()))
	linesAllItems := h - linesOther
	return linesAllItems / itemHeight
}

func NewModel(cfg *config.Config) *Model {

	w, h, _ := term.GetSize(int(os.Stderr.Fd()))
	tabs := []string{"Groups", "Users"}

	p := paginator.New()
	p.Type = paginator.Arabic
	p.PerPage = itemsPerPage(7)

	return &Model{
		Config:     cfg,
		Tabs:       tabs,
		ActiveTab:  0,
		Groups:     nil,
		Paginator:  p,
		WindowSize: tea.WindowSizeMsg{Width: w, Height: h},
		TUIState:   ViewingGroups,
	}
}
