package tui

import "github.com/timskovjacobsen/ldapget/style"

var (
	GroupsViewControls = style.Controls.Render("\n" +
		"  • ↑/j: up                      • ↓/k: down\n" +
		"  • ←/h: previous page           • →/l: next page\n" +
		"  • ctrl-←/ctrl-h: previous tab  • ctrl-→/ctrl-l: next tab\n" +
		"  • enter: see group members     • /: search\n" +
		"  • q: quit\n")

	// Remove chars valid for searching typing as controls
	SearchViewControls = style.Controls.Render("\n" +
		"  • ↑: up                      • ↓: down\n" +
		"  • ←: previous page           • →: next page\n" +
		"  • ctrl-←/ctrl-h: previous tab  • ctrl-→/ctrl-l: next tab\n" +
		"  • enter: see group members     • esc: quit search\n" +
		"  • q: quit\n")
)
