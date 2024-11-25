package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

type LdapServer struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type Search struct {
	RootDN string `toml:"root_dn"`
}

type Client struct {
	Search Search `toml:"search"`
}

type Config struct {
	LdapServer LdapServer `toml:"ldap-server"`
	Client     Client     `toml:"client"`
}

func defaultConfigFile() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("no user config directory found: %v", err)
	}
	return filepath.Join(configDir, "ldapget", "config.toml"), nil
}

func ReadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		var err error
		configPath, err = defaultConfigFile()
		fmt.Println(configPath)
		if err != nil {
			return &Config{}, fmt.Errorf("no config file found in %s", configPath)
		}
	}

	fmt.Println(configPath)
	f, err := os.Open(configPath)
	if err != nil {
		return &Config{}, fmt.Errorf("failed to read '%s': %v", configPath, err)
	}
	defer f.Close()

	var cfg Config
	decoder := toml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return &Config{}, fmt.Errorf("failed to decode config file: %v", err)
	}
	return &cfg, nil
}
