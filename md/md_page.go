package md

type MDPage struct {
	ID         string `gorm:"primary_key;size:50" json:"id"`
	CreatedAt  Time   `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt  *Time  `gorm:"name:更新时间" json:"updated_at"`
	Type       string `gorm:"size:50;name:类型" json:"type"` //page，ref，app
	Domain     string `gorm:"size:50;name:领域" json:"domain"`
	EntID      string `gorm:"size:50;name:企业" json:"ent_id"`
	Code       string `gorm:"size:50" json:"code"`
	Name       string `gorm:"size:50" json:"name"`
	Widgets    SJson  `gorm:"type:text" json:"widgets"` //JSON
	MainEntity string `gorm:"size:50;name:主实体" json:"main_entity"`
	Templated  SBool  `json:"templated"`
	System     SBool  `json:"system"`
}

func (s *MDPage) MD() *Mder {
	return &Mder{ID: "md.page", Domain: md_domain, Name: "页面"}
}
