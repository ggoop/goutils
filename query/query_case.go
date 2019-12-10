package query

import (
	"fmt"
	"strings"

	"github.com/ggoop/goutils/context"
	"github.com/ggoop/goutils/di"
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/md"
	"github.com/ggoop/goutils/repositories"
)

type QueryCase struct {
	md.Model
	EntID     string           `gorm:"size:50" json:"ent_id"`
	UserID    string           `gorm:"size:50" json:"user_id"`
	QueryID   string           `gorm:"size:50;name:查询" json:"query_id"`
	Query     *Query           `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;name:查询" json:"query"`
	Name      string           `gorm:"name:名称" json:"name"`
	ScopeID   string           `gorm:"size:50;name:范围ID" json:"scope_id"`
	ScopeType string           `gorm:"size:50;name:范围类型" json:"scope_type"`
	Memo      string           `gorm:"name:备注" json:"memo"`
	Page      int              `gorm:"name:页码" json:"page"`
	PageSize  int              `gorm:"name:每页显示记录数" json:"page_size"`
	Columns   []QueryColumn    `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:OwnerID;name:栏目集合" json:"columns"`
	Orders    []QueryOrder     `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:OwnerID;name:排序集合" json:"orders"`
	Wheres    []QueryWhere     `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:OwnerID;name:条件集合" json:"wheres"`
	Context   *context.Context `gorm:"type:text;name:上下文" json:"context"` //上下文参数
	Export    bool             `gorm:"name:是否导出" json:"export"`
}

func (s *QueryCase) MD() *md.Mder {
	return &md.Mder{ID: "query.case", Name: "查询方案"}
}
func (s *QueryCase) Format() *QueryCase {
	if s.Query == nil && s.QueryID != "" {
		if err := di.Global.Invoke(func(db *repositories.MysqlRepo) {
			q := Query{}
			if err := db.Preload("Columns").Preload("Orders").Preload("Wheres").Take(&q, "id=? or code=?", s.QueryID, s.QueryID).Error; err != nil {
				glog.Errorf("query error:%s", err)
			}
			if q.ID != "" {
				s.Query = &q
			}
		}); err != nil {
			glog.Errorf("di Provide error:%s", err)
		}
	}
	if s.Columns == nil && len(s.Columns) == 0 {
		s.Columns = make([]QueryColumn, 0)
		for _, v := range s.Query.Columns {
			s.Columns = append(s.Columns, v)
		}
	}
	for _, v := range s.Columns {
		if v.Name == "" {
			v.Name = strings.Replace(strings.ToLower(v.Field), ".", "_", -1)
		}
	}
	if s.Orders == nil && len(s.Orders) == 0 {
		s.Orders = make([]QueryOrder, 0)
		for _, v := range s.Query.Orders {
			s.Orders = append(s.Orders, v)
		}
	}
	if s.Page <= 0 {
		s.Page = 1
	}
	if s.PageSize <= 0 {
		s.PageSize = s.Query.PageSize
	}
	if s.PageSize <= 0 {
		s.PageSize = 30
	}
	return s
}
func (s *QueryCase) GetExector() IExector {
	s.Format()
	if s.Query == nil {
		return nil
	}
	exector := NewExector(s.Query.Entry)
	if s.Page > 0 && s.PageSize > 0 {
		exector.Page(s.Page, s.PageSize)
	}
	exector.SetContext(s.Context)
	for _, v := range s.Columns {
		if v.Name != "" {
			exector.Select(v.Field + " as " + v.Name)
		} else {
			exector.Select(v.Field)
		}
	}
	for _, v := range s.Wheres {
		iw := s.queryWhereToIWhere(v)
		if iw != nil {
			if iw.GetLogical() == "or" {
				iw = exector.OrWhere(iw.GetQuery(), iw.GetArgs())
			} else {
				iw = exector.Where(iw.GetQuery(), iw.GetArgs())
			}
			if v.Children != nil && len(v.Children) > 0 {
				for _, item := range v.Children {
					s.addSubItemToIWhere(iw, item)
				}
			}
		}
	}
	if len(s.Orders) > 0 {
		for _, v := range s.Orders {
			if v.Order != "" {
				exector.Order(v.Field + " " + v.Order)
			} else {
				exector.Order(v.Field)
			}
		}
	} else {
		exector.Order("ID")
	}

	return exector
}
func (s *QueryCase) addSubItemToIWhere(iw IQWhere, subValue QueryWhere) {
	newIw := s.queryWhereToIWhere(subValue)
	if iw.GetLogical() == "or" {
		newIw = iw.OrWhere(newIw.GetQuery(), newIw.GetArgs())
	} else {
		newIw = iw.Where(newIw.GetQuery(), newIw.GetArgs())
	}
	if subValue.Children != nil && len(subValue.Children) > 0 {
		for _, item := range subValue.Children {
			s.addSubItemToIWhere(newIw, item)
		}
	}
}
func (s *QueryCase) queryWhereToIWhere(value QueryWhere) IQWhere {
	item := qWhere{Logical: value.Logical}
	if value.Field != "" && value.Value != "" && value.Operator == "contains" {
		item.Query = fmt.Sprintf("%v like ?", value.Field)
		item.Args = []interface{}{"%" + value.Value + "%"}
	} else if value.Field != "" && value.Value != "" && value.Operator == "like" {
		item.Query = fmt.Sprintf("%v like ?", value.Field)
		item.Args = []interface{}{"%" + value.Value + "%"}
	} else if value.Field != "" && value.Value != "" && value.Operator == "not like" {
		item.Query = fmt.Sprintf("%v not like ?", value.Field)
		item.Args = []interface{}{"%" + value.Value + "%"}
	} else if value.Field != "" && value.Value != "" && (value.Operator == "in" || value.Operator == "not in") {
		item.Query = fmt.Sprintf("%v %s (?)", value.Field, value.Operator)
		item.Args = []interface{}{value.Value}
	} else if value.Field != "" && value.Value != "" && (value.Operator == "=" || value.Operator == "<>" || value.Operator == ">" || value.Operator == ">=" || value.Operator == "<" || value.Operator == "<=") {
		item.Query = fmt.Sprintf("%v %s ?", value.Field, value.Operator)
		item.Args = []interface{}{value.Value}
	} else if value.Field != "" && value.Value != "" && s.Context != nil && (value.Operator == "=p" || value.Operator == "<>p" || value.Operator == ">p" || value.Operator == ">=p" || value.Operator == "<p" || value.Operator == "<=p") {
		item.Query = fmt.Sprintf("%v %s ?", value.Field, strings.Replace(value.Operator, "p", "", -1))
		item.Args = []interface{}{s.Context.GetValue(strings.Replace(value.Value, "@", "", -1))}
	}
	return &item
}
