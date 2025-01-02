package client

import (
	"github.com/go-ldap/ldap/v3"
	"github.com/timskovjacobsen/ldapget/config"
)

func Users(cfg *config.Config) ([]UserInfo, error) {
	baseDN := cfg.Client.Search.BaseDN
	conn, err := BindToLdapServer(*cfg)
	if err != nil {
		return nil, err
	}

	// Search filter for all user objects
	filter := "(&(objectClass=user)(objectCategory=person))"

	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"sAMAccountName", "distinguishedName", "mail"},
		nil,
	)

	result, err := conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	users := make([]UserInfo, 0, len(result.Entries))
	for _, entry := range result.Entries {
		user := UserInfo{
			Name:  entry.GetAttributeValue("sAMAccountName"),
			DN:    entry.GetAttributeValue("distinguishedName"),
			Email: entry.GetAttributeValue("mail"),
		}
		users = append(users, user)
	}

	return users, nil
}

// Usage example:
/*
func main() {
    // Assuming you have an established LDAP connection and baseDN
    users, err := GetAllUsers(conn, baseDN)
    if err != nil {
        log.Fatal(err)
    }

    for _, user := range users {
        fmt.Printf("Name: %s, Email: %s\n", user.Name, user.Email)
    }
}
*/
