package configs

import (
	"os"
	"strings"

	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/utils"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Token  string `mapstructure:"token"`
	Code   string `mapstructure:"code"`
	Name   string `mapstructure:"name"`
	Port   string `mapstructure:"port"`
	Locale string `mapstructure:"locale"`
	Debug  bool   `mapstructure:"debug"`
}
type DbConfig struct {
	Driver    string `mapstructure:"driver"`
	Host      string `mapstructure:"host"`
	Port      string `mapstructure:"port"`
	Database  string `mapstructure:"database"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	Charset   string `mapstructure:"charset"`
	Collation string `mapstructure:"collation"`
}
type LogConfig struct {
	Level string `mapstructure:"level"`
}

const appConfigName = "app"

type Config struct {
	App  AppConfig
	Db   DbConfig
	Log  LogConfig
	data map[string]string
}

func (s *Config) GetValue(name string) string {
	if s.data == nil {
		return ""
	}
	return s.data[strings.ToLower(name)]
}
func (s *Config) SetValue(name, value string) *Config {
	if s.data == nil {
		s.data = make(map[string]string)
	}
	s.data[strings.ToLower(name)] = value
	return s
}

var Default *Config

func init() {
	New()
}
func New() {
	Default = &Config{}
	viper.SetConfigType("yaml")

	viper.SetConfigName(appConfigName)
	viper.AddConfigPath(utils.JoinCurrentPath("env"))
	err := viper.ReadInConfig()
	if err != nil {
		glog.Errorf("Fatal error when reading %s config file:%s", appConfigName, err)
		os.Exit(1)
	}
	viper.Unmarshal(&Default)
	if Default.App.Port == "" {
		Default.App.Port = "8080"
	}
	if Default.App.Locale == "" {
		Default.App.Locale = "zh-cn"
	}
	if Default.App.Token == "" {
		Default.App.Token = "01e8f6a984101b20b24e4d172ec741be"
	}
	if Default.Db.Driver == "" {
		Default.Db.Driver = "mysql"
	}
	if Default.Db.Host == "" {
		Default.Db.Host = "localhost"
	}
	if Default.Db.Port == "" {
		Default.Db.Port = "3306"
	}
	if Default.Db.Charset == "" {
		Default.Db.Charset = "utf8mb4"
	}
	if Default.Db.Collation == "" {
		Default.Db.Collation = "utf8mb4_general_ci"
	}
	kvs := make(map[string]string)
	viper.Unmarshal(&kvs)
	if len(kvs) > 0 {
		for k, v := range kvs {
			Default.SetValue(k, v)
		}
	}
}
