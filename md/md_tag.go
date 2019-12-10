package md

type MDTag struct {
	ID        string `gorm:"primary_key;size:50" json:"id"`
	CreatedAt Time   `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt *Time  `gorm:"name:更新时间" json:"updated_at"`
	Name      string `gorm:"size:50"`
	ParentID  string `gorm:"size:50;name:父模型" json:"parent_id"`
}

func (s *MDTag) MD() *Mder {
	return &Mder{ID: "md.tag", Domain: md_domain, Name: "标签"}
}
