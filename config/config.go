package config

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"sync"
)

type Config struct {
	mu           sync.Mutex
	loaded       bool
	App struct{
		Env string
		Host string
		Name string
		PidFile string
		LogFile string
		DebugLogFile string
	}
	Dir struct{
		UploadImage string
	}
	Mysql struct {
		Host     string
		Port     string
		User     string
		Password string
		Database string
		Charset  string
	}
}

var c = &Config{}

func load() {
	if c.loaded {
		panic("config has already loaded")
	}
	viper.SetConfigFile("v-blog.ini")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath(".")

	// 默认设置
	viper.SetDefault("app.pidFile", "/var/tmp/v-blog.pid")
	viper.SetDefault("app.logFile", "/var/logs/v-blog.log")
	viper.SetDefault("app.debugLogFile", "/var/logs/v-blog-debug.log")
	viper.SetDefault("dir.uploadImage", "./upload/images")

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
	err := viper.Unmarshal(c)
	if err != nil {
		panic(fmt.Sprintf("config unmarshal err: %s", err))
	}
}

func InitConfig() {
	load()
}

func Get() (*Config, error) {
	if c.loaded == false {
		return c, errors.New("config has not loaded")
	}

	return c, nil
}
