package gorm

import (
	"fmt"
	"time"
)

// Model base model definition, including fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`, which could be embedded in your models
//    type User struct {
//      gorm.Model
//    }
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

const (
	//mssql
	DRIVER_MSSQL = "mssql"
	//mysql
	DRIVER_MYSQL = "mysql"
	//oracle :godror
	DRIVER_GODROR = "godror"
)

type SqlStruct struct {
	SelectSql string
	FromSql   string
	JoinSql   string
	GroupSql  string
	HavingSql string
	OrderSql  string
	LimitSql  string
	WhereSql  string
	LimitNum  int64
	OffsetNum int64
}

func (sql SqlStruct) CombinedConditionSql() string {
	return sql.JoinSql + sql.WhereSql + sql.GroupSql + sql.HavingSql + sql.OrderSql + sql.LimitSql
}
func (sql SqlStruct) CombinedSql() string {
	return fmt.Sprintf("SELECT %v FROM %v %v", sql.SelectSql, sql.FromSql, sql.CombinedConditionSql())
}
