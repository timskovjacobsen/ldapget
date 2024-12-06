package tui

import (
	"os"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/timskovjacobsen/ldapget/client"
	"github.com/timskovjacobsen/ldapget/config"
	"golang.org/x/term"
)

type Model struct {
	Config         *config.Config
	Tabs           []string
	TabContent     [][]string
	ActiveTab      int
	Groups         []client.GroupInfo
	Paginator      paginator.Model
	Cursor         int
	WindowSize     tea.WindowSizeMsg
	StatusMsg      string
	SelectedGroup  *client.GroupInfo
	ViewingMembers bool
	ViewingGroups  bool
	IsSearching    bool
	FilteredGroups []client.GroupInfo
	SearchInput    string
}

func NewModel(cfg *config.Config) *Model {

	w, h, _ := term.GetSize(int(os.Stderr.Fd()))
	tabs := []string{"Groups", "Users"}
	groups := client.Groups(cfg)
	var groupsContent []string
	for _, group := range groups {
		groupsContent = append(groupsContent, FormatGroup(group, w))
	}
	tabContent := [][]string{groupsContent, {"Users"}}

	p := paginator.New()
	p.Type = paginator.Arabic
	p.PerPage = 5
	p.SetTotalPages(len(groups))

	return &Model{
		Config:        cfg,
		Tabs:          tabs,
		TabContent:    tabContent,
		ActiveTab:     0,
		Groups:        groups,
		Paginator:     p,
		WindowSize:    tea.WindowSizeMsg{Width: w, Height: h},
		ViewingGroups: true,
	}
}
