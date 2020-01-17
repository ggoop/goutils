package utils

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/ggoop/goutils/glog"
)

type AppConfig struct {
	Token   string `mapstructure:"token" json:"token"`
	Code    string `mapstructure:"code" json:"code"`
	Active  string `mapstructure:"active" json:"active"`
	Name    string `mapstructure:"name" json:"name"`
	Port    string `mapstructure:"port" json:"port"`
	Locale  string `mapstructure:"locale" json:"locale"`
	Debug   bool   `mapstructure:"debug" json:"debug"`
	Storage string `mapstructure:"storage" json:"storage"`
	//注册中心
	Registry string `mapstructure:"registry" json:"registry"`
	//服务地址，带端口号
	Address string `mapstructure:"address" json:"address"`
}
type DbConfig struct {
	Driver    string `mapstructure:"driver" json:"driver"`
	Host      string `mapstructure:"host" json:"host"`
	Port      string `mapstructure:"port" json:"port"`
	Database  string `mapstructure:"database" json:"database"`
	Username  string `mapstructure:"username" json:"username"`
	Password  string `mapstructure:"password" json:"password"`
	Charset   string `mapstructure:"charset" json:"charset"`
	Collation string `mapstructure:"collation" json:"collation"`
}
type LogConfig struct {
	Level string `mapstructure:"level" json:"level"`
	Path  string `mapstructure:"path" json:"path"`
	Stack bool   `mapstructure:"stack" json:"stack"`
}
type AuthConfig struct {
	//权限中心地址
	Address string `mapstructure:"address" json:"address"`
	//权限中心编码
	Code string `mapstructure:"code" json:"code"`
}

const AppConfigName = "app"

type jsonConfig struct {
	App  AppConfig              `mapstructure:"app" json:"app"`
	Db   DbConfig               `mapstructure:"db" json:"db"`
	Log  LogConfig              `mapstructure:"log" json:"log"`
	Auth AuthConfig             `mapstructure:"auth" json:"auth"`
	Data map[string]interface{} `json:"data"`
}

type Config struct {
	App  AppConfig  `mapstructure:"app" json:"app"`
	Db   DbConfig   `mapstructure:"db" json:"db"`
	Log  LogConfig  `mapstructure:"log" json:"log"`
	Auth AuthConfig `mapstructure:"auth" json:"auth"`
	data map[string]interface{}
}

func (c Config) MarshalJSON() ([]byte, error) {
	jsonMap := jsonConfig{}
	jsonMap.App = c.App
	jsonMap.Db = c.Db
	jsonMap.Log = c.Log
	jsonMap.Auth = c.Auth
	jsonMap.Data = c.data
	return json.Marshal(jsonMap)
}

func (c *Config) UnmarshalJSON(b []byte) error {
	jsonMap := jsonConfig{}
	json.Unmarshal(b, &jsonMap)
	c.App = jsonMap.App
	c.Db = jsonMap.Db
	c.Log = jsonMap.Log
	c.Auth = jsonMap.Auth
	c.data = jsonMap.Data
	return nil
}
func (s *Config) GetValue(name string) string {
	v := s.GetObject(name)
	if v == nil {
		return ""
	}
	switch s := v.(type) {
	case string:
		return s
	case []byte:
		return string(s)
	}
	refV := reflectTarget(reflect.ValueOf(v))
	if !refV.IsValid() {
		return ""
	}
	switch refV.Kind() {
	case reflect.String:
		return refV.String()
	case reflect.Bool:
		if refV.Bool() {
			return "true"
		} else {
			return "false"
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", refV.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fmt.Sprintf("%d", refV.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", refV.Float())
	}
	return ""
}
func (s *Config) UnmarshalValue(name string, rawVal interface{}) error {
	return viper.UnmarshalKey(name, rawVal)
}

func (s *Config) GetBool(name string) bool {
	ov := s.GetObject(name)
	if ov == nil {
		return false
	}
	return ToBool(ov)
}
func (s *Config) GetObject(name string) interface{} {
	if s.data == nil {
		return nil
	}
	if v, ok := s.data[strings.ToLower(name)]; ok {
		return v
	}
	var d interface{}
	if err := s.UnmarshalValue(name, &d); err == nil {
		return d
	}
	return nil
}
func (s *Config) SetValue(name string, value interface{}) *Config {
	if s.data == nil {
		s.data = make(map[string]interface{})
	}
	s.data[strings.ToLower(name)] = value
	return s
}

var DefaultConfig *Config

func createEnvFile() {
	appFile := JoinCurrentPath("env/app.yaml")
	paths := []string{"env/app.prod.yaml", "env/app.dev.yaml"}
	for _, p := range paths {
		//不存在时，自动由dev创建
		if !PathExists(appFile) {
			devFile := JoinCurrentPath(p)
			if PathExists(devFile) {
				if s, err := os.Open(devFile); err == nil {
					defer s.Close()
					if newEnv, err := os.Create(appFile); err == nil {
						defer newEnv.Close()
						io.Copy(newEnv, s)
						break
					}
				}
			}
		}
	}
}
func NewInitConfig() {
	createEnvFile()

	DefaultConfig = &Config{}
	viper.SetConfigType("yaml")

	viper.SetConfigName(AppConfigName)
	viper.AddConfigPath(JoinCurrentPath("env"))
	if err := viper.ReadInConfig(); err != nil {
		glog.Errorf("Fatal error when reading %s config file:%s", AppConfigName, err)
	}
	if err := viper.Unmarshal(&DefaultConfig); err != nil {
		glog.Errorf("Fatal error when reading %s config file:%s", AppConfigName, err)
	}
	if DefaultConfig.App.Port == "" {
		DefaultConfig.App.Port = "8080"
	}
	if DefaultConfig.App.Locale == "" {
		DefaultConfig.App.Locale = "zh-cn"
	}
	if DefaultConfig.App.Token == "" {
		DefaultConfig.App.Token = "01e8f6a984101b20b24e4d172ec741be"
	}
	if DefaultConfig.App.Storage == "" {
		DefaultConfig.App.Storage = "./storage"
	}
	if DefaultConfig.Db.Driver == "" {
		DefaultConfig.Db.Driver = "mysql"
	}
	if DefaultConfig.Db.Host == "" {
		DefaultConfig.Db.Host = "localhost"
	}
	if DefaultConfig.Db.Port == "" {
		if DefaultConfig.Db.Driver == "mysql" {
			DefaultConfig.Db.Port = "3306"
		}
		if DefaultConfig.Db.Driver == "mssql" {
			DefaultConfig.Db.Port = "1433"
		}
	}
	if DefaultConfig.Db.Charset == "" {
		DefaultConfig.Db.Charset = "utf8mb4"
	}
	if DefaultConfig.Db.Collation == "" {
		DefaultConfig.Db.Collation = "utf8mb4_general_ci"
	}
	if DefaultConfig.Auth.Code == "" {
		DefaultConfig.Auth.Code = DefaultConfig.App.Code
	}
	kvs := make(map[string]interface{})
	if err := viper.Unmarshal(&kvs); err != nil {
		glog.Errorf("Fatal error when reading %s config file:%s", AppConfigName, err)
	}
	if len(kvs) > 0 {
		for k, v := range kvs {
			DefaultConfig.SetValue(k, v)
		}
	}
}
