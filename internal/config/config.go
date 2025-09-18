package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host struct {
		Addr string `yaml:"addr"`
		Port string `yaml:"port"`
	} `yaml:"host"`
	Name string   `yaml:"name"`
	Peer []string `yaml:"peers"`
	Dir  string   `yaml:"dir"`
}

func ReadFrom(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}

	c := Config{}
	if err := yaml.NewDecoder(f).Decode(&c); err != nil {
		return Config{}, err
	}

	return c, nil
}
