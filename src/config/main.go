package config

import (
	"os"
)

type InstagramConfig struct {
	Username string
	Password string
}

type DropboxConfig struct {
	Token string
}

type FilesystemConfig struct {
	ImageDirectory string
}

type Config struct {
	Instagram InstagramConfig
	Filesystem FilesystemConfig
	Dropbox DropboxConfig
}

// New returns a new Config struct
func New() *Config {
    return &Config{
    Instagram: InstagramConfig{
	    Username: getEnv("INSTAGRAM_USERNAME", ""),
	    Password: getEnv("INSTAGRAM_PASSWORD", ""),
	},
	Dropbox: DropboxConfig{
		Token: getEnv("DROPBOX_ACCESS_TOKEN", ""),
	},
	Filesystem: FilesystemConfig{
		ImageDirectory: "./images/" + getEnv("INSTAGRAM_USERNAME", ""),
	},
    }
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
    if value, exists := os.LookupEnv(key); exists {
	return value
    }

    return defaultVal
}