package tui

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/timskovjacobsen/ldapget/client"
	"github.com/timskovjacobsen/ldapget/style"
)

// Format a list of group members
//
// E.g.
//  1. Axel Foley     (af@example.com)
//  2. Rowan Atkinson (ra@example.com)
func formatMemberList(members []client.UserInfo) string {
	if len(members) == 0 {
		return ""
	}
	// Find the maximum lengths needed for each column
	maxNameWidth := 0
	maxEmailWidth := 0
	maxNumberWidth := len(fmt.Sprintf("%d. ", len(members)))

	for _, member := range members {
		nameWidth := utf8.RuneCountInString(member.Name)
		emailWidth := utf8.RuneCountInString(member.Email)

		if nameWidth > maxNameWidth {
			maxNameWidth = nameWidth
		}
		if emailWidth > maxEmailWidth {
			maxEmailWidth = emailWidth
		}
	}

	// Build the formatted list
	var builder strings.Builder
	nameFormat := fmt.Sprintf("%%-%ds", maxNameWidth)
	numFormat := fmt.Sprintf("%%%dd.", maxNumberWidth)
	for i, member := range members {
		enu := style.Enumerate.Render(fmt.Sprintf(numFormat, i+1))
		fmt.Fprintf(&builder, "%s %s (%s)\n",
			enu,
			fmt.Sprintf(nameFormat, member.Name),
			member.Email)
	}
	return builder.String()
}

// Return a styled view of the members in the selected group
func (m *Model) renderMembersView() string {
	var b strings.Builder

	// Render header with group name
	header := fmt.Sprintf("Members of: %s", m.SelectedGroup.Name)
	b.WriteString(style.ItemTitle.Render(header))
	b.WriteString("\n\n")

	members, _ := client.GroupMembers(m.SelectedGroup.Name, m.Config)
	if len(members) == 0 {
		b.WriteString("No members found\n")
	} else {
		b.WriteString(formatMemberList(members))
	}

	// Add controls help
	b.WriteString(GroupMembersViewControls)

	// Wrap in the same content style as groups
	contentWidth := min(m.WindowSize.Width-4, m.WindowSize.Width+4)
	return style.Content.Width(contentWidth).Render(b.String())
}

// Return a styled view of the groups in the current filter
//
// Modifies the input builder and returns it in string form afterwards
func (m *Model) renderGroupsView(b *strings.Builder) string {
	height := m.WindowSize.Height - 6
	width := m.WindowSize.Width - 4
	var content strings.Builder
	var groupList []client.GroupInfo
	var controls string
	if m.TUIState == SearchingGroups {
		input := style.ItemTitle.Render(fmt.Sprintf("%s", m.SearchInput))
		content.WriteString(fmt.Sprintf("Search: %s_\n\n", input))
		groupList = m.FilteredGroups
		controls = SearchViewControls
	} else {
		content.WriteString("\n\n")
		groupList = m.Groups
		controls = GroupsViewControls
	}
	// Get the entries for the current page
	m.Paginator.SetTotalPages(len(groupList))
	start, end := m.Paginator.GetSliceBounds(len(groupList))

	// NOTE: There are some weird cases where [start:end] comes out something like
	// [10:0], which makes no sense. We'll handle it here to avoid a crash
	if start > end {
		start, end = 0, 0
	}
	visibleGroups := groupList[start:end]
	content.WriteString(style.SecondaryText.Render(fmt.Sprintf("Showing %d groups\n", len(groupList))))
	// Render entries
	for i, group := range visibleGroups {
		var itemStyle = style.InactiveItem

		if i == m.Cursor {
			itemStyle = style.ActiveItem
		}
		content.WriteString(itemStyle.Render(FormatGroup(group, m.WindowSize.Width)))
		content.WriteString(Hrule("#555555", m.WindowSize.Width-16))
		content.WriteString("\n")
	}
	content.WriteString("  " + m.Paginator.View() + "\n")
	content.WriteString(controls)

	// Render content in content area
	b.WriteString(style.Content.Width(width).Height(height).UnsetAlign().Align(lipgloss.Left).Render(content.String()))
	return b.String()
}

func (m *Model) renderUsersView(b *strings.Builder) string {
	height := m.WindowSize.Height - 6
	width := m.WindowSize.Width - 4
	var content strings.Builder
	var userList []client.UserInfo
	var controls string
	if m.TUIState == SearchingUsers {
		input := style.ItemTitle.Render(fmt.Sprintf("%s", m.SearchInput))
		content.WriteString(fmt.Sprintf("Search: %s_\n\n", input))
		userList = m.FilteredUsers
		controls = SearchViewControls
	} else {
		content.WriteString("\n\n")
		userList = m.Users
		controls = GroupsViewControls
	}
	// Get the entries for the current page
	m.Paginator.SetTotalPages(len(userList))
	start, end := m.Paginator.GetSliceBounds(len(userList))

	// NOTE: There are some weird cases where [start:end] comes out something like
	// [10:0], which makes no sense. We'll handle it here to avoid a crash
	if start > end {
		start, end = 0, 0
	}
	content.WriteString(style.SecondaryText.Render(fmt.Sprintf("Showing %d users\n", len(userList))))
	// Render visible entries given by pager
	for i, user := range userList[start:end] {
		var itemStyle = style.InactiveItem

		if i == m.Cursor {
			itemStyle = style.ActiveItem
		}
		content.WriteString(itemStyle.Render(user.Name))
		content.WriteString("\n")
		content.WriteString(Hrule("#555555", m.WindowSize.Width-16))
		content.WriteString("\n")
	}
	content.WriteString("  " + m.Paginator.View() + "\n")
	content.WriteString(controls)

	// Render content in content area
	b.WriteString(style.Content.Width(width).Height(height).UnsetAlign().Align(lipgloss.Left).Render(content.String()))
	return b.String()
}

func (m *Model) View() string {
	var b strings.Builder

	var tabs []string
	for i, tab := range m.Tabs {
		tabStyle := style.InactiveTab
		if i == m.ActiveTab {
			tabStyle = style.ActiveTab
		}
		tabs = append(tabs, tabStyle.Render(tab))
	}
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, tabs...))
	b.WriteString("\n")

	if m.ActiveTab == 0 { // Groups tab
		if m.Groups == nil {
			// The groups are loaded by defualt when the TUI starts, but we might have
			// cleared them from memory during the lifetime of the TUI, we'll fetch
			// them again
			m.Groups = client.Groups(m.Config)
		}
		// Render the appropriate view
		if m.TUIState == ViewingGroupMembers {
			return m.renderMembersView() // note: tabs are not shown here
		} else if m.TUIState == ViewingGroups {
			// We must be viewing the list of groups
			return m.renderGroupsView(&b)
		} else {
			return m.renderGroupsView(&b)
		}
	} else if m.ActiveTab == 1 { // Users tab
		if m.Users == nil {
			m.Users, _ = client.Users(m.Config)
		}
		// if m.TUIState == ViewingUsers {
		return m.renderUsersView(&b)
		// }
	}

	return style.Window.
		Width(m.WindowSize.Width).
		Height(m.WindowSize.Height).
		Render(b.String())
}
