package client

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/go-ldap/ldap/v3"
	"github.com/timskovjacobsen/ldapget/config"
)

type GroupInfo struct {
	// Note: descriptions below are based an Active Directory
	Name          string
	DN            string // Distinguished Name, i.e. the locaiton of the group
	SystemCreated bool   // whether group is created by the system or not
	Scope         string // "DomainLocal", "Global", or "Universal"
	Kind          string // "Security" or "Distribution"
	Description   string
	Members       int // members count
}

// Based on the Microsoft documentation from:
// https://learn.microsoft.com/en-us/windows/win32/adschema/a-grouptype?redirectedfrom=MSDN
// Value                   Description
// -----                   -----------
// 1 (0x00000001)          group that is created by the system
// 2 (0x00000002)          group with global scope
// 4 (0x00000004)          group with domain local scope
// 8 (0x00000008)          group with universal scope
// 16 (0x00000010)         APP_BASIC group for Windows Server Authorization Manager
// 32 (0x00000020)         APP_QUERY group for Windows Server Authorization Manager
// 2147483648 (0x80000000) if set, the group is a security group, else a distr. group
const (
	GROUP_SYSTEM       = 0x00000001
	GROUP_GLOBAL       = 0x00000002
	GROUP_DOMAIN_LOCAL = 0x00000004
	GROUP_UNIVERSAL    = 0x00000008
	GROUP_APP_BASIC    = 0x00000010
	GROUP_APP_QUERY    = 0x00000020
	GROUP_SECURITY     = 0x80000000
)

func groupTypeInfo(groupType int64) (scope, kind string, isSystem bool) {
	switch {
	case groupType&GROUP_GLOBAL != 0:
		scope = "Global"
	case groupType&GROUP_DOMAIN_LOCAL != 0:
		scope = "Domain Local"
	case groupType&GROUP_UNIVERSAL != 0:
		scope = "Universal"
	default:
		scope = "Unknown"
	}
	if groupType&GROUP_SECURITY != 0 {
		kind = "Security"
	} else {
		kind = "Distribution"
	}
	isSystem = groupType&GROUP_SYSTEM != 0
	return scope, kind, isSystem
}

func GroupsRequest(conn *ldap.Conn, baseDN string) (*ldap.SearchResult, error) {
	filter := fmt.Sprintf("(&(objectClass=group))")
	attributes := []string{
		"cn",
		"distinguishedName",
		"groupType",
		"description",
		"member",
		"sAMAccountName",
	}

	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		attributes,
		nil,
	)
	return conn.Search(searchRequest)
}

// Return a list of groups with all relevant info attached
func Groups(cfg *config.Config) []GroupInfo {
	baseDN := cfg.Client.Search.RootDN
	conn, err := BindToLdapServer(*cfg)
	if err != nil {
		log.Fatalf("failed to bind to ldap server: %v", err)
	}
	result, err := GroupsRequest(conn, baseDN)
	if err != nil {
		log.Fatalf("failed to search LDAP server for groups: %v", err)
	}
	if len(result.Entries) == 0 {
		log.Fatalf("no groups found")
	}
	var groups []GroupInfo
	for _, entry := range result.Entries {
		// Decode the insane group type that is returned
		groupType, _ := strconv.ParseInt(entry.GetAttributeValue("groupType"), 10, 64)
		scope, kind, isSystem := groupTypeInfo(groupType)

		group := GroupInfo{
			Name:          entry.GetAttributeValue("cn"),
			DN:            entry.GetAttributeValue("distinguishedName"),
			Scope:         scope,
			Kind:          kind,
			SystemCreated: isSystem,
			Description:   entry.GetAttributeValue("description"),
			Members:       len(entry.GetAttributeValues("member")),
		}
		groups = append(groups, group)
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})
	return groups
}
