package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func (c *Config) Read() {
	filename, _ := filepath.Abs("../internal/server/config/config.yml")
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("can't read server config. Error: %e", err)
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		fmt.Printf("can't unmarshal server config. Error: %e", err)
	}
}
