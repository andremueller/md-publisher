package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config is the configuration of the md-publisher tool.
type Config struct {
	NoImages          bool   // do not upload images
	MediumAccessToken string // access token for medium
}

// GetDefaultConfigFile returns the path to the default md-publisher.conf file
func GetDefaultConfigFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Cannot find user home directory", err)
	}
	return filepath.Join(home, ".config/md-publisher/md-publisher.conf")
}

// ReadConfig reads the configuration from a TOML (a ini like file) file.
func ReadConfig(configFile string, config *Config) error {
	_, err := os.Stat(configFile)
	if err != nil {
		return fmt.Errorf("Config file is missing: %s", configFile)
	}

	_, err = toml.DecodeFile(configFile, config)

	return err
}
