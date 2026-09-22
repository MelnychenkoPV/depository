package cmd

import (
	"encoding/json"
	"os"
)

type (
	Config struct {
		Server Server `json:"server"`
		DB     DB     `json:"db"`
	}
	Server struct {
		Addr string `json:"addr"`
	}
	DB struct {
		DSN string `json:"dsn"`
	}
)

func ParseConfig() (Config, error) {
	f, err := os.Open("config.json")
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var cfg Config
	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
