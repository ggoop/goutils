package md

import (
	"strings"
)

const (
	md_domain string = "md"
)

type MDEnum struct {
	ID       string `gorm:"primary_key;size:50" json:"id"`
	EntityID string `gorm:"size:50;unique_index:uix" json:"entity_id"`
	Code     string `gorm:"size:50;unique_index:uix" json:"code"`
	Name     string `gorm:"size:50" json:"name"`
	Sequence int    `json:"sequence"`
}

func (t MDEnum) TableName() string {
	return "md_fields"
}
func (s *MDEnum) MD() *Mder {
	return &Mder{ID: "md.enum", Domain: md_domain, Name: "枚举", Type: TYPE_ENUM}
}

type MDEntity struct {
	ID        string `gorm:"primary_key;size:50" json:"id"`
	CreatedAt Time   `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt *Time  `gorm:"name:更新时间" json:"updated_at"`
	Type      string `gorm:"size:50"` // simple，entity，enum，interface，dto,view
	Domain    string `gorm:"size:50;name:领域" json:"domain"`
	Code      string `gorm:"size:100"`
	Name      string `gorm:"size:100"`
	TableName string `gorm:"size:50"`
	Memo      string `gorm:"size:500"`
	Tags      string `gorm:"size:500"`
	System    bool
	Fields    []MDField `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:EntityID"`
	cache     map[string]MDField
}

func (s *MDEntity) MD() *Mder {
	return &Mder{ID: "md.entity", Domain: md_domain, Name: "实体"}
}
func (s *MDEntity) GetField(code string) *MDField {
	if s.cache == nil {
		s.cache = make(map[string]MDField)
	}
	code = strings.ToLower(code)
	if s.Fields != nil && len(s.Fields) > 0 && len(s.cache) == 0 {
		for _, v := range s.Fields {
			s.cache[strings.ToLower(v.Code)] = v
		}
	}
	if v, ok := s.cache[code]; ok {
		return &v
	}
	return nil
}

type MDEntityRelation struct {
	ID        string `gorm:"primary_key;size:50" json:"id"`
	CreatedAt Time   `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt *Time  `gorm:"name:更新时间" json:"updated_at"`
	ParentID  string `gorm:"size:50;name:页面"`
	ChildID   string `gorm:"size:50;name:动作ID"`
	Kind      string `gorm:"name:参数"` //inherit，interface，
}

func (s *MDEntityRelation) MD() *Mder {
	return &Mder{ID: "md.entity.relation", Domain: md_domain, Name: "实体关系"}
}

type MDField struct {
	ID             string    `gorm:"primary_key;size:50" json:"id"`
	CreatedAt      Time      `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt      *Time     `gorm:"name:更新时间" json:"updated_at"`
	EntityID       string    `gorm:"size:50;unique_index:uix"`
	Entity         *MDEntity `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Code           string    `gorm:"size:50;unique_index:uix"`
	Name           string    `gorm:"size:50"`
	DbName         string    `gorm:"size:50"`
	IsNormal       bool
	IsPrimaryKey   bool
	ForeignKey     string    `gorm:"size:50"` //外键
	AssociationKey string    `gorm:"size:50"` //Association
	Kind           string    `gorm:"size:50"`
	TypeID         string    `gorm:"size:50"`
	TypeType       string    `gorm:"size:50"`
	Type           *MDEntity `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Limit          string    `gorm:"size:500;name:限制"`
	Memo           string    `gorm:"size:500"`
	Tags           string    `gorm:"size:500"`
	Sequence       int
	Nullable       bool
	Length         int
	Precision      int
	DefaultValue   string
	MinValue       string
	MaxValue       string
}

func (s *MDField) MD() *Mder {
	return &Mder{ID: "md.field", Domain: md_domain, Name: "属性"}
}
