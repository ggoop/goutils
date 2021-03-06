package md

import "github.com/ggoop/goutils/utils"

type UIWidget struct {
	ID        string      `gorm:"primary_key;size:36" json:"id"`
	CreatedAt utils.Time  `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt utils.Time  `gorm:"name:更新时间" json:"updated_at"`
	EntID     string      `gorm:"size:36;name:企业" json:"ent_id"`
	Code      string      `gorm:"size:36;name:编码" json:"code"`
	Name      string      `gorm:"size:50;name:名称" json:"name"`
	Type      string      `gorm:"size:20;name:类型" json:"type"`
	FilterID  string      `gorm:"size:36;name:过滤器" json:"filter_id"`
	AuthID    string      `gorm:"size:36;name:权限ID" json:"auth_id"`
	IDField   string      `gorm:"size:36;name:主键字段" json:"id_field"`
	Extras    utils.SJson `gorm:"type:text" json:"extras"` //JSON
	System    utils.SBool `gorm:"name:系统的" json:"system"`
}

func (s *UIWidget) MD() *Mder {
	return &Mder{ID: "md.ui.widget", Domain: md_domain, Name: "组件"}
}

type UIEntity struct {
	ID         string      `gorm:"primary_key;size:36" json:"id"`
	CreatedAt  utils.Time  `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt  utils.Time  `gorm:"name:更新时间" json:"updated_at"`
	EntID      string      `gorm:"size:36;name:企业" json:"ent_id"`
	WidgetID   string      `gorm:"size:36;name:组件ID" json:"widget_id"`
	Code       string      `gorm:"size:36;name:编码" json:"code"`
	Name       string      `gorm:"size:50;name:名称" json:"name"`
	DsType     string      `gorm:"size:36;name:数据源类型" json:"ds_type"`
	DsEntry    string      `gorm:"size:36;name:数据源实体" json:"ds_entry"`
	IsMain     utils.SBool `gorm:"name:是否主实体" json:"is_main"`
	PrimaryKey string      `gorm:"size:36;name:主键" json:"primary_key"`
	ForeignKey string      `gorm:"size:36;name:外键" json:"foreign_key"`
	ParentID   string      `gorm:"size:36;name:上级" json:"parent_id"`
	Kind       string      `gorm:"size:20;name:关联关系" json:"kind"`
	Limit      string      `gorm:"type:text;name:条件限制" json:"limit"`
	System     utils.SBool `gorm:"name:系统的" json:"system"`
}

func (s *UIEntity) MD() *Mder {
	return &Mder{ID: "md.ui.entity", Domain: md_domain, Name: "组件数据源"}
}

type UILayout struct {
	ID        string      `gorm:"primary_key;size:36" json:"id"`
	CreatedAt utils.Time  `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt utils.Time  `gorm:"name:更新时间" json:"updated_at"`
	EntID     string      `gorm:"size:36;name:企业" json:"ent_id"`
	WidgetID  string      `gorm:"size:36;name:组件ID" json:"widget_id"`
	ParentID  string      `gorm:"size:36;name:上级" json:"parent_id"`
	Children  []UILayout  `gorm:"-" json:"children"`
	Code      string      `gorm:"size:36;name:编码" json:"code"`
	Name      string      `gorm:"size:50;name:名称" json:"name"`
	Type      string      `gorm:"size:20;name:布局类型" json:"type"`
	Sequence  int         `gorm:"size:3;name:顺序" json:"sequence"`
	Align     string      `gorm:"size:20;name:对齐方式" json:"align"`
	Cols      int         `gorm:"size:3;name:列数" json:"cols"`
	Style     utils.SJson `gorm:"type:text;name:样式" json:"style"`
	System    utils.SBool `gorm:"name:系统的" json:"system"`
	Items     []UIItem    `gorm:"-" json:"items"`
}

func (s *UILayout) MD() *Mder {
	return &Mder{ID: "md.ui.layout", Domain: md_domain, Name: "组件布局"}
}

type UIItem struct {
	ID          string      `gorm:"primary_key;size:36" json:"id"`
	CreatedAt   utils.Time  `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt   utils.Time  `gorm:"name:更新时间" json:"updated_at"`
	EntID       string      `gorm:"size:36;name:企业" json:"ent_id"`
	WidgetID    string      `gorm:"size:36;name:组件ID" json:"widget_id"`
	LayoutID    string      `gorm:"size:36;name:布局ID" json:"layout_id"`
	ParentID    string      `gorm:"size:36;name:上级" json:"parent_id"`
	Children    []UIItem    `gorm:"-" json:"children"`
	Code        string      `gorm:"size:36;name:编码" json:"code"`
	Name        string      `gorm:"size:50;name:名称" json:"name"`
	Type        string      `gorm:"size:20;name:类型" json:"type"`
	Caption     string      `gorm:"size:50;name:标题" json:"caption"`
	DsType      string      `gorm:"size:36;name:数据源类型" json:"ds_type"`
	DsEntry     string      `gorm:"size:36;name:数据源实体" json:"ds_entry"`
	DsField     string      `gorm:"size:36;name:数据源字段" json:"ds_field"`
	RefType     string      `gorm:"size:36;name:参照类型" json:"ref_type"`
	RefCode     string      `gorm:"size:36;name:参照编码" json:"ref_code"`
	RefReturn   string      `gorm:"size:36;name:参照返回" json:"ref_return"`
	RefFilter   string      `gorm:"type:text;name:参照查询条件" json:"ref_filter"`
	Sequence    int         `gorm:"size:3;name:顺序" json:"sequence"`
	Value1      utils.SJson `gorm:"type:text;name:值1" json:"value1"`
	Value2      utils.SJson `gorm:"type:text;name:值2" json:"value2"`
	Precision   int         `gorm:"size:3;name:精度" json:"precision"`
	Format      string      `gorm:"size:36;name:格式化" json:"format"`
	Placeholder string      `gorm:"size:50;name:占位符" json:"placeholder"`
	Length      int         `gorm:"size:3;name:长度" json:"length"`
	Hidden      utils.SBool `gorm:"name:隐藏" json:"hidden"`
	Multiple    utils.SBool `gorm:"name:多选" json:"multiple"`
	Nullable    utils.SBool `gorm:"name:可空" json:"nullable"`
	Editable    utils.SBool `gorm:"name:可编辑" json:"editable"`
	Fixed       utils.SBool `gorm:"name:固定的" json:"fixed"`
	Width       string      `gorm:"name:宽度" json:"width"`
	Align       string      `gorm:"size:20;name:对齐方式" json:"align"`
	Style       utils.SJson `gorm:"type:text;name:样式" json:"style"`
	Extras      utils.SJson `gorm:"type:text;name:扩展属性" json:"extras"` //JSON
	System      utils.SBool `gorm:"name:系统的" json:"system"`
}

func (s *UIItem) MD() *Mder {
	return &Mder{ID: "md.ui.item", Domain: md_domain, Name: "组件元素"}
}
