package gorm

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/utils"
)

func (s *DB) BatchInsert(objArr []interface{}) error {
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
			if mainField.Name == "ID" && mainField.Field.Type().Kind() == reflect.String {

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
				if field.Name == "ID" && field.Field.Type().Kind() == reflect.String {
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
