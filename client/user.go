package client

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
)

type UserInfo struct {
	Name  string
	DN    string
	Email string
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
