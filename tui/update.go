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
		m.FilteredGroups = nil
		searchLower := strings.ToLower(m.SearchInput)
		for _, group := range m.Groups {
			if strings.Contains(strings.ToLower(group.Name), searchLower) {
				m.FilteredGroups = append(m.FilteredGroups, group)
			}
		}
	}
	// Update paginator with new filtered length
	m.Paginator.SetTotalPages(len(m.FilteredGroups))
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case fetchMembersMsg:
		m.GroupMembers = msg.Members
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
			return m, nil

		case "down", "j":
			start, end := m.Paginator.GetSliceBounds(len(m.Groups))
			currentPageEntries := m.Groups[start:end]
			if m.Cursor < len(currentPageEntries)-1 {
				m.Cursor++
			}
			return m, nil

		case "right", "l":
			if m.Paginator.Page != m.Paginator.TotalPages-1 {
				m.Cursor = 0 // Reset cursor when changing pages
				m.Paginator.NextPage()
			}
			return m, nil

		case "left", "h":
			if m.Paginator.Page != 0 {
				m.Cursor = 0 // Reset cursor when changing pages
				m.Paginator.PrevPage()
			}
			return m, nil

		case "ctrl+left", "ctrl+h":
			m.ActiveTab = max(m.ActiveTab-1, 0)
			return m, nil
		case "ctrl+right", "ctrl+l":
			m.ActiveTab = min(m.ActiveTab+1, len(m.Tabs)-1)
			return m, nil
		case "/":
			if !m.IsSearching {
				m.IsSearching = true
				m.SearchInput = ""
				m.filterGroups()
			}
			return m, nil
		case "b", "esc":
			if m.ViewingMembers {
				// Return to group list
				m.ViewingMembers = false
				// m.SelectedGroup = nil
				m.ViewingGroups = true
				return m, nil
			}
			if m.IsSearching {
				// Reset filter
				m.IsSearching = false
				m.SearchInput = ""
				m.FilteredGroups = m.Groups
				m.Paginator.SetTotalPages(len(m.FilteredGroups))
			}
		case "enter":
			if m.ViewingGroups && len(m.Groups) > 0 {
				start, end := m.Paginator.GetSliceBounds(len(m.Groups))
				visibleGroups := m.Groups[start:end]
				if m.Cursor < len(visibleGroups) {
					m.SelectedGroup = &visibleGroups[m.Cursor]
					m.ViewingMembers = true
					m.ViewingGroups = false
					return m, fetchMembers(m.SelectedGroup, m.Config)
				}
			}
		}
		if m.IsSearching {
			switch msg.Type {
			case tea.KeyBackspace:
				// User is deleting a char from the search input
				if len(m.SearchInput) > 0 {
					m.SearchInput = m.SearchInput[:len(m.SearchInput)-1]
					m.filterGroups()
				}
			case tea.KeyRunes:
				// User is typing a char into the search input
				m.SearchInput += string(msg.Runes)
				m.filterGroups()

			case tea.KeySpace: // Space key must be handled separately
				m.SearchInput += " "
				m.filterGroups()
			}
		}
	}

	m.Paginator, cmd = m.Paginator.Update(msg)
	return m, cmd
}
