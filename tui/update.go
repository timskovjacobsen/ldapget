package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/timskovjacobsen/ldapget/client"
	"github.com/timskovjacobsen/ldapget/config"
)

type fetchMembersMsg struct {
	Members []client.UserInfo
}

func fetchMembers(group *client.GroupInfo, cfg *config.Config) tea.Cmd {
	return func() tea.Msg {
		members, _ := client.GroupMembers(group.Name, cfg)
		return fetchMembersMsg{Members: members}
	}
}

func (m *Model) filterGroups() {
	if m.SearchInput == "" {
		m.FilteredGroups = m.Groups
	} else {
		m.FilteredGroups = nil // reset filtering slice to empty
		searchLower := strings.ToLower(m.SearchInput)
		for _, group := range m.Groups {
			if strings.Contains(strings.ToLower(group.Name), searchLower) {
				m.FilteredGroups = append(m.FilteredGroups, group)
			}
		}
	}
	// Update paginator with new filtered length
	m.Paginator.Page = 0
	if len(m.FilteredGroups) == 0 {
		m.Paginator.TotalPages = 1
	} else {
		m.Paginator.SetTotalPages(len(m.FilteredGroups))
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c": // quit action is common for all TUI views
			return m, tea.Quit
		case "q":
			if !m.IsSearching() { // Typing "q" in interactive search shouldn't quit
				return m, tea.Quit
			}
		// Entering search must be handled first, so we're not rendering the "/"
		case "/":
			if m.TUIState == ViewingGroups {
				m.TUIState = SearchingGroups
				m.SearchInput = ""
				m.filterGroups()
				return m, nil
			}
			if m.TUIState == ViewingUsers {
				m.TUIState = SearchingUsers
				m.SearchInput = ""
				// m.filterUsers()
				return m, nil
			}
		}
		// Apply controls based on what TUI view is currently active
		if m.TUIState == ViewingGroups {
			m.SetGroupsViewControls(msg)
		} else if m.TUIState == ViewingGroupMembers {
			m.SetMemberViewControls(msg)
		} else if m.IsSearching() {
			m.SetSearchControls(msg)
		}
	}

	m.Paginator, cmd = m.Paginator.Update(msg)
	return m, cmd
}
