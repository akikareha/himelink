package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Route struct {
	Path     string `yaml:"path"`
	Protocol string `yaml:"protocol"`
	API      string `yaml:"api"`
	Raw      string `yaml:"raw,omitempty"`
}

type Config struct {
	App struct {
		Mode string `yaml:"mode"`
		Addr string `yaml:"addr"`
	} `yaml:"app"`

	Site struct {
		Base   string `yaml:"base"`
		Name   string `yaml:"name"`
		Static string `yaml:"static"`
		Card   string `yaml:"card"`
	} `yaml:"site"`

	Routes []Route `yaml:"routes"`
}

func Load(path string) *Config {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	return &cfg
}
