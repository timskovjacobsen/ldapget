package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadConfig(t *testing.T) {
	tests := []struct {
		name        string
		configData  string
		wantErr     bool
		expectedCfg *Config
	}{
		{
			name: "valid config",
			configData: `
[ldap-server]
host = "ldaps://test.example.com"
port = 636
username = "testuser"
password = "testpass"

[client.search]
base_dn = "DC=test,DC=com"
`,
			wantErr: false,
			expectedCfg: &Config{
				LdapServer: LdapServer{
					Host:     "ldaps://test.example.com",
					Port:     636,
					Username: "testuser",
					Password: "testpass",
				},
				Client: Client{
					Search: Search{
						BaseDN: "DC=test,DC=com",
					},
				},
			},
		},
		{
			name:       "invalid toml",
			configData: `this is not valid toml [[[`,
			wantErr:    true,
		},
		{
			name: "partial config",
			configData: `
[ldap-server]
host = "ldaps://test.example.com"
port = 636
`,
			wantErr: false,
			expectedCfg: &Config{
				LdapServer: LdapServer{
					Host: "ldaps://test.example.com",
					Port: 636,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.toml")

			if err := os.WriteFile(configPath, []byte(tt.configData), 0644); err != nil {
				t.Fatalf("failed to write test config: %v", err)
			}

			cfg, err := ReadConfig(configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.expectedCfg != nil {
				if cfg.LdapServer.Host != tt.expectedCfg.LdapServer.Host {
					t.Errorf("Host = %v, want %v", cfg.LdapServer.Host, tt.expectedCfg.LdapServer.Host)
				}
				if cfg.LdapServer.Port != tt.expectedCfg.LdapServer.Port {
					t.Errorf("Port = %v, want %v", cfg.LdapServer.Port, tt.expectedCfg.LdapServer.Port)
				}
				if cfg.LdapServer.Username != tt.expectedCfg.LdapServer.Username {
					t.Errorf("Username = %v, want %v", cfg.LdapServer.Username, tt.expectedCfg.LdapServer.Username)
				}
				if cfg.Client.Search.BaseDN != tt.expectedCfg.Client.Search.BaseDN {
					t.Errorf("BaseDN = %v, want %v", cfg.Client.Search.BaseDN, tt.expectedCfg.Client.Search.BaseDN)
				}
			}
		})
	}
}

func TestReadConfig_NonExistentFile(t *testing.T) {
	_, err := ReadConfig("/this/path/does/not/exist/config.toml")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}
