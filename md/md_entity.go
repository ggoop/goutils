package md

type MDEntity struct {
	ModelUnscoped
	Type      string // enum,entity
	Code      string
	Name      string
	FullName  string
	TableName string
	Memo      string
	Fields    []MDField `gorm:"foreignkey:EntityID"`

	cache map[string]MDField
}

func (s *MDEntity) MD() *Mder {
	return &Mder{ID: "01e8f30669c91c409d0a79d0b60bd6f0", Name: "实体"}
}
func (s *MDEntity) GetField(code string) *MDField {
	if s.cache == nil {
		s.cache = make(map[string]MDField)
	}
	if s.Fields != nil && len(s.Fields) > 0 && len(s.cache) == 0 {
		for _, v := range s.Fields {
			s.cache[v.Code] = v
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
	Type           *MDEntity `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Memo           string
}

func (s *MDField) MD() *Mder {
	return &Mder{ID: "01e8f30673ef1ad092b90553bb1bf432", Name: "属性"}
}
