package config

import (
	"os"
)

type Config struct {
	Host string
	Port string
}

func (c *Config) Read() {
	c.Host = os.Getenv("SERVER_HOST")
	c.Port = os.Getenv("SERVER_PORT")
}
