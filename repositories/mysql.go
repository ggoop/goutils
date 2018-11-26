package repositories

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ggoop/goutils/configs"
	"github.com/ggoop/goutils/glog"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type MysqlRepo struct {
	*gorm.DB
}

func NewMysqlRepo() *MysqlRepo {
	// 生成数据库
	initDb()

	// 创建连接
	config := mysql.Config{User: configs.Default.Db.Username, Passwd: configs.Default.Db.Password, Net: "tcp", Addr: configs.Default.Db.Host, DBName: configs.Default.Db.Database, AllowNativePasswords: true, ParseTime: true}
	if configs.Default.Db.Port != "" {
		config.Addr = fmt.Sprintf("%s:%s", configs.Default.Db.Host, configs.Default.Db.Port)
	}
	str := config.FormatDSN()
	db, err := gorm.Open("mysql", str)
	if err != nil {
		glog.Errorf("orm failed to initialized: %v", err)
	}

	db.LogMode(configs.Default.App.Debug)
	return &MysqlRepo{db}
}
func initDb() {
	config := mysql.Config{User: configs.Default.Db.Username, Passwd: configs.Default.Db.Password, Net: "tcp", Addr: configs.Default.Db.Host, AllowNativePasswords: true, ParseTime: true}
	if configs.Default.Db.Port != "" {
		config.Addr = fmt.Sprintf("%s:%s", configs.Default.Db.Host, configs.Default.Db.Port)
	}
	str := config.FormatDSN()
	db, err := gorm.Open("mysql", str)
	if err != nil {
		glog.Errorf("orm failed to initialized: %v", err)
	}
	db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET %s COLLATE %s;", configs.Default.Db.Database, configs.Default.Db.Charset, configs.Default.Db.Collation))

	defer db.Close()
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
		if (mainField.IsPrimaryKey && mainField.IsBlank) || (mainField.IsIgnored) ||
			(mainField.Field.Kind() == reflect.Struct && mainField.Field.Type().Name() != "Time") ||
			(mainField.Field.Kind() == reflect.Slice && mainField.Field.Type().Elem().Kind() == reflect.Struct) {
			continue
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

			if (field.IsPrimaryKey && field.IsBlank) || (field.IsIgnored) ||
				(field.Field.Kind() == reflect.Struct && field.Field.Type().Name() != "Time") ||
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
			return err
		}
	}
	return nil
}
