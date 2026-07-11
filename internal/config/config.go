// Package config resolves garminctl's tiny configuration: which profiles (Garmin accounts)
// exist and which is the default. Credentials never live here — they are in the OS keyring,
// keyed by profile. Precedence for the active profile is flag > env > config > "default".
package config

import (
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v4"
)

// Config is the on-disk state: the known profiles and the default one.
type Config struct {
	DefaultProfile string   `yaml:"default_profile"`
	Profiles       []string `yaml:"profiles"`
}

// Dir returns the configuration directory: $XDG_CONFIG_HOME/garminctl if set, else
// ~/.garminctl-cli.
func Dir() (string, error) {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "garminctl"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".garminctl-cli"), nil
}

// Path returns the config file path.
func Path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

// Load reads the config, returning an empty Config (not an error) when the file is absent.
func Load() (*Config, error) {
	p, err := Path()
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(p) // #nosec G304 -- path derived from HOME/XDG, not user input
	if os.IsNotExist(err) {
		return &Config{}, nil
	}
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// Save writes the config, creating the directory if needed.
func Save(c *Config) error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	p, err := Path()
	if err != nil {
		return err
	}
	return os.WriteFile(p, b, 0o600)
}

// AddProfile records a profile name and, if it's the first, makes it the default.
func (c *Config) AddProfile(name string) {
	for _, p := range c.Profiles {
		if p == name {
			return
		}
	}
	c.Profiles = append(c.Profiles, name)
	if c.DefaultProfile == "" {
		c.DefaultProfile = name
	}
}

// Resolve returns the active profile from flag > env > config default > "default".
func Resolve(flag string) string {
	if flag != "" {
		return flag
	}
	if env := os.Getenv("GARMINCTL_PROFILE"); env != "" {
		return env
	}
	if c, err := Load(); err == nil && c.DefaultProfile != "" {
		return c.DefaultProfile
	}
	return "default"
}
