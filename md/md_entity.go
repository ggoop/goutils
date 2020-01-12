package md

import (
	"strings"

	"github.com/ggoop/goutils/di"
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/repositories"
)

const (
	md_domain string = "md"
)

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

var mdCache map[string]*MDEntity

func GetEntity(id string) *MDEntity {
	defer func() {
		if err := recover(); err != nil {
			glog.Error(err)
		}
	}()
	if mdCache == nil {
		mdCache = make(map[string]*MDEntity)
	}
	id = strings.ToLower(id)
	if v, ok := mdCache[id]; ok {
		return v
	}
	item := &MDEntity{}
	if err := di.Global.Invoke(func(db *repositories.MysqlRepo) {
		db.Preload("Fields").Order("id").Take(item, "id=?", id)
	}); err != nil {
		repositories.Default().Preload("Fields").Order("id").Take(item, "id=?", id)
	}
	if item.ID != "" {
		mdCache[strings.ToLower(item.ID)] = item
		if item.TableName != "" {
			mdCache[strings.ToLower(item.TableName)] = item
		}
		return item
	}
	return nil
}
