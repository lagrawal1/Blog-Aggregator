package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() Config {
	var conf Config
	hom_dir, err := os.UserHomeDir()

	if err != nil {
		fmt.Println(err)
		return Config{}
	}
	data, err := os.ReadFile(hom_dir + "/.gatorconfig.json")

	if err != nil {
		fmt.Println(err)
		return Config{}
	}

	err = json.Unmarshal(data, &conf)

	if err != nil {
		fmt.Println(err)
		return Config{}
	}

	return conf
}

func (conf *Config) SetUser(curr_user string) error {
	conf.CurrentUserName = curr_user
	hom_dir, err := os.UserHomeDir()

	if err != nil {
		return err
	}
	config_file := hom_dir + "/.gatorconfig.json"

	data, err := json.Marshal(conf)

	if err != nil {
		return err
	}

	permissions := 0644

	err = os.WriteFile(config_file, data, os.FileMode(permissions))

	if err != nil {
		return err
	}

	return nil

}
