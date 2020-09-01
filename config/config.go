package config

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"sync"
)

const DefaultConfigName = "v-blog"
const DefaultConfigType = "json"
const DefaultConfigPath = "/etc"

type Param struct {
	Name string
	Type string
	Path []string
}

type Config struct {
	mu      sync.Mutex
	loaded  bool
	Version float64
	Name    string
	Db      struct {
		Host     string
		Port     string
		User     string
		Password string
		Database string
		Charset  string
	}
}

var c *Config = &Config{}

func load(param *Param) {
	if c.loaded {
		panic("config has already loaded")
	}
	viper.SetConfigName(param.Name)
	viper.SetConfigType(param.Type)
	for _, path := range param.Path {
		viper.AddConfigPath(path)
	}
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	c.loaded = true

	read()
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		read()
	})
}

func read() {
	defer c.mu.Unlock()
	c.mu.Lock()
	c.Version = viper.GetFloat64("version")
	c.Name = viper.GetString("name")
	c.Db.Host = viper.GetString("db.host")
	c.Db.Port = viper.GetString("db.port")
	c.Db.User = viper.GetString("db.user")
	c.Db.Password = viper.GetString("db.password")
	c.Db.Database = viper.GetString("db.database")
	c.Db.Charset = viper.GetString("db.charset")
}

func InitConfig(param *Param) {
	if param.Name == "" {
		param.Name = DefaultConfigName
	}
	if param.Type == "" {
		param.Type = DefaultConfigType
	}
	if param.Path == nil {
		param.Path = []string{DefaultConfigPath, "."}
	}

	load(param)
}

func Get() (*Config, error) {
	if c.loaded == false {
		return c, errors.New("config has not loaded")
	}

	return c, nil
}
