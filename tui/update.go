package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
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
		}
	}

	m.Paginator, cmd = m.Paginator.Update(msg)
	return m, cmd
}
