package config

import (
    "log"
    "os"

    "gopkg.in/yaml.v3"
)

type Config struct {
    Host string `yaml:"host"`
    Port int    `yaml:"port"`
    Env  string `yaml:"env"`
}

func DefaultConfig() *Config {
    return &Config{
        Host: "0.0.0.0",
        Port: 8080,
        Env:  "development",
    }
}

func New(path string) *Config {
    cfg := DefaultConfig()

    file, err := os.ReadFile(path)
    if err != nil {
        log.Printf("YAML config not found, using defaults: %v", err)
        return cfg
    }

    err = yaml.Unmarshal(file, cfg)
    if err != nil {
        log.Fatalf("Failed to parse YAML: %v", err)
    }

    return cfg
}