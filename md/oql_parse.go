package md

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/utils"
)

func (s *OQL) parseOQL(value interface{}) *OQL {
	for _, v := range s.froms {
		s.parseFromField(v)
	}
	for _, v := range s.joins {
		s.parseJoinField(v)
	}
	for _, v := range s.selects {
		s.parseSelectField(v)
	}
	for _, v := range s.wheres {
		s.parseWhereField(v)
	}
	for _, v := range s.having {
		s.parseWhereField(v)
	}
	for _, v := range s.orders {
		s.parseOrderField(v)
	}
	for _, v := range s.groups {
		s.parseGroupField(v)
	}
	return s
}

func (s *OQL) parseFromField(value *OQLFrom) error {
	//主表，使用别名作路径
	form := s.parseEntity(value.Query, value.Alias)
	if s.Error != nil {
		return s.Error
	}
	form.IsMain = true
	value.expr = fmt.Sprintf("%s  %s", form.Entity.TableName, form.Alias)
	return nil
}
func (s *OQL) parseJoinField(value *OQLJoin) error {
	//主表，使用别名作路径
	form := s.parseEntity(value.Query, value.Alias)
	if s.Error != nil {
		return s.Error
	}
	joins := make([]string, 0)
	switch value.Type {
	case OQL_LEFT_JOIN:
		joins = append(joins, "left join")
	case OQL_RIGHT_JOIN:
		joins = append(joins, "left join")
	case OQL_FULL_JOIN:
		joins = append(joins, "join")
	}
	joins = append(joins, fmt.Sprintf("%s  %s", form.Entity.TableName, form.Alias))
	if value.Condition != "" {
		condition := s.parseFieldExpr(value.Condition)
		if condition != "" {
			joins = append(joins, "on ", condition)
		}
	}
	value.expr = strings.Join(joins, " ")
	return nil
}
func (s *OQL) parseWhereField(value *OQLWhere) {
	value.expr = s.parseFieldExpr(value.Query)
	if len(value.Children) > 0 {
		for _, v := range value.Children {
			s.parseWhereField(v)
		}
	}
}
func (s *OQL) parseSelectField(value *OQLSelect) {
	value.expr = s.parseFieldExpr(value.Query)
}
func (s *OQL) parseGroupField(value *OQLGroup) {
	value.expr = s.parseFieldExpr(value.Query)
}
func (s *OQL) parseOrderField(value *OQLOrder) {
	value.expr = s.parseFieldExpr(value.Query)
}

// 解析字段表达式，如
//	a.fieldA+fieldB+sum(b.fieldA)   =>a.fieldA ,fieldB, b.fieldA
//	$$a.fieldA + sum( c.fieldA )	=>$$a.fieldA, c.fieldA
// 函数与左括号之间不能有空格
// 多级字段.号不能有空格
func (s *OQL) parseFieldExpr(expr string) string {
	if expr == "" {
		return expr
	}
	r, _ := regexp.Compile(`([\$]?[A-Za-z._]+[0-9A-Za-z|\(])`)
	matches := r.FindAllStringSubmatch(expr, -1)
	for _, match := range matches {
		str := match[1]
		//带有括号的是函数，不需要解析
		if strings.Index(str, utils.PARENTHESIS_LEFT) > 0 {
			field, _ := s.parseField(str)
			if field != nil {
				expr = strings.ReplaceAll(expr, str, fmt.Sprintf("%s.%s", field.Entity.Alias, field.Field.DbName))
			}
		}
	}
	return expr
}

// 解析实体
func (s *OQL) formatEntity(entity *MDEntity) *oqlEntity {
	e := oqlEntity{Entity: entity}
	return &e
}
func (s *OQL) formatField(entity *oqlEntity, field *MDField) *oqlField {
	e := oqlField{Entity: entity, Field: field}
	return &e
}
func (s *OQL) parseEntity(id, path string) *oqlEntity {
	path = strings.ToLower(strings.TrimSpace(path))
	if v, ok := s.entities[path]; ok {
		return v
	}
	entity := GetEntity(id)
	if entity == nil {
		err := glog.Errorf("找不到实体 %v", id)
		s.AddErr(err)
		return nil
	}
	v := s.formatEntity(entity)
	v.Sequence = len(s.entities) + 1
	v.Alias = fmt.Sprintf("a%v", v.Sequence)
	v.Path = path
	s.entities[path] = v
	return v
}

// 解析字段
func (s *OQL) parseField(fieldPath string) (*oqlField, error) {
	fieldPath = strings.ToLower(strings.TrimSpace(fieldPath))
	if v, ok := s.fields[fieldPath]; ok {
		return v, nil
	}
	start := 0
	parts := strings.Split(fieldPath, ".")
	var mainFrom *OQLFrom
	if len(parts) > 1 {
		//如果主表有别名，则第一个字段为表
		for i, v := range s.froms {
			if v.Alias != "" && strings.ToLower(v.Alias) == parts[0] {
				mainFrom = s.froms[i]
				start = 1
				break
			}
		}
	}
	if mainFrom == nil {
		//如果没有找到主表，则说明字段没有 表作为导引
		for i, v := range s.froms {
			if v.Alias == "" {
				mainFrom = s.froms[i]
				break
			}
		}
	}
	if mainFrom == nil {
		mainFrom = s.froms[0]
	}
	//主实体
	entity := s.entities[strings.ToLower(mainFrom.Alias)]

	path := ""
	for i, part := range parts {
		if i > 0 {
			path += "."
		}
		path += part
		if i < start {
			continue
		}
		mdField := entity.Entity.GetField(part)
		if mdField == nil {
			return nil, nil
		}
		field := s.formatField(entity, mdField)
		field.Path = path
		s.fields[path] = field
		if i < len(parts)-1 {
			entity = s.parseEntity(mdField.TypeID, path)
			if s.Error != nil {
				return nil, nil
			}
		} else {
			return field, nil
		}
	}
	return nil, nil
}
