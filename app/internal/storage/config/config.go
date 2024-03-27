package configdb

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ConfigDb struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}

func (c *ConfigDb) Read() {
	filename, _ := filepath.Abs("../internal/storage/config/config.yml")
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("can't read storage config. Error: %e", err)
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		fmt.Printf("can't unmarshal storage config. Error: %e", err)
	}
}
