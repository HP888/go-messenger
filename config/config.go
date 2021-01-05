package config

import "github.com/pelletier/go-toml"

type ServerConfig struct {
	Host string
	Port int
}

type Config struct {
	Server ServerConfig
}

func Load() Config {
	config := Config{}
	tomlConfig, _ := toml.LoadFile("config.toml")
	_ = toml.Unmarshal([]byte(tomlConfig.String()), &config)
	return config
}
