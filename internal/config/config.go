package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server     Server     `yaml:"server"`
	Database   Database   `yaml:"database"`
	Repository Repository `yaml:"repository"`
	Gateway    Gateway    `yaml:"gateway"`
}

type Server struct {
	Port    string `yaml:"port:"`
	Timeout int64  `yaml:"timeout"`
}

type Database struct {
	HostLocal  string `yaml:"host-local"`
	HostDocker string `yaml:"host-local"`
	Port       string `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	DBName     string `yaml:"dbname"`
}

type Repository struct {
	Timeout int64 `yaml:"timeout"`
}

type Gateway struct {
	Timeout int64 `yaml:"timeout"`
}

func ReadConfigYml(filePath string) (Config, error) {
	f, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	cfg := Config{}
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
