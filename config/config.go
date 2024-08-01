package config

import (
	"logger"
	"runtime"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Name string
}

func (c *Config) initConfig() error {
	if c.Name != "" {
		// 若指定了配置文件，则解析指定配置文件
		viper.SetConfigFile(c.Name)
	} else {
		// 否则解析默认配置文件
		viper.AddConfigPath("conf")
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
func (c *Config) listenConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logger.Logger.Printf("Config file changed: %s\n", e.Name)
	})
}

func ConfigInit(name ...string) {
	c := &Config{}
	if len(name) != 0 {
		c.Name = name[0]
	} else {
		c.Name = ""
	}

	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		logger.Logger.Println(err)
	}

	if viper.GetInt("maxThreadNum") == 0 {
		viper.Set("maxThreadNum", runtime.GOMAXPROCS(0))
	}

	// 监控配置文件变化并热加载程序
	c.listenConfig()
}
