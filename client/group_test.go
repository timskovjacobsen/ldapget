package client

import (
	"testing"
)

func TestGroupTypeInfo(t *testing.T) {
	tests := []struct {
		name          string
		groupType     int64
		expectedScope string
		expectedKind  string
		expectedSys   bool
	}{
		{
			name:          "global security group",
			groupType:     GROUP_GLOBAL | GROUP_SECURITY,
			expectedScope: "Global",
			expectedKind:  "Security",
			expectedSys:   false,
		},
		{
			name:          "domain local distribution group",
			groupType:     GROUP_DOMAIN_LOCAL,
			expectedScope: "Domain Local",
			expectedKind:  "Distribution",
			expectedSys:   false,
		},
		{
			name:          "universal security group",
			groupType:     GROUP_UNIVERSAL | GROUP_SECURITY,
			expectedScope: "Universal",
			expectedKind:  "Security",
			expectedSys:   false,
		},
		{
			name:          "system created global security group",
			groupType:     GROUP_SYSTEM | GROUP_GLOBAL | GROUP_SECURITY,
			expectedScope: "Global",
			expectedKind:  "Security",
			expectedSys:   true,
		},
		{
			name:          "global distribution group",
			groupType:     GROUP_GLOBAL,
			expectedScope: "Global",
			expectedKind:  "Distribution",
			expectedSys:   false,
		},
		{
			name:          "unknown scope",
			groupType:     GROUP_SECURITY,
			expectedScope: "Unknown",
			expectedKind:  "Security",
			expectedSys:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope, kind, isSystem := groupTypeInfo(tt.groupType)

			if scope != tt.expectedScope {
				t.Errorf("scope = %v, want %v", scope, tt.expectedScope)
			}
			if kind != tt.expectedKind {
				t.Errorf("kind = %v, want %v", kind, tt.expectedKind)
			}
			if isSystem != tt.expectedSys {
				t.Errorf("isSystem = %v, want %v", isSystem, tt.expectedSys)
			}
		})
	}
}
