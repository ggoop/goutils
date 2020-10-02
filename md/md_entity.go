package md

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strings"

	"github.com/ggoop/goutils/di"
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/repositories"
	"github.com/ggoop/goutils/utils"
)

const (
	md_domain string = "md"
)

type MDEntity struct {
	ID        string     `gorm:"primary_key;size:50" json:"id"`
	CreatedAt utils.Time `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt utils.Time `gorm:"name:更新时间" json:"updated_at"`
	Type      string     `gorm:"size:50"` // simple，entity，enum，interface，dto,view
	Domain    string     `gorm:"size:50;name:领域" json:"domain"`
	Code      string     `gorm:"size:100;index:code_idx"`
	Name      string     `gorm:"size:100"`
	TableName string     `gorm:"size:50"`
	Memo      string     `gorm:"size:500"`
	Tags      string     `gorm:"size:500"`
	System    utils.SBool
	Fields    []MDField `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:EntityID"`
	cache     map[string]MDField
}

func (s MDEntity) String() string {
	return fmt.Sprintf("%s-%s-%s", s.Domain, s.Code, s.ID)
}
func (s *MDEntity) MD() *Mder {
	return &Mder{ID: "md.entity", Domain: md_domain, Name: "实体"}
}
func (s *MDEntity) GetField(code string) *MDField {
	if s.cache == nil {
		s.cache = make(map[string]MDField)
	}

	if s.Fields != nil && len(s.Fields) > 0 && len(s.cache) == 0 {
		for i, v := range s.Fields {
			s.cache[strings.ToLower(v.Code)] = s.Fields[i]
			if v.DbName != "" {
				s.cache[strings.ToLower(v.DbName)] = s.Fields[i]
			}
		}
	}
	if v, ok := s.cache[strings.ToLower(code)]; ok {
		return &v
	}
	return nil
}

type MDEntityRelation struct {
	ID        string     `gorm:"primary_key;size:50" json:"id"`
	CreatedAt utils.Time `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt utils.Time `gorm:"name:更新时间" json:"updated_at"`
	ParentID  string     `gorm:"size:50;name:页面"`
	ChildID   string     `gorm:"size:50;name:动作ID"`
	Kind      string     `gorm:"name:参数"` //inherit，interface，
}

func (s *MDEntityRelation) MD() *Mder {
	return &Mder{ID: "md.entity.relation", Domain: md_domain, Name: "实体关系"}
}

type MDField struct {
	ID             string     `gorm:"primary_key;size:50" json:"id"`
	CreatedAt      utils.Time `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt      utils.Time `gorm:"name:更新时间" json:"updated_at"`
	EntityID       string     `gorm:"size:50;unique_index:uix"`
	Entity         *MDEntity  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Code           string     `gorm:"size:50;unique_index:uix"`
	Name           string     `gorm:"size:50"`
	DbName         string     `gorm:"size:50"`
	IsNormal       utils.SBool
	IsPrimaryKey   utils.SBool
	ForeignKey     string    `gorm:"size:50"` //外键
	AssociationKey string    `gorm:"size:50"` //Association
	Kind           string    `gorm:"size:50"`
	TypeID         string    `gorm:"size:50"`
	TypeType       string    `gorm:"size:50"`
	Type           *MDEntity `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Limit          string    `gorm:"size:500;name:限制"`
	Memo           string    `gorm:"size:500"`
	Tags           string    `gorm:"size:500"` // code,name,ent,import
	Sequence       int
	Nullable       utils.SBool
	Length         int
	Precision      int
	DefaultValue   string
	MinValue       string
	MaxValue       string
}

func (s MDField) String() string {
	return fmt.Sprintf("%s-%s-%s", s.Code, s.Name, s.TypeType)
}
func (s MDField) CompileValue(value interface{}) interface{} {
	if value == nil || value == "" || s.TypeID == "" {
		return nil
	}
	if s.TypeID == FIELD_TYPE_STRING || s.TypeID == FIELD_TYPE_TEXT || s.TypeID == FIELD_TYPE_XML {
		return value
	}
	if s.TypeID == FIELD_TYPE_INT {
		return utils.ToInt(value)
	}
	if s.TypeID == FIELD_TYPE_BOOL {
		return utils.SBool_Parse(value)
	}
	if s.TypeID == FIELD_TYPE_JSON {
		return utils.SJson_Parse(value)
	}
	if s.TypeID == FIELD_TYPE_DATE || s.TypeID == FIELD_TYPE_DATETIME {
		return utils.CreateTime(value)
	}
	if s.TypeID == FIELD_TYPE_DECIMAL {
		if v, err := decimal.NewFromString(utils.ToString(value)); err != nil {
			return glog.Error(err)
		} else {
			return v
		}
	}
	return nil
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
	if v, ok := mdCache[strings.ToLower(id)]; ok {
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

type MDFilter struct {
	ID         string      `gorm:"primary_key;size:50" json:"id"`
	CreatedAt  utils.Time  `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt  utils.Time  `gorm:"name:更新时间" json:"updated_at"`
	OwnerID    string      `gorm:"size:50;index:ownerIdx;name:拥有者ID" json:"owner_id"`   //查询
	OwnerType  string      `gorm:"size:50;index:ownerIdx;name:拥有者类型" json:"owner_type"` //查询
	OwnerField string      `gorm:"size:50;index:ownerIdx;name:所属字段" json:"owner_field"` //字段名称
	Field      string      `gorm:"size:50;name:字段" json:"field"`                        //如果#开始，则表示标记，需要使用表达式
	Expr       string      `gorm:"size:200;name:表达式" json:"expr"`
	Title      string      `gorm:"name:显示名称" json:"title"`
	DataType   string      `gorm:"size:50;name:字段类型" json:"data_type"` //string,bool,datetime
	DataSource string      `gorm:"name:数据来源" json:"data_source"`       //条件：where| 常量：const| 变量：var|context:上下文
	Operator   string      `gorm:"size:50;name:操作符号" json:"operator"`
	Value      utils.SJson `gorm:"name:值" json:"value"`
	Enabled    utils.SBool `gorm:"default:true;name:启用" json:"enabled"`
}

func (s *MDFilter) MD() *Mder {
	return &Mder{ID: "md.filter", Domain: md_domain, Name: "限制条件"}
}
