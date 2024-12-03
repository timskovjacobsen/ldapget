package tui

import (
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/timskovjacobsen/ldapget/client"
	"github.com/timskovjacobsen/ldapget/config"
)

type Model struct {
	Tabs           []string
	TabContent     [][]string
	ActiveTab      int
	Groups         []client.GroupInfo
	Paginator      paginator.Model
	Cursor         int
	WindowSize     tea.WindowSizeMsg
	SearchMode     bool
	SearchTerm     string
	StatusMsg      string
	GroupMembers   []client.UserInfo
	ViewingMembers bool
}

type fetchMembersMsg struct {
	Members []client.UserInfo
}

func fetchMembers(group client.GroupInfo, cfg *config.Config) tea.Cmd {
	return func() tea.Msg {
		members, _ := client.GroupMembers(group.Name, cfg)
		return fetchMembersMsg{Members: members}
	}
}

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
