package ebs

import (
	"fmt"
	"os"
)

// configuration
// to make environment variables available to testing
const EXTENSION_ID = "EXTENSION_ID"

type Config struct {
	extensionId string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) ExtensionId(id string) {
	c.extensionId = id
}

func (c *Config) Hostname() string {
	// format the hostname for the CORS allow-origin
	// 1. For Netlify, EXTENSION_ID environment variable should be defined
	// 2. Locally for testing, rely on configuration field
	if c.extensionId != "" {
		return fmt.Sprintf("https://%s.ext-twitch.tv", c.extensionId)
	}

	cid := os.Getenv(EXTENSION_ID)
	if cid != "" {
		c.extensionId = cid
		return fmt.Sprintf("https://%s.ext-twitch.tv", c.extensionId)
	}

	return ""
}
