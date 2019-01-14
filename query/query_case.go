package query

import "github.com/ggoop/goutils/md"

type QueryCase struct {
	md.ModelUnscoped
	EntID     string        `gorm:"size:100" json:"ent_id"`
	UserID    string        `gorm:"size:100" json:"user_id"`
	QueryID   string        `gorm:"name:查询ID" json:"case_id"`
	Query     *Query        `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;name:查询" json:"query"`
	Name      string        `gorm:"name:名称" json:"name"`
	Content   string        `gorm:"name:内容" json:"content"`
	ScopeType string        `gorm:"name:范围类型" json:"scope_type"`
	ScopeID   string        `gorm:"name:范围ID" json:"scope_id"`
	Memo      string        `gorm:"name:备注" json:"memo"`
	Columns   []QueryColumn `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:CaseID;name:栏目集合" json:"columns"`
	Orders    []QueryOrder  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:CaseID;name:排序集合" json:"orders"`
	Wheres    []QueryWhere  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:CaseID;name:条件集合" json:"wheres"`
}

func (s *QueryCase) MD() *md.Mder {
	return &md.Mder{ID: "01e916da3fbb092f44be8cec4b7174de", Name: "查询方案"}
}
