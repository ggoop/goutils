package md

type Entity struct {
	ModelUnscoped
	Type      string // enum,entity
	Code      string
	Name      string
	TableName string
	Memo      string

	Fields []EntityField

	cache map[string]EntityField
}

func (s *Entity) MD() *Mder {
	return &Mder{ID: "01e8f0b45e12835fe7fd8cec4b7174de", Name: "实体"}
}
func (s *Entity) GetField(name string) *EntityField {
	if s.cache == nil {
		s.cache = make(map[string]EntityField)
	}
	if s.Fields != nil && len(s.Fields) > 0 && len(s.cache) == 0 {
		for _, v := range s.Fields {
			s.cache[v.Name] = v
		}
	}
	if v, ok := s.cache[name]; ok {
		return &v
	}
	return nil
}

type EntityField struct {
	ModelUnscoped
	EntityID       string
	Entity         *Entity `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Code           string
	Name           string
	DBName         string
	IsNormal       bool
	IsPrimaryKey   bool
	ForeignKey     string //外键
	AssociationKey string //Association ForeignKey
	Kind           string
	TypeID         string
	Type           *Entity `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Memo           string
}

func (s *EntityField) MD() *Mder {
	return &Mder{ID: "01e8f0b45e1456e1fc4d8cec4b7174de", Name: "属性"}
}
