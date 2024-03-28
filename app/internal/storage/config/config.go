package configdb

import (
	"os"
)

type ConfigDb struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func (c *ConfigDb) Read() {
	c.Host = os.Getenv("POSTGRES_HOST")
	c.Port = os.Getenv("POSTGRES_PORT")
	c.Database = os.Getenv("POSTGRES_DB")
	c.Username = os.Getenv("POSTGRES_USER")
	c.Password = os.Getenv("POSTGRES_PASSWORD")
}
