package godror

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/godror/godror"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ggoop/goutils/gorm"
)

func init() {
	gorm.RegisterDialect("godror", &godror{})
}

type godror struct {
	db gorm.SQLCommon
	gorm.DefaultForeignKeyNamer
}

func (godror) GetName() string {
	return "godror"
}

func (s *godror) SetDB(db gorm.SQLCommon) {
	s.db = db
}

func (godror) BindVar(i int) string {
	return fmt.Sprintf(":%d", i)
}

func (godror) Quote(key string) string {
	return fmt.Sprintf(`"%s"`, key)
}

func (s *godror) DataTypeOf(field *gorm.StructField) string {
	var dataValue, sqlType, size, additionalType = gorm.ParseFieldStructForDialect(field, s)

	if sqlType == "" {
		switch dataValue.Kind() {
		case reflect.Bool:
			sqlType = "NUMBER"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
			sqlType = "NUMBER"
		case reflect.Int64, reflect.Uint64:
			sqlType = "NUMBER"
		case reflect.Float32, reflect.Float64:
			sqlType = "NUMBER"
		case reflect.String:
			if size > 0 {
				sqlType = fmt.Sprintf("varchar2(%d)", size)
			} else {
				sqlType = "varchar2(50)"
			}
		case reflect.Struct:
			if _, ok := dataValue.Interface().(time.Time); ok {
				sqlType = "TIMESTAMP"
			}
		default:
			if gorm.IsByteArrayOrSlice(dataValue) {
				if size > 0 && size < 8000 {
					sqlType = fmt.Sprintf("varchar2(%d)", size)
				} else {
					sqlType = "CLOB"
				}
			}
		}
	}
	if sqlType == "" {
		panic(fmt.Sprintf("invalid sql type %s (%s) for godror", dataValue.Type().Name(), dataValue.Kind().String()))
	}
	//指定替换
	switch strings.ToUpper(sqlType) {
	case "TEXT":
		sqlType = "CLOB"
	}
	return sqlType
	if strings.TrimSpace(additionalType) == "" {
		return sqlType
	}
	return fmt.Sprintf("%v %v", sqlType, additionalType)
}

func (s godror) fieldCanAutoIncrement(field *gorm.StructField) bool {
	if value, ok := field.TagSettingsGet("AUTO_INCREMENT"); ok {
		return value != "FALSE"
	}
	return field.IsPrimaryKey
}

func (s godror) HasIndex(tableName string, indexName string) bool {
	var count int
	s.db.QueryRow("SELECT count(*) FROM user_indexes WHERE index_name=? AND table_name=?", strings.ToUpper(indexName), tableName).Scan(&count)
	return count > 0
}

func (s godror) RemoveIndex(tableName string, indexName string) error {
	_, err := s.db.Exec(fmt.Sprintf("DROP INDEX %v ON %v", indexName, s.Quote(tableName)))
	return err
}

func (s godror) HasForeignKey(tableName string, foreignKeyName string) bool {
	var count int
	s.db.QueryRow("SELECT COUNT(*) FROM USER_CONSTRAINTS WHERE CONSTRAINT_TYPE = 'R' AND TABLE_NAME = :1 AND CONSTRAINT_NAME = :2", tableName, foreignKeyName).Scan(&count)
	return count > 0
}

func (s godror) HasTable(tableName string) bool {
	var count int
	s.db.QueryRow("SELECT COUNT(*) FROM USER_TABLES WHERE TABLE_NAME = :1", tableName).Scan(&count)
	return count > 0
}

func (s godror) HasColumn(tableName string, columnName string) bool {
	var count int
	s.db.QueryRow("SELECT COUNT(*) FROM USER_TAB_COLUMNS WHERE TABLE_NAME = :1 AND COLUMN_NAME = :2", tableName, columnName).Scan(&count)
	return count > 0
}

func (s godror) ModifyColumn(tableName string, columnName string, typ string) error {
	_, err := s.db.Exec(fmt.Sprintf("ALTER TABLE %v ALTER COLUMN %v %v", tableName, columnName, typ))
	return err
}

func (s godror) CurrentDatabase() (name string) {
	s.db.QueryRow("SELECT DB_NAME() AS [Current Database]").Scan(&name)
	return
}

func (s godror) LimitAndOffsetSQL(limit, offset interface{}) (sql string, err error) {
	if limit != nil {
		if parsedLimit, err := strconv.ParseInt(fmt.Sprint(limit), 0, 0); err == nil && parsedLimit >= 0 {
			sql += fmt.Sprintf("ROWNUM <= %d", limit)
		}
	}
	return
}

func (godror) SelectFromDummyTable() string {
	return "FROM dual"
}

func (godror) LastInsertIDOutputInterstitial(tableName, columnName string, columns []string) string {
	return ""
}

func (godror) LastInsertIDReturningSuffix(tableName, columnName string) string {
	return ""
}

func (godror) DefaultValueStr() string {
	return "DEFAULT VALUES"
}

// NormalizeIndexAndColumn returns argument's index name and column name without doing anything
func (godror) NormalizeIndexAndColumn(indexName, columnName string) (string, string) {
	return indexName, columnName
}
func currentDatabaseAndTable(dialect gorm.Dialect, tableName string) (string, string) {
	if strings.Contains(tableName, ".") {
		splitStrings := strings.SplitN(tableName, ".", 2)
		return splitStrings[0], splitStrings[1]
	}
	return dialect.CurrentDatabase(), tableName
}

// JSON type to support easy handling of JSON data in character table fields
// using golang json.RawMessage for deferred decoding/encoding
type JSON struct {
	json.RawMessage
}

// Value get value of JSON
func (j JSON) Value() (driver.Value, error) {
	if len(j.RawMessage) == 0 {
		return nil, nil
	}
	return j.MarshalJSON()
}

// Scan scan value into JSON
func (j *JSON) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value (strcast):", value))
	}
	bytes := []byte(str)
	return json.Unmarshal(bytes, j)
}
