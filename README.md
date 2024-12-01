# ldapget

Example config file in `~/.config/ldapget/config.toml`:

```shell
[ldap-server]
host = "ldaps://host.domain.com"
port = 636
username = "some-user"
password = "..."

[client.search]
base_dn = "DC=company,DC=com"
```
