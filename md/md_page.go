package md

type MDPage struct {
	ID         string `gorm:"primary_key;size:50" json:"id"`
	CreatedAt  Time   `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt  *Time  `gorm:"name:更新时间" json:"updated_at"`
	Type       string `gorm:"size:50;name:类型"` //page，ref，app
	Domain     string `gorm:"size:50;name:领域" json:"domain"`
	EntID      string `gorm:"size:50;name:企业"`
	Name       string `gorm:"size:50"`
	Widgets    SJson  `gorm:"type:text"` //JSON
	MainEntity string `gorm:"size:50;name:主实体"`
	Templated  SBool
	System     SBool
}

func (s *MDPage) MD() *Mder {
	return &Mder{ID: "md.page", Domain: md_domain, Name: "页面"}
}

type MDPageView struct {
	ID         string `gorm:"primary_key;size:50" json:"id"`
	CreatedAt  Time   `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt  *Time  `gorm:"name:更新时间" json:"updated_at"`
	PageID     string `gorm:"size:50;name:页面"`
	Code       string `gorm:"size:50;name:编码"`
	Name       string `gorm:"size:50"`
	EntityID   string `gorm:"size:50;name:实体ID"`
	PrimaryKey string `gorm:"size:50;name:主键"`
	Data       SJson  `gorm:"type:text"` //JSON
	Multiple   SBool  `json:"multiple"`
	Nullable   SBool
	IsMain     SBool
}
