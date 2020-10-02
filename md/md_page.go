package md

import "github.com/ggoop/goutils/utils"

type MDPage struct {
	ID         string         `gorm:"primary_key;size:50" json:"id"`
	CreatedAt  utils.Time     `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt  utils.Time     `gorm:"name:更新时间" json:"updated_at"`
	Type       string         `gorm:"size:50;name:类型" json:"type"` //page，ref，app
	Domain     string         `gorm:"size:50;name:领域" json:"domain"`
	EntID      string         `gorm:"size:50;name:企业" json:"ent_id"`
	Code       string         `gorm:"size:50;unique_index:uix_code" json:"code"`
	Name       string         `gorm:"size:50" json:"name"`
	Children   []MDPageWidget `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:PageID"`
	Widgets    utils.SJson    `gorm:"type:text" json:"widgets"` //JSON
	Extras     utils.SJson    `gorm:"type:text" json:"extras"`  //JSON
	MainEntity string         `gorm:"size:50;name:主实体" json:"main_entity"`
	Element    string         `gorm:"size:20;name:元素" json:"element"`
	System     utils.SBool    `json:"system"`
}

func (s *MDPage) MD() *Mder {
	return &Mder{ID: "md.page", Domain: md_domain, Name: "页面"}
}

type MDPageWidget struct {
	ID          string         `gorm:"primary_key;size:50" json:"id"`
	CreatedAt   utils.Time     `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt   utils.Time     `gorm:"name:更新时间" json:"updated_at"`
	EntID       string         `gorm:"size:36;name:企业" json:"ent_id"`
	PageID      string         `gorm:"size:36;name:页面" json:"page_id"` //page，ref，app
	Element     string         `gorm:"size:20;name:元素" json:"element"`
	ParentCode  string         `gorm:"size:36;name:上级元素" json:"parent_code"`
	Children    []MDPageWidget `gorm:"-" json:"children"`
	Code        string         `gorm:"size:50" json:"code"`
	Name        string         `gorm:"size:50" json:"name"`
	Entity      string         `gorm:"size:50;name:实体" json:"entity"`
	Field       string         `gorm:"size:50;name:字段" json:"field"`
	Extras      utils.SJson    `gorm:"type:text;name:扩展属性" json:"extras"` //JSON
	Required    utils.SBool    `gorm:"name:是否必填" json:"required"`
	Hidden      utils.SBool    `gorm:"name:是否隐藏" json:"hidden"`
	Editable    utils.SBool    `gorm:"name:可编辑" json:"editable"`
	Placeholder string         `gorm:"size:50;name:占位符" json:"placeholder"`
	Sequence    int            `gorm:"name:顺序号" json:"sequence"`
	Value       utils.SJson    `gorm:"type:text;name:默认值" json:"value"`
	Align       string         `gorm:"size:20;name:对齐" json:"align"`
	Width       string         `gorm:"size:10;name:宽度" json:"width"`
	InputType   string         `gorm:"size:10;name:输入类型" json:"input_type"`
	DataSource  string         `gorm:"size:50;name:数据来源" json:"data_source"`
	DataType    string         `gorm:"size:50;name:数据类型" json:"data_type"`
}

func (s *MDPageWidget) MD() *Mder {
	return &Mder{ID: "md.page.widget", Domain: md_domain, Name: "页面部件"}
}
