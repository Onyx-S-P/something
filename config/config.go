package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Url      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func ReadConfig(path string) (Config, error) {

	c := Config{}
	data, err := os.ReadFile(path)
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(data, &c)
	return c, err
}
