package md

import "github.com/ggoop/goutils/utils"

type UIActionCommand struct {
	ID        string      `gorm:"primary_key;size:36" json:"id"`
	CreatedAt utils.Time  `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt utils.Time  `gorm:"name:更新时间" json:"updated_at"`
	EntID     string      `gorm:"size:36;name:企业" json:"ent_id"`
	WidgetID  string      `gorm:"size:36;name:组件ID" json:"widget_id"`
	Code      string      `gorm:"size:36;name:编码" json:"code"`
	Name      string      `gorm:"size:50;name:名称" json:"name"`
	Type      string      `gorm:"size:20;name:类型" json:"type"`
	Action    string      `gorm:"size:50;name:动作" json:"action"`
	Url       string      `gorm:"size:100;name:服务路径" json:"url"`
	Parameter utils.SJson `gorm:"type:text;name:参数" json:"parameter"`
	Method    string      `gorm:"size:20;name:请求方式" json:"method"`
	Target    string      `gorm:"size:36;name:目标" json:"target"`
	Script    string      `gorm:"type:text;name:脚本" json:"script"`
	Rules     string      `gorm:"type:text;name:规则链" json:"rules"`
	System    utils.SBool `gorm:"name:系统的" json:"system"`
}

func (s *UIActionCommand) MD() *Mder {
	return &Mder{ID: "md.ui.action.command", Domain: md_domain, Name: "组件命令"}
}

type UIActionRule struct {
	ID        string      `gorm:"primary_key;size:50" json:"id"` //领域.规则：md.save，ui.save
	CreatedAt utils.Time  `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt utils.Time  `gorm:"name:更新时间" json:"updated_at"`
	Domain    string      `gorm:"size:50;name:领域" json:"domain"` //common为公共动作
	Code      string      `gorm:"size:50;name:编码" json:"code"`
	Name      string      `gorm:"size:50;name:名称" json:"name"`
	Async     utils.SBool `gorm:"name:异步的" json:"async"`
	System    utils.SBool `gorm:"name:系统的" json:"system"`
}

func (s *UIActionRule) MD() *Mder {
	return &Mder{ID: "md.ui.action.rule", Domain: md_domain, Name: "动作规则"}
}
