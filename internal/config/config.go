package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const configFile = ".gatorconfig.json"

type Config struct {
	DBurl    string `json:"db_url"`
	Username string `json:"username"`
}

func Read() *Config {
	data, err := os.ReadFile(configLocation())
	if err != nil {
		log.Fatal(err)
	}
	var config *Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func (c *Config) SetUser(user string) {
	c.Username = user
	data, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(configLocation(), data, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User has been set to: %s\n", user)
}

func configLocation() string {
	home, _ := os.UserHomeDir()
	return fmt.Sprintf("%s/%s", home, configFile)
}
