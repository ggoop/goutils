package repositories

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ggoop/goutils/configs"
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/utils"

	"github.com/go-sql-driver/mysql"

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
func updateIDForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		if field, ok := scope.FieldByName("ID"); ok {
			if field.IsBlank && !field.IsIgnored {
				if field.Field.Type().Kind() == reflect.String {
					field.Set(utils.GUID())
				}
			}
		}
	}
}
func createMysqlRepo(db *gorm.DB) *MysqlRepo {
	repo := &MysqlRepo{db}
	repo.Callback().Create().Before("gorm:create").Register("md:update_id", updateIDForCreateCallback)
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

// Begin begin a transaction
func (s *MysqlRepo) Begin() *MysqlRepo {
	return createMysqlRepo(s.DB.Begin())
}

// New clone a new db connection without search conditions
func (s *MysqlRepo) New() *MysqlRepo {
	return createMysqlRepo(s.DB.New())
}
func (s *MysqlRepo) BatchInsert(objArr []interface{}) error {
	if len(objArr) == 0 {
		return nil
	}
	var itemCount uint = 0
	var MaxBatchs uint = 100

	mainObj := objArr[0]
	mainScope := s.NewScope(mainObj)
	mainFields := mainScope.Fields()
	quoted := make([]string, 0, len(mainFields))
	for i := range mainFields {
		mainField := mainFields[i]
		if !mainField.IsNormal || mainField.IsIgnored {
			continue
		}
		if (mainField.IsIgnored) || (mainField.Relationship != nil) ||
			(mainField.Field.Kind() == reflect.Slice && mainField.Field.Type().Elem().Kind() == reflect.Struct) {
			continue
		}
		if mainField.IsPrimaryKey && mainField.IsBlank {
			if (mainField.Name == "ID" && mainField.Field.Type().Kind() == reflect.String) {

			} else {
				continue
			}
		}
		quoted = append(quoted, mainScope.Quote(mainFields[i].DBName))
	}

	placeholdersArr := make([]string, 0, MaxBatchs)

	for _, obj := range objArr {
		itemCount++
		scope := s.NewScope(obj)
		fields := scope.Fields()
		placeholders := make([]string, 0, len(fields))
		for i := range fields {
			field := fields[i]
			if field.Name == "CreatedAt" && field.IsBlank {
				fields[i].Set(time.Now())
			}
			if field.Name == "UpdatedAt" && field.IsBlank {
				fields[i].Set(time.Now())
			}
			if field.IsPrimaryKey && field.IsBlank {
				if (field.Name == "ID" && field.Field.Type().Kind() == reflect.String) {
					fields[i].Set(utils.GUID())
				}
			}
			if !field.IsNormal || field.IsIgnored {
				continue
			}

			if (field.IsPrimaryKey && field.IsBlank) || (field.IsIgnored) || (field.Relationship != nil) ||
				(field.Field.Kind() == reflect.Slice && field.Field.Type().Elem().Kind() == reflect.Struct) {
				continue
			}
			placeholders = append(placeholders, mainScope.AddToVars(fields[i].Field.Interface()))
		}
		placeholdersStr := "(" + strings.Join(placeholders, ", ") + ")"
		placeholdersArr = append(placeholdersArr, placeholdersStr)
		mainScope.SQLVars = append(mainScope.SQLVars, scope.SQLVars...)

		if itemCount >= MaxBatchs {
			mainScope.Raw(fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
				mainScope.QuotedTableName(),
				strings.Join(quoted, ", "),
				strings.Join(placeholdersArr, ", "),
			))
			_, err := mainScope.SQLDB().Exec(mainScope.SQL, mainScope.SQLVars...)

			itemCount = 0
			placeholdersArr = make([]string, 0, MaxBatchs)
			mainScope.SQLVars = make([]interface{}, 0)

			if err != nil {
				return err
			}
		}
	}
	if len(placeholdersArr) > 0 && itemCount > 0 {
		mainScope.Raw(fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
			mainScope.QuotedTableName(),
			strings.Join(quoted, ", "),
			strings.Join(placeholdersArr, ", "),
		))
		if _, err := mainScope.SQLDB().Exec(mainScope.SQL, mainScope.SQLVars...); err != nil {
			glog.Error(err)
			return err
		}
	}
	return nil
}
