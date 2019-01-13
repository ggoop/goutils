package query

import (
	"github.com/ggoop/goutils/md"
)

type Query struct {
	md.ModelUnscoped
	Code    string        `gorm:"size:100;name:编码"  json:"code"`
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
	QueryID string `gorm:"name:查询ID"`
	Query   *Query `gorm:"name:查询"`
	Type    string `gorm:"size:100;name:字段类型"` //string,enum,entity,bool,datetime
	Field   string `gorm:"size:100;name:字段"`
	Name    string `gorm:"name:名称"`
	Memo    string `gorm:"name:备注"`
	Hidden  bool   `gorm:"name:隐藏"`
}

func (s *QueryField) MD() *md.Mder {
	return &md.Mder{ID: "01e8f30683641760aa9261a2b248c5f0", Name: "查询字段"}
}

type QueryColumn struct {
	md.ModelUnscoped
	QueryID  string `gorm:"name:查询ID"`
	CaseID   string `gorm:"name:方案ID"`
	Field    string `gorm:"size:100;name:字段"`
	Title    string `gorm:"name:显示名称"`
	Sequence int    `gorm:"name:顺序"`
	Width    int    `gorm:"name:宽度"`
}

func (s *QueryColumn) MD() *md.Mder {
	return &md.Mder{ID: "01e916da3fb0b00455b78cec4b7174de", Name: "查询栏目"}
}

type QueryOrder struct {
	md.ModelUnscoped
	QueryID  string `gorm:"name:查询ID"`
	CaseID   string `gorm:"name:方案ID"`
	Field    string `gorm:"size:100;name:字段"`
	IsDesc   bool   `gorm:"name:降序"`
	Sequence int    `gorm:"name:顺序"`
}

func (s *QueryOrder) MD() *md.Mder {
	return &md.Mder{ID: "01e916da3fa5e1f4c3118cec4b7174de", Name: "查询排序"}
}

type QueryWhere struct {
	md.ModelUnscoped
	QueryID  string `gorm:"name:查询ID"`
	CaseID   string `gorm:"name:方案ID"`
	Field    string `gorm:"size:100;name:字段"`
	Type     string `gorm:"size:100;name:类型"`
	Operator string `gorm:"size:100;name:操作符号"`
	Value    string `gorm:"name:值"`
	Sequence int    `gorm:"name:顺序"`
	Fixed    bool   `gorm:"name:固定"`
	Hidden   bool   `gorm:"name:隐藏"`
}

func (s *QueryWhere) MD() *md.Mder {
	return &md.Mder{ID: "01e916da3fb51af46f288cec4b7174de", Name: "查询条件"}
}

type QueryCase struct {
	md.ModelUnscoped
	EntID     string        `gorm:"size:100"`
	UserID    string        `gorm:"size:100"`
	QueryID   string        `gorm:"name:查询ID"`
	Name      string        `gorm:"name:名称"`
	Content   string        `gorm:"name:内容"`
	ScopeType string        `gorm:"name:范围类型"`
	ScopeID   string        `gorm:"name:范围ID"`
	Memo      string        `gorm:"name:备注"`
	Columns   []QueryColumn `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:CaseID;name:栏目集合" json:"columns"`
	Orders    []QueryOrder  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:CaseID;name:排序集合" json:"orders"`
	Wheres    []QueryWhere  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:CaseID;name:条件集合" json:"wheres"`
}

func (s *QueryCase) MD() *md.Mder {
	return &md.Mder{ID: "01e916da3fbb092f44be8cec4b7174de", Name: "查询方案"}
}
