package md

type Entity struct {
	ModelUnscoped
	Type      string // enum,entity
	Code      string
	Name      string
	TableName string
	Memo      string
}

func (s *Entity) GetMD() string {
	return "6d28bb8009b711e78e9151efd7044098"
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

func (s *EntityField) GetMD() string {
	return "f532f4a009b611e78197f91f9f35de3a"
}
