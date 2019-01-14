package query

import (
	"github.com/ggoop/goutils/md"
)

type Query struct {
	md.ModelUnscoped
	Code    string        `gorm:"size:100;name:编码" json:"code"`
	Name    string        `gorm:"name:名称"  json:"name"`
	Type    string        `gorm:"size:100;name:查询类型"  json:"type"` //entity,service
	Entry   string        `gorm:"size:100;name:入口"  json:"entry"`
	Memo    string        `gorm:"name:备注"  json:"memo"`
	Fields  []QueryField  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:QueryID;name:字段集合" json:"fields"`
	Columns []QueryColumn `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:QueryID;name:栏目集合" json:"columns"`
	Orders  []QueryOrder  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:QueryID;name:排序集合" json:"orders"`
	Wheres  []QueryWhere  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:QueryID;name:条件集合" json:"wheres"`
}

func (s *Query) MD() *md.Mder {
	return &md.Mder{ID: "01e8f3067a691a50b46b697fa9f73d01", Name: "查询"}
}

type QueryField struct {
	md.ModelUnscoped
	QueryID string `gorm:"name:查询ID" json:"query_id"`
	Query   *Query `gorm:"name:查询" json:"case_id"`
	Type    string `gorm:"size:100;name:字段类型" json:"type"` //string,enum,entity,bool,datetime
	Field   string `gorm:"size:100;name:字段" json:"field"`
	Name    string `gorm:"name:名称" json:"name"`
	Memo    string `gorm:"name:备注" json:"memo"`
	Hidden  bool   `gorm:"name:隐藏" json:"hidden"`
}

func (s *QueryField) MD() *md.Mder {
	return &md.Mder{ID: "01e8f30683641760aa9261a2b248c5f0", Name: "查询字段"}
}

type QueryColumn struct {
	md.ModelUnscoped
	QueryID    string `gorm:"name:查询ID" json:"query_id"`
	CaseID     string `gorm:"name:方案ID" json:"case_id"`
	Field      string `gorm:"size:100;name:字段" json:"field"`
	ColumnName string `gorm:"size:100;name:栏目名称" json:"column_name"`
	Title      string `gorm:"name:显示名称" json:"title"`
	Sequence   int    `gorm:"name:顺序" json:"sequence"`
	Width      string `gorm:"name:宽度" json:"width"`
}

func (s *QueryColumn) MD() *md.Mder {
	return &md.Mder{ID: "01e916da3fb0b00455b78cec4b7174de", Name: "查询栏目"}
}

type QueryOrder struct {
	md.ModelUnscoped
	QueryID  string `gorm:"name:查询ID" json:"query_id"`
	CaseID   string `gorm:"name:方案ID" json:"case_id"`
	Field    string `gorm:"size:100;name:字段"  json:"field"`
	IsDesc   bool   `gorm:"name:降序"  json:"is_desc"`
	Sequence int    `gorm:"name:顺序"  json:"sequence"`
}

func (s *QueryOrder) MD() *md.Mder {
	return &md.Mder{ID: "01e916da3fa5e1f4c3118cec4b7174de", Name: "查询排序"}
}

type QueryWhere struct {
	md.ModelUnscoped
	QueryID  string       `gorm:"name:查询ID" json:"query_id"`
	CaseID   string       `gorm:"name:方案ID" json:"case_id"`
	ParentID string       `gorm:"size:100" json:"parent_id"`
	Field    string       `gorm:"size:100;name:字段" json:"field"`
	Type     string       `gorm:"size:100;name:类型" json:"type"`
	Operator string       `gorm:"size:100;name:操作符号" json:"operator"`
	Value    string       `gorm:"name:值" json:"value"`
	Sequence int          `gorm:"name:顺序" json:"sequence"`
	Fixed    bool         `gorm:"name:固定" json:"fixed"`
	Hidden   bool         `gorm:"name:隐藏" json:"hidden"`
	Children []QueryWhere `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:ParentID;name:子条件" json:"children"`
}

func (s *QueryWhere) MD() *md.Mder {
	return &md.Mder{ID: "01e916da3fb51af46f288cec4b7174de", Name: "查询条件"}
}
