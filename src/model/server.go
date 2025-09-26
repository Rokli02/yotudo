package model

import "fmt"

type ServerConfig struct {
	Port        int
	IsHosted    bool
	ListeningOn string
}

func (c *ServerConfig) Merge(other *ServerConfig) {
	c.Port = other.Port
	c.IsHosted = other.IsHosted
}

func (c *ServerConfig) String() string {
	return fmt.Sprintf(
		"ServerConfig(Port=%d, IsHosted=%t, ListeningOn=\"%s\")",
		c.Port, c.IsHosted, c.ListeningOn,
	)
}
