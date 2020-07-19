package md

/**
动作
*/
type MDActionCommand struct {
	ID        string `gorm:"primary_key;size:50" json:"id"` //save,delete
	CreatedAt Time   `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt *Time  `gorm:"name:更新时间" json:"updated_at"`
	PageID    string `gorm:"size:50;name:页面" json:"page_id"` //common为公共动作
	Code      string `gorm:"size:50;name:编码" json:"code"`
	Name      string `gorm:"size:50;name:名称" json:"name"`
	Type      string `gorm:"size:50;name:标签" json:"type"` //ui,sv
	Path      string `gorm:"name:路径" json:"path"`
	Content   string `gorm:"type:text;name:规则内容" json:"content"` //js语法
	Rules     string `gorm:"type:text;name:规则链" json:"rules"`
	System    SBool  `gorm:"name:系统的" json:"system"`
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
	Code      string `gorm:"size:50;name:编码" json:"code"`
	Name      string `gorm:"size:50;name:名称" json:"name"`
	Async     SBool  `gorm:"name:异步的" json:"async"`
	System    SBool  `gorm:"name:系统的" json:"system"`
}

func (s *MDActionRule) MD() *Mder {
	return &Mder{ID: "md.action.rule", Domain: md_domain, Name: "动作规则"}
}
