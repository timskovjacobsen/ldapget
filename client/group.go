package client

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
	"github.com/timskovjacobsen/ldapget/config"
)

type GroupInfo struct {
	// Note: descriptions below are based an Active Directory
	Name          string
	DN            string // Distinguished Name, i.e. the locaiton of the group
	SystemCreated bool   // whether group is created by the system or not
	Scope         string // "DomainLocal", "Global", or "Universal"
	Type          string // "Security" or "Distribution"
	Description   string
	MemberCount   int
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

func groupRequest(conn *ldap.Conn, baseDN string, groupName string) (*ldap.SearchResult, error) {
	// First get the group's DN and members
	filter := fmt.Sprintf("(&(objectClass=group)(cn=%s))", ldap.EscapeFilter(groupName))
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"distinguishedName", "member"},
		nil,
	)

	result, err := conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	if len(result.Entries) == 0 {
		return nil, fmt.Errorf("group not found: %s", groupName)
	}

	// Use the found group DN in a memberOf filter
	groupDN := result.Entries[0].DN
	memberFilter := fmt.Sprintf("(&(objectClass=user)(memberOf=%s))", ldap.EscapeFilter(groupDN))

	userRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		memberFilter,
		[]string{"cn", "sAMAccountName", "mail"},
		nil,
	)
	return conn.Search(userRequest)
}

func GroupMembers(groupName string, cfg *config.Config) ([]UserInfo, error) {
	baseDN := cfg.Client.Search.BaseDN
	conn, err := BindToLdapServer(*cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to bind to ldap server")
	}
	membersResult, err := groupRequest(conn, baseDN, groupName)
	if err != nil {
		return nil, fmt.Errorf("group search failed")
	}
	if len(membersResult.Entries) == 0 {
		return nil, fmt.Errorf("group not found: %s", groupName)
	}

	var Users []UserInfo
	for _, entry := range membersResult.Entries {
		user := UserInfo{
			Name:  entry.GetAttributeValue("cn"),
			Email: entry.GetAttributeValue("mail"),
			DN:    entry.DN,
		}
		Users = append(Users, user)
	}
	return Users, nil
}
