package godror

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/godror/godror"
	"reflect"
	"strings"
	"time"

	"github.com/ggoop/goutils/gorm"
	"github.com/ggoop/goutils/utils"
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
	return fmt.Sprintf(`"%s"`, strings.ToUpper(key))
}
func (s *godror) DataTypeOf(field *gorm.StructField) string {
	var dataValue, sqlType, size, additionalType = gorm.ParseFieldStructForDialect(field, s)

	if sqlType == "" {
		switch dataValue.Kind() {
		case reflect.Bool:
			sqlType = "CHAR(1)"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
			sqlType = "INTEGER"
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
	s.db.QueryRow("SELECT count(*) FROM user_indexes WHERE index_name = :1 AND table_name =:2", strings.ToUpper(indexName), strings.ToUpper(tableName)).Scan(&count)
	return count > 0
}

func (s godror) RemoveIndex(tableName string, indexName string) error {
	_, err := s.db.Exec(fmt.Sprintf("DROP INDEX %v ON %v", indexName, s.Quote(tableName)))
	return err
}

func (s godror) HasForeignKey(tableName string, foreignKeyName string) bool {
	var count int
	s.db.QueryRow("SELECT COUNT(*) FROM USER_CONSTRAINTS WHERE CONSTRAINT_TYPE = 'R' AND TABLE_NAME = :1 AND CONSTRAINT_NAME = :2", strings.ToUpper(tableName), strings.ToUpper(foreignKeyName)).Scan(&count)
	return count > 0
}

func (s godror) HasTable(tableName string) bool {
	var count int
	s.db.QueryRow("SELECT COUNT(*) FROM USER_TABLES WHERE TABLE_NAME = :1", strings.ToUpper(tableName)).Scan(&count)
	return count > 0
}

func (s godror) HasColumn(tableName string, columnName string) bool {
	var count int
	s.db.QueryRow("SELECT COUNT(*) FROM USER_TAB_COLUMNS WHERE TABLE_NAME = :1 AND COLUMN_NAME = :2", strings.ToUpper(tableName), strings.ToUpper(columnName)).Scan(&count)
	return count > 0
}

func (s godror) ModifyColumn(tableName string, columnName string, typ string) error {
	_, err := s.db.Exec(fmt.Sprintf("ALTER TABLE %v ALTER COLUMN %v %v", strings.ToUpper(tableName), strings.ToUpper(columnName), typ))
	return err
}

func (s godror) CurrentDatabase() (name string) {
	s.db.QueryRow("SELECT DB_NAME() AS [Current Database]").Scan(&name)
	return
}

func (s godror) LimitAndOffsetSQL(limit, offset interface{}) (sql string, err error) {
	return
}

func (godror) AdjustSql(sql *utils.SqlStruct) *utils.SqlStruct {
	if sql.LimitNum > 0 {
		offset := sql.OffsetNum
		if offset < 0 {
			offset = 0
		}
		needPager := true
		if sql.LimitNum == 1 { // id查询不需要
			where := strings.ToLower(strings.ReplaceAll(sql.WhereSql, " ", ""))
			if strings.Contains(where, "(id=") {
				needPager = false
			}
		}
		if sql.OrderSql == "" && needPager { //无排序
			//SELECT *
			//FROM (SELECT ROWNUM AS rowno, t.*
			//          FROM DONORINFO t
			//          WHERE t.BIRTHDAY BETWEEN TO_DATE ('19800101', 'yyyymmdd')
			//          AND TO_DATE ('20060731', 'yyyymmdd')
			//          AND ROWNUM <= page*size) table_alias
			//WHERE table_alias.rowno > (page-1)*size;
			if sql.SelectSql == "*" && !strings.Contains(sql.FromSql, " ") {
				sql.FromSql = fmt.Sprintf("%s t", sql.FromSql)
				sql.SelectSql = "t.*"
			}
			sql.SelectSql = fmt.Sprintf("ROWNUM as row_num,%v", sql.SelectSql)
			sql.AddWhere(fmt.Sprintf("ROWNUM <= %v ", offset+sql.LimitNum))
			sqlStr := sql.CombinedSql()

			sql.SelectSql = "*"
			sql.FromSql = fmt.Sprintf("(%s) table_pager", sqlStr)
			sql.JoinSql = ""
			sql.GroupSql = ""
			sql.HavingSql = ""
			sql.OrderSql = ""
			sql.WhereSql = fmt.Sprintf("where table_pager.row_num> %v", sql.OffsetNum)

			sqlStr = sql.CombinedSql()
		} else if sql.OrderSql != "" && needPager { //有排序
			//SELECT *
			//FROM (SELECT ROWNUM AS rowno,r.*
			//           FROM(SELECT * FROM DONORINFO t
			//                    WHERE t.BIRTHDAY BETWEEN TO_DATE ('19800101', 'yyyymmdd')
			//                    AND TO_DATE ('20060731', 'yyyymmdd')
			//                    ORDER BY t.BIRTHDAY desc
			//                   ) r
			//           where ROWNUM <= page*size
			//          ) table_alias
			//WHERE table_alias.rowno > (page-1)*size;

			if sql.SelectSql == "*" && !strings.Contains(sql.FromSql, " ") {
				sql.FromSql = fmt.Sprintf("%s t", sql.FromSql)
				sql.SelectSql = "t.*"
			}

			sqlStr := sql.CombinedSql()

			sql.SelectSql = "*"
			sql.FromSql = fmt.Sprintf("(select ROWNUM as row_num,r.* from  (%s)  r where ROWNUM <= %v) table_pager", sqlStr, offset+sql.LimitNum)
			sql.JoinSql = ""
			sql.GroupSql = ""
			sql.HavingSql = ""
			sql.OrderSql = ""
			sql.WhereSql = fmt.Sprintf("where table_pager.row_num> %v", sql.OffsetNum)
			sqlStr = sql.CombinedSql()
		}
	}
	return sql
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
