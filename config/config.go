package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type config struct {
	Version string
	Name string
}

var Config config

func init()  {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("read config fail [%s]", err))
	}

	writeConfig()
}

func writeConfig() {
	Config = config{
		Name: viper.GetString("name"),
		Version: viper.GetString("version"),
	}
}

