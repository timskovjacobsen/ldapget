package tui

import (
	"testing"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/timskovjacobsen/ldapget/client"
)

func TestFilterGroups(t *testing.T) {
	tests := []struct {
		name          string
		groups        []client.GroupInfo
		searchInput   string
		expectedCount int
		expectedNames []string
	}{
		{
			name: "empty search returns all groups",
			groups: []client.GroupInfo{
				{Name: "Admin"},
				{Name: "Users"},
				{Name: "Developers"},
			},
			searchInput:   "",
			expectedCount: 3,
			expectedNames: []string{"Admin", "Users", "Developers"},
		},
		{
			name: "case insensitive search",
			groups: []client.GroupInfo{
				{Name: "Admin"},
				{Name: "Users"},
				{Name: "Developers"},
			},
			searchInput:   "admin",
			expectedCount: 1,
			expectedNames: []string{"Admin"},
		},
		{
			name: "partial match",
			groups: []client.GroupInfo{
				{Name: "Admin"},
				{Name: "AdminGroup"},
				{Name: "Users"},
				{Name: "SuperAdmin"},
			},
			searchInput:   "admin",
			expectedCount: 3,
			expectedNames: []string{"Admin", "AdminGroup", "SuperAdmin"},
		},
		{
			name: "no matches",
			groups: []client.GroupInfo{
				{Name: "Admin"},
				{Name: "Users"},
			},
			searchInput:   "xyz",
			expectedCount: 0,
			expectedNames: []string{},
		},
		{
			name: "search with spaces",
			groups: []client.GroupInfo{
				{Name: "Domain Admins"},
				{Name: "Domain Users"},
				{Name: "LocalUsers"},
			},
			searchInput:   "domain",
			expectedCount: 2,
			expectedNames: []string{"Domain Admins", "Domain Users"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := paginator.New()
			p.PerPage = 10

			m := &Model{
				Groups:      tt.groups,
				SearchInput: tt.searchInput,
				Paginator:   p,
			}

			m.filterGroups()

			if len(m.FilteredGroups) != tt.expectedCount {
				t.Errorf("FilteredGroups count = %d, want %d", len(m.FilteredGroups), tt.expectedCount)
			}

			for i, expectedName := range tt.expectedNames {
				if i >= len(m.FilteredGroups) {
					t.Errorf("Expected group %s at index %d, but FilteredGroups is too short", expectedName, i)
					continue
				}
				if m.FilteredGroups[i].Name != expectedName {
					t.Errorf("FilteredGroups[%d].Name = %s, want %s", i, m.FilteredGroups[i].Name, expectedName)
				}
			}

			if tt.expectedCount > 0 {
				if m.Paginator.TotalPages == 0 {
					t.Error("Paginator.TotalPages should not be 0 when there are filtered results")
				}
			}

			if m.Paginator.Page != 0 {
				t.Errorf("Paginator.Page = %d, want 0 (should reset to first page)", m.Paginator.Page)
			}
		})
	}
}

func TestIsSearching(t *testing.T) {
	tests := []struct {
		name     string
		state    ViewState
		expected bool
	}{
		{
			name:     "searching groups",
			state:    SearchingGroups,
			expected: true,
		},
		{
			name:     "searching users",
			state:    SearchingUsers,
			expected: true,
		},
		{
			name:     "viewing groups",
			state:    ViewingGroups,
			expected: false,
		},
		{
			name:     "viewing group members",
			state:    ViewingGroupMembers,
			expected: false,
		},
		{
			name:     "viewing users",
			state:    ViewingUsers,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Model{TUIState: tt.state}
			result := m.IsSearching()
			if result != tt.expected {
				t.Errorf("IsSearching() = %v, want %v", result, tt.expected)
			}
		})
	}
}
