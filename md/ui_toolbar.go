package md

import "github.com/ggoop/goutils/utils"

type UIToolbars struct {
	ID        string          `gorm:"primary_key;size:36" json:"id"`
	CreatedAt utils.Time      `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt utils.Time      `gorm:"name:更新时间" json:"updated_at"`
	EntID     string          `gorm:"size:36;name:企业" json:"ent_id"`
	WidgetID  string          `gorm:"size:36;name:组件ID" json:"widget_id"`
	LayoutID  string          `gorm:"size:36;name:布局ID" json:"layout_id"`
	Code      string          `gorm:"size:36;name:编码" json:"code"`
	Name      string          `gorm:"size:50;name:名称" json:"name"`
	Type      string          `gorm:"size:20;name:类型" json:"type"`
	Mount     string          `gorm:"size:20;name:加载方式" json:"mount"`
	Sequence  int             `gorm:"size:3;name:顺序" json:"sequence"`
	Align     string          `gorm:"size:20;name:对齐方式" json:"align"`
	Style     utils.SJson     `gorm:"type:text;name:样式" json:"style"`
	System    utils.SBool     `gorm:"name:系统的" json:"system"`
	Items     []UIToolbarItem `gorm:"-" json:"items"`
}

func (s *UIToolbars) MD() *Mder {
	return &Mder{ID: "md.ui.toolbars", Domain: md_domain, Name: "组件工具集"}
}

type UIToolbarItem struct {
	ID        string          `gorm:"primary_key;size:36" json:"id"`
	CreatedAt utils.Time      `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt utils.Time      `gorm:"name:更新时间" json:"updated_at"`
	EntID     string          `gorm:"size:36;name:企业" json:"ent_id"`
	WidgetID  string          `gorm:"size:36;name:组件ID" json:"widget_id"`
	LayoutID  string          `gorm:"size:36;name:布局ID" json:"layout_id"`
	ParentID  string          `gorm:"size:36;name:上级" json:"parent_id"`
	Children  []UIToolbarItem `gorm:"-" json:"children"`
	Code      string          `gorm:"size:36;name:编码" json:"code"`
	Name      string          `gorm:"size:50;name:名称" json:"name"`
	Type      string          `gorm:"size:20;name:类型" json:"type"`
	Caption   string          `gorm:"size:50;name:标题" json:"caption"`
	Command   string          `gorm:"size:36;name:命令" json:"command"`
	Parameter utils.SJson     `gorm:"type:text;name:参数" json:"parameter"`
	Sequence  int             `gorm:"size:3;name:顺序" json:"sequence"`
	Icon      string          `gorm:"size:100;name:图标" json:"icon"`
	Align     string          `gorm:"size:20;name:对齐方式" json:"align"`
	Style     utils.SJson     `gorm:"type:text;name:样式" json:"style"`
	AuthID    string          `gorm:"size:36;name:权限ID" json:"auth_id"`
	Extras    utils.SJson     `gorm:"type:text;name:扩展属性" json:"extras"` //JSON
	System    utils.SBool     `gorm:"name:系统的" json:"system"`
}

func (s *UIToolbarItem) MD() *Mder {
	return &Mder{ID: "md.ui.toolbar.item", Domain: md_domain, Name: "组件工具条"}
}
