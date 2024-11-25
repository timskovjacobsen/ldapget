package cmd

import (
	"fmt"
	"github.com/charmbracelet/lipgloss/list"
)

func Arabic(_ list.Items, i int) string {
	return fmt.Sprintf("  %d.", i+1)
}
