package md

type MDPage struct {
	ID        string `gorm:"primary_key;size:50" json:"id"`
	CreatedAt Time   `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt *Time  `gorm:"name:更新时间" json:"updated_at"`
	Type      string `gorm:"size:50;name:类型"` //page，ref，app
	Domain    string `gorm:"size:50;name:领域" json:"domain"`
	EntID     string `gorm:"size:50;name:企业"`
	Name      string `gorm:"size:50"`
	Widgets   string `gorm:"type:text"` //JSON
	Tags      string `gorm:"size:500"`
	Templated bool
	System    bool
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
	Data       string `gorm:"type:text"` //JSON
	Multiple   bool   `json:"multiple"`
	Nullable   bool
	IsMain     bool
}

func (s *MDPageView) MD() *Mder {
	return &Mder{ID: "md.page.entity", Domain: md_domain, Name: "页面视图"}
}

type MDPageWidget struct {
	ID        string `gorm:"primary_key;size:50" json:"id"`
	CreatedAt Time   `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt *Time  `gorm:"name:更新时间" json:"updated_at"`
	PageID    string `gorm:"size:50;name:页面"`
	ParentID  string `gorm:"size:50;name:上线ID"`
	PartID    string `gorm:"size:50;name:部件ID"`
	Sequence  int    `gorm:"name:排序" json:"sequence"`
	Title     string `gorm:"name:元素的额外信息"`
	Text      string `gorm:"name:文本内容"`
	Align     string `gorm:"name:对齐"`
	Color     string `gorm:"name:颜色"`
	BgColor   string `gorm:"name:背景颜色"`
	Classes   string `gorm:"name:元素的类名"`
	Styles    string `gorm:"name:元素的行内样式"`
	Props     string `gorm:"name:属性集合"`
}

func (s *MDPageWidget) MD() *Mder {
	return &Mder{ID: "md.page.widget", Domain: md_domain, Name: "页面组件"}
}
