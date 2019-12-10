package md

type MDPart struct {
	ID        string `gorm:"primary_key;size:50" json:"id"`
	CreatedAt Time   `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt *Time  `gorm:"name:更新时间" json:"updated_at"`
	Element   string `gorm:"size:50"` // layout,button,table,icon,pagination,checkbox,radio,switch,select,avatar,card,list,divider
	Domain    string `gorm:"size:50;name:领域" json:"domain"`
	Name      string `gorm:"size:50"`
	Content   string `gorm:"name:内容"`
	Type      string `gorm:"size:50"`
	Classes   string `gorm:"name:元素的类名"`
	Styles    string `gorm:"name:元素的行内样式"`
	System    bool
}

func (s *MDPart) MD() *Mder {
	return &Mder{ID: "md.part", Domain: md_domain, Name: "部件"}
}

type MDPartProps struct {
	ID           string `gorm:"primary_key;size:50" json:"id"`
	CreatedAt    Time   `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt    *Time  `gorm:"name:更新时间" json:"updated_at"`
	PartID       string `gorm:"size:50;name:部件"`
	Code         string `gorm:"size:50;name:编码"`
	Name         string `gorm:"size:50"`
	Type         string `gorm:"size:50"` //bool,string,color,function,part
	Kind         string `gorm:"size:50"` //prop,state,method
	Nullable     bool
	DefaultValue string
	MinValue     string
	MaxValue     string
}

func (s *MDPartProps) MD() *Mder {
	return &Mder{ID: "md.part.prop", Domain: md_domain, Name: "部件属性"}
}
