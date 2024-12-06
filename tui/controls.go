package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/timskovjacobsen/ldapget/style"
)

// The controls and navigation functionality below is separated into what
// view they are associated with.
// This results in some code duplication, but increases readability, becuase
// it's easier to see what code is run when each TUI view is active.

var (
	GroupsViewControls = style.Controls.Render("\n" +
		"  • ↑/j: up                      • ↓/k: down\n" +
		"  • ←/h: previous page           • →/l: next page\n" +
		"  • ctrl-←/ctrl-h: previous tab  • ctrl-→/ctrl-l: next tab\n" +
		"  • enter: see group members     • /: search\n" +
		"  • q: quit\n")

	// Note: no control chars that might be used in search requests
	SearchViewControls = style.Controls.Render("\n" +
		"  • ↑: up                      • ↓: down\n" +
		"  • ←: previous page           • →: next page\n" +
		"  • ctrl-←/ctrl-h: previous tab  • ctrl-→/ctrl-l: next tab\n" +
		"  • enter: see group members     • esc: quit search\n" +
		"  • ctrl-c: quit\n")

	GroupMembersViewControls = style.Controls.Render("\n" +
		"  • b/esc: back to group view    • q: quit\n")
)

// Apply controls and functionality for groups list view
func (m *Model) SetGroupsViewControls(msg tea.KeyMsg) {
	switch msg.String() {
	case "q":
		tea.Quit()
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}
	case "down", "j":
		start, end := m.Paginator.GetSliceBounds(len(m.Groups))
		currentPageEntries := m.Groups[start:end]
		if m.Cursor < len(currentPageEntries)-1 {
			m.Cursor++
		}
	case "right", "l":
		if m.Paginator.Page != m.Paginator.TotalPages-1 {
			m.Cursor = 0 // Reset cursor when changing pages
			m.Paginator.NextPage()
		}
	case "left", "h":
		if m.Paginator.Page != 0 {
			m.Cursor = 0 // Reset cursor when changing pages
			m.Paginator.PrevPage()
		}
	case "ctrl+left", "ctrl+h":
		m.ActiveTab = max(m.ActiveTab-1, 0)
	case "ctrl+right", "ctrl+l":
		m.ActiveTab = min(m.ActiveTab+1, len(m.Tabs)-1)
	case "enter":
		if m.ViewingGroups && len(m.Groups) > 0 {
			start, end := m.Paginator.GetSliceBounds(len(m.Groups))
			visibleGroups := m.Groups[start:end]
			if m.Cursor < len(visibleGroups) {
				m.SelectedGroup = &visibleGroups[m.Cursor]
				m.ViewingMembers = true
				m.ViewingGroups = false
			}
		}
	}
}

// Apply controls and functionality for search filter view
func (m *Model) SetSearchControls(msg tea.KeyMsg) {

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
		m.Cursor = 0

	case tea.KeySpace: // Space key must be handled separately (not in KeyRunes)
		m.SearchInput += " "
		m.filterGroups()

	case tea.KeyEsc: // Reset search filter and return to previous view
		m.IsSearching = false
		m.ViewingGroups = true
		m.ViewingMembers = false
		m.SearchInput = ""
		m.FilteredGroups = m.Groups
		m.Paginator.SetTotalPages(len(m.FilteredGroups))

	}
	switch msg.String() {
	case "up":
		if m.Cursor > 0 {
			m.Cursor--
		}
	case "down":
		start, end := m.Paginator.GetSliceBounds(len(m.Groups))
		currentPageEntries := m.Groups[start:end]
		if m.Cursor < len(currentPageEntries)-1 {
			m.Cursor++
		}
	case "right":
		if m.Paginator.Page != m.Paginator.TotalPages-1 {
			m.Cursor = 0 // Reset cursor when changing pages
			m.Paginator.NextPage()
		}
	case "left":
		if m.Paginator.Page != 0 {
			m.Cursor = 0 // Reset cursor when changing pages
			m.Paginator.PrevPage()
		}
	case "ctrl+left", "ctrl+h":
		m.ActiveTab = max(m.ActiveTab-1, 0)
	case "ctrl+right", "ctrl+l":
		m.ActiveTab = min(m.ActiveTab+1, len(m.Tabs)-1)
	case "enter":
		if len(m.FilteredGroups) > 0 {
			start, end := m.Paginator.GetSliceBounds(len(m.FilteredGroups))
			visibleGroups := m.FilteredGroups[start:end]
			if m.Cursor < len(visibleGroups) {
				m.SelectedGroup = &visibleGroups[m.Cursor]
				m.ViewingMembers = true
				m.ViewingGroups = false
				m.IsSearching = false
			}
		}
	}
}

// Apply controls and functionality for viewing members of a group
func (m *Model) SetMemberViewControls(msg tea.KeyMsg) {
	switch msg.String() {
	case "b", "esc":
		m.ViewingGroups = true
		m.ViewingMembers = false
		m.IsSearching = false

		// We clear the filter when returning, to simplify things. In the future we
		// could revert to the search, if that was active before
		m.FilteredGroups = m.Groups
		m.Paginator.SetTotalPages(len(m.Groups))
	}
}
