package client

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/go-ldap/ldap/v3"
	"github.com/timskovjacobsen/ldapget/config"
)

func groupsRequest(conn *ldap.Conn, baseDN string) (*ldap.SearchResult, error) {
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
func Groups(cfg *config.Config) ([]GroupInfo, error) {
	baseDN := cfg.Client.Search.BaseDN
	conn, err := BindToLdapServer(*cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to bind to ldap server: %w", err)
	}
	defer conn.Close()

	result, err := groupsRequest(conn, baseDN)
	if err != nil {
		return nil, fmt.Errorf("failed to search LDAP server for groups: %w", err)
	}
	if len(result.Entries) == 0 {
		return nil, fmt.Errorf("no groups found")
	}
	var groups []GroupInfo
	for _, entry := range result.Entries {
		groupType, _ := strconv.ParseInt(entry.GetAttributeValue("groupType"), 10, 64)
		scope, kind, isSystem := groupTypeInfo(groupType)

		group := GroupInfo{
			Name:          entry.GetAttributeValue("cn"),
			DN:            entry.GetAttributeValue("distinguishedName"),
			Scope:         scope,
			Type:          kind,
			SystemCreated: isSystem,
			Description:   entry.GetAttributeValue("description"),
			MemberCount:   len(entry.GetAttributeValues("member")),
		}
		groups = append(groups, group)
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})
	return groups, nil
}
