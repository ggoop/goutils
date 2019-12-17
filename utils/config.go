package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"

	"github.com/ggoop/goutils/glog"
)

type AppConfig struct {
	Token   string `mapstructure:"token"`
	Code    string `mapstructure:"code"`
	Active  string `mapstructure:"active"`
	Name    string `mapstructure:"name"`
	Port    string `mapstructure:"port"`
	Locale  string `mapstructure:"locale"`
	Debug   bool   `mapstructure:"debug"`
	Storage string `mapstructure:"storage"`
	//注册中心
	Registry string `mapstructure:"registry"`
	//服务地址，带端口号
	Address string `mapstructure:"address"`
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
	Path  string `mapstructure:"path"`
	Stack bool   `mapstructure:"stack"`
}
type AuthConfig struct {
	//权限中心地址
	Address string `mapstructure:"address"`
	//权限中心编码
	Code string `mapstructure:"code"`
}

const appConfigName = "app"

type Config struct {
	App  AppConfig
	Db   DbConfig
	Log  LogConfig
	Auth AuthConfig
	data map[string]interface{}
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
	ov := s.GetObject(name)
	if ov == nil {
		return nil
	}
	if b, err := json.Marshal(ov); err != nil {
		return err
	} else {
		if err := json.Unmarshal(b, rawVal); err != nil {
			return err
		}
	}
	return nil
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

func NewInitConfig() {
	appFile := JoinCurrentPath("env/app.yaml")
	//不存在时，自动由dev创建
	if !PathExists(appFile) {
		devFile := JoinCurrentPath("env/app.yaml.dev")
		if PathExists(devFile) {
			if s, err := os.Open(devFile); err == nil {
				defer s.Close()
				if newEnv, err := os.Create(appFile); err == nil {
					defer newEnv.Close()
					io.Copy(newEnv, s)
				}
			}
		}
	}
	DefaultConfig = &Config{}
	viper.SetConfigType("yaml")

	viper.SetConfigName(appConfigName)
	viper.AddConfigPath(JoinCurrentPath("env"))
	if err := viper.ReadInConfig(); err != nil {
		glog.Errorf("Fatal error when reading %s config file:%s", appConfigName, err)
	}
	if err := viper.Unmarshal(&DefaultConfig); err != nil {
		glog.Errorf("Fatal error when reading %s config file:%s", appConfigName, err)
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
		DefaultConfig.Db.Port = "3306"
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
		glog.Errorf("Fatal error when reading %s config file:%s", appConfigName, err)
	}
	if len(kvs) > 0 {
		for k, v := range kvs {
			if k == "app" || k == "db" || k == "log" || k == "auth" {
				continue
			}
			DefaultConfig.SetValue(k, v)
		}
	}
}
