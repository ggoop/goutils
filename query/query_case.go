package query

import (
	"fmt"

	"github.com/ggoop/goutils/md"
)

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
func NewCaseExector(item QueryCase) IExector {
	if item.Query == nil {
		return nil
	}
	exector := NewExector(item.Query.Entry)
	if item.Columns != nil && len(item.Columns) > 0 {
		for _, v := range item.Columns {
			if v.ColumnName != "" {
				exector.Select(v.Field + " as " + v.ColumnName)
			} else {
				exector.Select(v.Field)
			}
		}
	}
	if item.Wheres != nil && len(item.Wheres) > 0 {
		for _, v := range item.Wheres {
			iw := queryWhereToIWhere(v)
			if iw != nil {
				if iw.GetLogical() == "or" {
					iw = exector.OrWhere(iw.GetQuery(), iw.GetArgs())
				} else {
					iw = exector.Where(iw.GetQuery(), iw.GetArgs())
				}
				if v.Children != nil && len(v.Children) > 0 {
					for _, item := range v.Children {
						addSubItemToIWhere(iw, item)
					}
				}
			}
		}
	}
	return exector
}
func addSubItemToIWhere(iw IWhere, subValue QueryWhere) {
	newIw := queryWhereToIWhere(subValue)
	if iw.GetLogical() == "or" {
		newIw = iw.OrWhere(newIw.GetQuery(), newIw.GetArgs())
	} else {
		newIw = iw.Where(newIw.GetQuery(), newIw.GetArgs())
	}
	if subValue.Children != nil && len(subValue.Children) > 0 {
		for _, item := range subValue.Children {
			addSubItemToIWhere(newIw, item)
		}
	}
}
func queryWhereToIWhere(value QueryWhere) IWhere {
	item := oqlWhere{Logical: value.Logical}
	if value.Field != "" && value.Operator == "contains" {
		item.Query = fmt.Sprintf("%v like ?", value.Field)
		item.Args = []interface{}{"%" + value.Value + "%"}
	} else if value.Field != "" && (value.Operator == "in" || value.Operator == "not in") {
		item.Query = fmt.Sprintf("%v %s (?)", value.Field, value.Operator)
		item.Args = []interface{}{value.Value}
	} else if value.Field != "" && (value.Operator == "=" || value.Operator == "<>" || value.Operator == ">" || value.Operator == ">=" || value.Operator == "<" || value.Operator == "<=") {
		item.Query = fmt.Sprintf("%v %s ?", value.Field, value.Operator)
		item.Args = []interface{}{value.Value}
	}
	return &item
}
