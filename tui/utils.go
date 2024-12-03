package tui

import (
	"fmt"
	"strings"

	"github.com/timskovjacobsen/ldapget/client"
	"github.com/timskovjacobsen/ldapget/style"
)

func (m *Model) search() {
	term := strings.ToLower(m.SearchTerm)
	for i, group := range m.Groups {
		// NOTE: only searches in group names for now!
		if strings.Contains(strings.ToLower(group.Name), term) {
			m.Cursor = i
			m.StatusMsg = ""
			return
		}
	}
	m.StatusMsg = "Not Found: " + m.SearchTerm
}

func truncateOrPad(s string, width int) string {
	if len(s) > width {
		return s[:width-3] + ">" // continuation char
	}
	// Pad with spaces for strings that do not take up all width
	return s + strings.Repeat(" ", width-len(s))
}

// wrap text at word boundaries to a specific width
func wordWrap(text string, width int, indent int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var lines []string
	currentLine := words[0]

	for _, word := range words[1:] {
		// Check if adding the next word exceeds the width
		if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			// Start a new line
			lines = append(lines, currentLine)
			currentLine = strings.Repeat(" ", indent) + word
		}
	}
	lines = append(lines, currentLine)

	// Pad each line to maintain consistent width
	for i, line := range lines {
		lines[i] = truncateOrPad(line, width)
	}

	return strings.Join(lines, "\n")
}

func GroupHeight(group client.GroupInfo, maxWidth int) int {
	descrLines := (len(group.Description) + maxWidth - 1) / maxWidth
	return 3 + descrLines
}

func FormatGroup(group client.GroupInfo) string {
	var b strings.Builder

	width := 80
	membersTxt := fmt.Sprintf("👥 %d members", group.MemberCount)

	padCount := width - len(group.Name) - len(membersTxt)
	padding := strings.Repeat(" ", padCount)

	groupNameTxt := style.ItemTitle.Render(group.Name)
	b.WriteString(fmt.Sprintf("%s %s %s\n", groupNameTxt, padding, membersTxt))
	wrappedDN := wordWrap(group.DN, width, 3)
	b.WriteString(fmt.Sprintf("🗺️ %s\n", wrappedDN))
	var description string
	if len(group.Description) == 0 {
		description = style.NotSet.Render("no description")
	} else {
		description = wordWrap(group.Description, width, 3)
	}
	b.WriteString(fmt.Sprintf("📝 %s\n", description))
	b.WriteString(fmt.Sprintf("🏷️ %s group\n", group.Type))
	b.WriteString(fmt.Sprintf("🎯 %s scope\n", group.Scope))
	b.WriteString(Hrule("#555555"))
	return b.String()
}
