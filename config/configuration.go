package config

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	DBDriver  string `json:"db_driver"`
	DBName    string `json:"db_name"`
	JWTSecret string `json:"jwt_secret"`
}

const configFilePath = "config/"

var validEnvs = []string{
	"dev",
	"test",
}

func isValidEnv(env string) bool {
	for _, validEnv := range validEnvs {
		if env == validEnv {
			return true
		}
	}
	return false
}
func GetConfig(env string) (config Configuration, err error) {
	if !isValidEnv(env) {
		// default env to development
		env = "dev"
	}
	//var config Configuration
	filename := configFilePath + env + ".json"

	configFile, err := os.Open(filename)
	if err != nil {
		return
	}

	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return
	}
	return
}
