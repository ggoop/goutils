package repositories

import (
	"bytes"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/ggoop/goutils/configs"
	"github.com/ggoop/goutils/glog"

	"github.com/ggoop/goutils/gorm"
)

const (
	//mssql
	Driver_MSSQL = "mssql"
	//mysql
	Driver_MYSQL = "mysql"
	//oci8
	Driver_OCI8 = "oci8"
)

type MysqlRepo struct {
	*gorm.DB
}

func NewMysqlRepo() *MysqlRepo {
	// 生成数据库
	CreateDB(configs.Default.Db.Database)

	db, err := gorm.Open(configs.Default.Db.Driver, getDsnString(true))
	if err != nil {
		glog.Errorf("orm failed to initialized: %v", err)
	}

	db.LogMode(configs.Default.App.Debug)
	return createMysqlRepo(db)
}
func getDsnString(inDb bool) string {
	//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	//mssql:   =>  sqlserver://username:password@localhost:1433?database=dbname
	//mysql => user:password@(localhost)/dbname?charset=utf8&parseTime=True&loc=Local
	str := ""
	// 创建连接
	if configs.Default.Db.Driver == Driver_MSSQL {
		var buf bytes.Buffer
		buf.WriteString("sqlserver://")
		buf.WriteString(configs.Default.Db.Username)
		if configs.Default.Db.Password != "" {
			buf.WriteByte(':')
			buf.WriteString(configs.Default.Db.Password)
		}
		if configs.Default.Db.Host != "" {
			buf.WriteByte('@')
			buf.WriteString(configs.Default.Db.Host)
			if configs.Default.Db.Port != "" {
				buf.WriteByte(':')
				buf.WriteString(configs.Default.Db.Port)
			} else {
				buf.WriteString(":1433")
			}
		}
		if configs.Default.Db.Database != "" && inDb {
			buf.WriteString("?database=")
			buf.WriteString(configs.Default.Db.Database)
		} else {
			buf.WriteString("?database=master")
		}
		str = buf.String()
	} else {
		config := mysql.Config{
			User:   configs.Default.Db.Username,
			Passwd: configs.Default.Db.Password, Net: "tcp", Addr: configs.Default.Db.Host,
			AllowNativePasswords: true,
			ParseTime:            true,
			Loc:                  time.Local,
		}
		if inDb {
			config.DBName = configs.Default.Db.Database
		}
		if configs.Default.Db.Port != "" {
			config.Addr = fmt.Sprintf("%s:%s", configs.Default.Db.Host, configs.Default.Db.Port)
		}
		str = config.FormatDSN()
	}
	return str
}
func createMysqlRepo(db *gorm.DB) *MysqlRepo {
	repo := &MysqlRepo{db}
	return repo
}
func DestroyDB(name string) error {
	db, err := gorm.Open(configs.Default.Db.Driver, getDsnString(false))
	if err != nil {
		glog.Errorf("orm failed to initialized: %v", err)
	}
	defer db.Close()
	return db.Exec(fmt.Sprintf("Drop Database if exists %s;", name)).Error
}
func CreateDB(name string) {
	db, err := gorm.Open(configs.Default.Db.Driver, getDsnString(false))
	if err != nil {
		glog.Errorf("orm failed to initialized: %v", err)
	}
	script := ""
	if configs.Default.Db.Driver == Driver_MSSQL {
		script = fmt.Sprintf("if not exists (select * from sysdatabases where name='%s') begin create database %s end;", name, name)
	} else {
		script = fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET %s COLLATE %s;", name, configs.Default.Db.Charset, configs.Default.Db.Collation)
	}
	err = db.Exec(script).Error
	if err != nil {
		glog.Errorf("create DATABASE err: %v", err)
	}

	defer db.Close()
}
