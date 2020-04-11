package md

/**
动作
*/
type MDActionCommand struct {
	ID        string `gorm:"primary_key;size:50" json:"id"` //save,delete
	CreatedAt Time   `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt *Time  `gorm:"name:更新时间" json:"updated_at"`
	PageID    string `gorm:"size:50;name:页面"` //common为公共动作
	Code      string `gorm:"size:50;name:编码"`
	Name      string `gorm:"size:50;name:名称"`
	Type      string `gorm:"size:50;name:标签"` //ui,sv
	Path      string `gorm:"name:路径"`
	Content   string `gorm:"type:text;name:规则内容"` //js语法
	Rules     string `gorm:"type:text;name:规则链"`
	System    SBool  `gorm:"name:系统的"`
}

func (s *MDActionCommand) MD() *Mder {
	return &Mder{ID: "md.action.command", Domain: md_domain, Name: "动作"}
}

/**
规则
*/
type MDActionRule struct {
	ID        string `gorm:"primary_key;size:50" json:"id"` //领域.规则：md.save，ui.save
	CreatedAt Time   `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt *Time  `gorm:"name:更新时间" json:"updated_at"`
	Domain    string `gorm:"size:50;name:领域" json:"domain"` //common为公共动作
	Code      string `gorm:"size:50;name:编码"`
	Name      string `gorm:"size:50;name:名称"`
	Async     SBool  `gorm:"name:异步的"`
	System    SBool  `gorm:"name:系统的"`
}

func (s *MDActionRule) MD() *Mder {
	return &Mder{ID: "md.action.rule", Domain: md_domain, Name: "动作规则"}
}
