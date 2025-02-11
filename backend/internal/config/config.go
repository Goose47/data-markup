// Package config describes required configuration and provides methods for reading configuration from files.
package config

import (
	"errors"
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

var errFileNotExists = errors.New("config file does not exist")
var errBadConfigFile = errors.New("unable to read config file")

// Config represents main app configuration.
type Config struct {
	Env  string `yaml:"env"`
	Port int    `yaml:"port"`
	DB   DB     `yaml:"db"`
}

// Config represents database configuration.
type DB struct {
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
	DBName string `yaml:"dbname"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
}

// LoadPath loads configuration from specified path and returns config instance and error.
func LoadPath(configPath string) (*Config, error) {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: %s", errFileNotExists, configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("%w: %s", errBadConfigFile, err.Error())
	}

	return &cfg, nil
}

// MustLoad fetches path, loads configuration and panics on any error.
func MustLoad() *Config {
	configPath := fetchConfigPath()

	cfg, err := LoadPath(configPath)
	if err != nil {
		panic(err)
	}

	return cfg
}

// MustLoadPath loads configuration form configPath and panics on any error.
func MustLoadPath(configPath string) *Config {
	cfg, err := LoadPath(configPath)
	if err != nil {
		panic(err)
	}

	return cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	if flag.Lookup("config") == nil {
		flag.StringVar(&res, "config", "", "path to config file")
		flag.Parse()
	} else {
		res = flag.Lookup("config").Value.String()
	}

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
