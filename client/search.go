package client

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
)

func Groups(conn *ldap.Conn, baseDN string) (*ldap.SearchResult, error) {
	filter := fmt.Sprintf("(&(objectClass=group))")
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		// []string{"mail", "memberOf"},
		// []string{"cn"},
		nil,
		nil,
	)
	return conn.Search(searchRequest)
}

func Group(conn *ldap.Conn, baseDN string, groupName string) (*ldap.SearchResult, error) {
	//
	filter := fmt.Sprintf("(&(objectClass=group)(cn=%s))", ldap.EscapeFilter(groupName))
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"member"},
		nil,
	)
	result, err := conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	if len(result.Entries) == 0 {
		return nil, fmt.Errorf("group not found: %s", groupName)
	}
	members := result.Entries[0].GetAttributeValues("member")
	if len(members) == 0 {
		return result, nil // group is empty
	}
	memberFilter := "(&(objectClass=user)(|"
	for _, member := range members {
		memberFilter += fmt.Sprintf("(distinguishedName=%s)", ldap.EscapeFilter(member))
	}
	memberFilter += "))"

	userRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		memberFilter,
		[]string{"cn", "sAMAccountName", "mail"},
		nil,
	)
	return conn.Search(userRequest)
}

func User(conn *ldap.Conn, baseDN string, user string) (*ldap.SearchResult, error) {
	// Search for the user and retrieve their groups
	filter := fmt.Sprintf("(&(objectClass=user)(sAMAccountName=%s))", ldap.EscapeFilter(user))
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"mail", "memberOf"},
		nil,
	)
	return conn.Search(searchRequest)
}
