package client

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"

	"github.com/timskovjacobsen/ldapget/config"
)

func BindToLdapServer(cfg config.Config) (*ldap.Conn, error) {
	// Connect to the LDAP server
	conn, err := ldap.DialURL(fmt.Sprintf("%s:%d", cfg.LdapServer.Host, cfg.LdapServer.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to LDAP server: %v", err)
	}

	// Bind with a read-only user
	err = conn.Bind(cfg.LdapServer.Username, cfg.LdapServer.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to bind to LDAP server: %v", err)
	}
	return conn, nil
}
