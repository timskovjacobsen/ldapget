package tui

import (
	"fmt"
	"github.com/timskovjacobsen/ldapget/client"
	"strings"
)

func (m *Model) ensureVisible() {
	// If cursor is above viewport, scroll up
	if m.Cursor < m.Viewport.Top {
		m.Viewport.Top = m.Cursor
		return
	}

	// If cursor is below viewport, scroll down
	if m.Cursor >= m.Viewport.Top+m.Viewport.Height {
		m.Viewport.Top = m.Cursor - m.Viewport.Height + 1
		return
	}
}

func (m *Model) search() {
	term := strings.ToLower(m.SearchTerm)
	for i, group := range m.Groups {
		// NOTE: only searches in group names for now!
		if strings.Contains(strings.ToLower(group.Name), term) {
			m.Cursor = i
			m.ensureVisible()
			m.StatusMsg = ""
			return
		}
	}
	m.StatusMsg = "Not Found: " + m.SearchTerm
}

func FormatGroup(group client.GroupInfo) string {
	return fmt.Sprintf(`%s
🗺️ %s
🏷️ %s group
🎯 %s scope
📝 %s
👥 %d members
%s
`,
		group.Name,
		group.DN,
		group.Type,
		group.Scope,
		group.Description,
		group.Members,
		Hrule(),
	)
}
