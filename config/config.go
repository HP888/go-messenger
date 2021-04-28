package config

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"os"
	"path/filepath"
)

type ServerConfig struct {
	Host string
	Port int
}

type Config struct {
	Server ServerConfig
}

func Load() Config {
	configPath, err := filepath.Abs("config.toml")
	if err != nil {
		fmt.Println("Cannot resolve absolute path to configuration file!")
		panic(err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		file, err := os.Create(configPath)
		if err != nil {
			fmt.Println("Cannot create configuration file!")
			panic(err)
		}

		defer file.Close()

		_, err = fmt.Fprintf(file,
			"#################################\n"+
				"#          GoMessenger\n"+
				"#          Version 0.1\n"+
				"#################################\n"+
				"\n"+
				"[server]\n"+
				"host = '0.0.0.0'\n"+
				"port = 9999",
		)

		if err != nil {
			fmt.Println("Cannot save configuration file!")
			panic(err)
		}
	}

	config := Config{}
	tomlConfig, _ := toml.LoadFile(configPath)
	_ = toml.Unmarshal([]byte(tomlConfig.String()), &config)
	return config
}
