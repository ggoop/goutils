package md

import "strings"
const md_domain string ="md"
type MDEntity struct {
	ModelUnscoped
	TypeID    string  // enum,entity
	Type      *MDEnum `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;limit:md.data.type"`
	Domain    string `gorm:"size:50;name:领域" json:"domain"`
	Code      string
	Name      string
	FullName  string
	TableName string
	Memo      string
	Fields    []MDField `gorm:"foreignkey:EntityID"`
	cache     map[string]MDField
}

func (s *MDEntity) MD() *Mder {
	return &Mder{ID: "01e8f30669c91c409d0a79d0b60bd6f0",Domain:md_domain, Name: "实体"}
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

type MDField struct {
	ModelUnscoped
	EntityID       string
	Entity         *MDEntity `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Code           string
	Name           string
	DBName         string
	IsNormal       bool
	IsPrimaryKey   bool
	ForeignKey     string //外键
	AssociationKey string //Association ForeignKey
	Kind           string
	TypeID         string
	TypeType       string
	Limit          string
	Type           *MDEntity `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Memo           string
}

func (s *MDField) MD() *Mder {
	return &Mder{ID: "01e8f30673ef1ad092b90553bb1bf432",Domain:md_domain, Name: "属性"}
}
