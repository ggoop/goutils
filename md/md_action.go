package md

/**
动作
*/
type MDActionCommand struct {
	ID        string `gorm:"primary_key;size:50" json:"id"` //save,delete
	CreatedAt Time   `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt *Time  `gorm:"name:更新时间" json:"updated_at"`
	Domain    string `gorm:"size:50;name:领域" json:"domain"` //common为公共动作
	Name      string `gorm:"name:名称"`
	Tags      string `gorm:"size:50;name:标签"` //ui,sv
	Path      string `gorm:"name:路径"`
	System    bool   `gorm:"name:系统的"`
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
	Name      string `gorm:"name:名称"`
	Content   string `gorm:"type:text;name:规则内容"` //js语法
	Tags      string `gorm:"size:50;name:标签"`     //ui,sv
	Async     bool   `gorm:"name:异步的"`
	System    bool   `gorm:"name:系统的"`
}

func (s *MDActionRule) MD() *Mder {
	return &Mder{ID: "md.action.rule", Domain: md_domain, Name: "动作规则"}
}

/**
规则
*/
type MDActionMaker struct {
	ID          string `gorm:"primary_key;size:50" json:"id"` //领域.规则：md.save，ui.save
	CreatedAt   Time   `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt   *Time  `gorm:"name:更新时间" json:"updated_at"`
	MakerID     string `gorm:"size:50;name:使用者ID" json:"maker_id"`   //pageID，EntityID
	MakerType   string `gorm:"size:50;name:使用者类型" json:"maker_type"` // page|entity
	CommandID   string `gorm:"size:50;name:动作ID" json:"command_id"`  //可空，表示公共规则
	RuleID      string `gorm:"size:50"`
	RuleContent string `gorm:"type:text;name:规则内容"` //js语法
	Sequence    int    `gorm:"name:排序" json:"sequence"`
	Group       int    `gorm:"name:分组" json:"group"`
}

func (s *MDActionMaker) MD() *Mder {
	return &Mder{ID: "md.action.maker", Domain: md_domain, Name: "动作规则"}
}
