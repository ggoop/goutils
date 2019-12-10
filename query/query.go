package query

import (
	"github.com/ggoop/goutils/md"
)

const md_domain string = "query"

type Query struct {
	md.Model
	ScopeID     string        `gorm:"size:50;name:范围" json:"scope_id"`
	ScopeType   string        `gorm:"size:50;name:范围" json:"scope_type"`
	Code        string        `gorm:"size:50;name:编码" json:"code"`
	Name        string        `gorm:"name:名称"  json:"name"`
	Type        string        `gorm:"size:50;name:查询类型"  json:"type"` //entity,service
	Entry       string        `gorm:"size:50;name:入口"  json:"entry"`
	Memo        string        `gorm:"name:备注"  json:"memo"`
	PageSize    int           `gorm:"default:30;name:每页显示记录数" json:"page_size"`
	Fields      []QueryField  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:QueryID;name:字段集合" json:"fields"`
	Columns     []QueryColumn `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:OwnerID;name:栏目集合" json:"columns"`
	Orders      []QueryOrder  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:OwnerID;name:排序集合" json:"orders"`
	Wheres      []QueryWhere  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:OwnerID;name:条件集合" json:"wheres"`
	ContextJson string        `gorm:"name:上下文" json:"context_json"` //上下文参数
}

func (s *Query) MD() *md.Mder {
	return &md.Mder{ID: "query.query", Domain: md_domain, Name: "查询"}
}

type QueryField struct {
	md.Model
	QueryID      string `gorm:"size:50;name:查询ID" json:"query_id"`
	Query        *Query `gorm:"name:查询" json:"case_id"`
	DataType     string `gorm:"size:50;name:字段类型" json:"data_type"` //sys.query.data.type，string,enum,entity,bool,datetime
	DataSourse   string `gorm:"name:数据来源" json:"data_sourse"`
	DefaultValue string `gorm:"name:默认值" json:"default_value"`
	Field        string `gorm:"size:50;name:字段" json:"field"`
	Group        string `gorm:"name:分组" json:"group"`
	Title        string `gorm:"name:显示名称" json:"title"`
	Memo         string `gorm:"name:备注" json:"memo"`
	Hidden       bool   `gorm:"name:隐藏" json:"hidden"`
	IsColumn     bool   `gorm:"name:栏目" json:"is_column"`
	IsOrder      bool   `gorm:"name:排序" json:"is_order"`
	IsWhere      bool   `gorm:"name:条件" json:"is_where"`
}

func (s *QueryField) MD() *md.Mder {
	return &md.Mder{ID: "query.field", Domain: md_domain, Name: "查询字段"}
}

type QueryColumn struct {
	md.Model
	OwnerID   string `gorm:"size:50;name:所属ID" json:"owner_id"`
	OwnerType string `gorm:"size:50;name:所属类型" json:"owner_type"`
	Field     string `gorm:"size:50;name:字段" json:"field"`
	Name      string `gorm:"size:50;name:栏目编码" json:"name"`
	Title     string `gorm:"name:显示名称" json:"title"`
	Sequence  int    `gorm:"name:顺序" json:"sequence"`
	Width     string `gorm:"name:宽度" json:"width"`
	Fixed     bool   `gorm:"name:固定" json:"fixed"`
	Hidden    bool   `gorm:"name:隐藏" json:"hidden"`
}

func (s *QueryColumn) MD() *md.Mder {
	return &md.Mder{ID: "query.column", Domain: md_domain, Name: "查询栏目"}
}

type QueryOrder struct {
	md.Model
	OwnerID   string `gorm:"size:50;name:所属ID" json:"owner_id"`
	OwnerType string `gorm:"size:50;name:所属类型" json:"owner_type"`
	Field     string `gorm:"size:50;name:字段"  json:"field"`
	Title     string `gorm:"name:显示名称" json:"title"`
	Order     string `gorm:"name:排序方式"  json:"order"`
	Fixed     bool   `gorm:"name:固定" json:"fixed"`
	Hidden    bool   `gorm:"name:隐藏" json:"hidden"`
	Sequence  int    `gorm:"name:顺序"  json:"sequence"`
}

func (s *QueryOrder) MD() *md.Mder {
	return &md.Mder{ID: "query.order", Domain: md_domain, Name: "查询排序"}
}

type QueryWhere struct {
	md.Model
	OwnerID    string       `gorm:"size:50;name:所属ID" json:"owner_id"`
	OwnerType  string       `gorm:"size:50;name:所属类型" json:"owner_type"`
	ParentID   string       `gorm:"size:50" json:"parent_id"`
	Logical    string       `gorm:"size:10;name:逻辑" json:"logical"` //and or
	Field      string       `gorm:"size:50;name:字段" json:"field"`
	Title      string       `gorm:"name:显示名称" json:"title"`
	DataType   string       `gorm:"size:50;name:字段类型" json:"data_type"` //sys.query.data.type，string,enum,entity,bool,datetime
	DataSourse string       `gorm:"name:数据来源" json:"data_sourse"`
	Operator   string       `gorm:"size:50;name:操作符号" json:"operator"`
	Value      string       `gorm:"name:值" json:"value"`
	Sequence   int          `gorm:"name:顺序" json:"sequence"`
	Fixed      bool         `gorm:"name:固定" json:"fixed"`
	Hidden     bool         `gorm:"name:隐藏" json:"hidden"`
	Children   []QueryWhere `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:ParentID;name:子条件" json:"children"`
}

func (s *QueryWhere) MD() *md.Mder {
	return &md.Mder{ID: "query.where", Domain: md_domain, Name: "查询条件"}
}
