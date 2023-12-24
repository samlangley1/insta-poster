package config

import (
	"os"
)

type InstagramConfig struct {
	Username string
	Password string
}

type FilesystemConfig struct {
	ImageDirectory string
}

type NetworkConfig struct {
	ProxyAddress string
}

type Config struct {
	Instagram  InstagramConfig
	Filesystem FilesystemConfig
	Network    NetworkConfig
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		Instagram: InstagramConfig{
			Username: getEnv("INSTAGRAM_USERNAME", ""),
			Password: getEnv("INSTAGRAM_PASSWORD", ""),
		},
		Filesystem: FilesystemConfig{
			ImageDirectory: "./images/" + getEnv("INSTAGRAM_USERNAME", ""),
		},
		Network: NetworkConfig{
			ProxyAddress: getEnv("HTTP_PROXY", ""),
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
