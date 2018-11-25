package md

type Entity struct {
	ModelUnscoped
	Type      string // enum,entity
	Code      string
	Name      string
	TableName string
	Memo      string
}

func (s *Entity) MDID() string {
	return "01e8f0b45e12835fe7fd8cec4b7174de"
}

type EntityField struct {
	ModelUnscoped
	EntityID   string
	Entity     *Entity `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Code       string
	Name       string
	FieldName  string
	ForeignKey string
	LocalKey   string
	TypeID     string
	Type       *Entity `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Memo       string
}

func (s *EntityField) MDID() string {
	return "01e8f0b45e1456e1fc4d8cec4b7174de"
}
