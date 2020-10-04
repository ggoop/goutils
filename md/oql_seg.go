package md

import (
	"regexp"
	"strings"
)

//公共查询
type OQL struct {
	Error   error
	from    OQLFrom
	joins   []OQLJoin
	selects []OQLField
	orders  []OQLOrder
	wheres  []OQLWhere
	groups  []OQLField
	having  []OQLWhere
	offset  int
	limit   int
}

// 设置 主 from ，示例：
//  tableA
//	tableA as a
//	tableA a
func (s *OQL) From(query interface{}, args ...interface{}) *OQL {
	if v, ok := query.(string); ok {
		seg := &oqlFrom{Origin: query, Args: args}
		r := regexp.MustCompile(regexp_oql_from)
		matches := r.FindStringSubmatch(v)
		if matches != nil && len(matches) == 4 {
			if matches[2] != "" {
				seg.Query = matches[1]
				seg.Alias = matches[2]
			} else {
				seg.Query = matches[3]
			}
		}
		s.from = seg
	} else if v, ok := query.(OQLFrom); ok {
		s.from = v
	}
	return s
}

// 添加字段，示例：
//	单字段：fieldA ，fieldA as A
//	复合字段：sum(fieldA) AS A，fieldA+fieldB as c
//
func (s *OQL) Select(query interface{}, args ...interface{}) *OQL {
	if v, ok := query.(string); ok {
		seg := &oqlField{Origin: query, Args: args}
		r := regexp.MustCompile(regexp_oql_select)
		matches := r.FindStringSubmatch(v)
		if matches != nil && len(matches) == 4 {
			if matches[2] != "" {
				seg.Query = matches[1]
				seg.Alias = matches[2]
			} else {
				seg.Query = matches[3]
			}
		}
		s.selects = append(s.selects, seg)
	} else if v, ok := query.(OQLField); ok {
		s.selects = append(s.selects, v)
	}
	return s
}
func (s *OQL) Order(query interface{}, args ...interface{}) *OQL {
	if v, ok := query.(string); ok {
		seg := &oqlOrder{Origin: query, Args: args}
		r := regexp.MustCompile(regexp_oql_order)
		matches := r.FindStringSubmatch(v)
		if matches != nil && len(matches) == 4 {
			if matches[2] != "" {
				seg.Query = matches[1]
				if strings.ToLower(matches[2]) == "desc" {
					seg.Sequence = -1
				} else {
					seg.Sequence = 1
				}
			} else {
				seg.Query = matches[3]
			}
		}
		s.orders = append(s.orders, seg)
	} else if v, ok := query.(OQLOrder); ok {
		s.orders = append(s.orders, v)
	}
	return s
}
func (s *OQL) Group(query interface{}, args ...interface{}) *OQL {
	if v, ok := query.(string); ok {
		seg := &oqlField{Origin: query, Args: args, Query: v}
		s.groups = append(s.groups, seg)
	} else if v, ok := query.(OQLField); ok {
		s.groups = append(s.groups, v)
	}
	return s
}
func (s *OQL) Where(query interface{}, args ...interface{}) OQLWhere {
	var seg OQLWhere
	if v, ok := query.(string); ok {
		seg = NewOQLWhere(v, args)
	} else if v, ok := query.(OQLWhere); ok {
		seg = v
	} else {
		seg = NewOQLWhere("", args)
	}
	s.wheres = append(s.wheres, seg)
	return seg
}

func (s *OQL) Having(query interface{}, args ...interface{}) OQLWhere {
	var seg OQLWhere
	if v, ok := query.(string); ok {
		seg = NewOQLWhere(v, args)
	} else if v, ok := query.(OQLWhere); ok {
		seg = v
	} else {
		seg = NewOQLWhere("", args)
	}
	s.having = append(s.having, seg)
	return seg
}
func (s *OQL) Count(value interface{}) *OQL {
	return GetOQLActuator().Count(s, value)
}
func (s *OQL) Pluck(column string, value interface{}) *OQL {
	return GetOQLActuator().Pluck(s, column, value)
}
func (s *OQL) Take(out interface{}) *OQL {
	return GetOQLActuator().Take(s, out)
}
func (s *OQL) Find(out interface{}) *OQL {
	return GetOQLActuator().Find(s, out)
}
func (s *OQL) Paginate(value interface{}, page int, pageSize int) *OQL {
	if pageSize > 0 && page <= 0 {
		page = 1
	} else if pageSize <= 0 {
		pageSize = 0
		page = 0
	}
	s.limit = pageSize
	s.offset = (page - 1) * pageSize

	return GetOQLActuator().Find(s, value)
}

//insert into table (aaa,aa,aa) values(aaa,aaa,aaa)
//field 从select 取， value 从 data 取
func (s *OQL) Create(data interface{}) *OQL {
	return GetOQLActuator().Create(s, data)
}

//update table set aa=bb
//field 从select 取， value 从 data 取
func (s *OQL) Update(data interface{}) *OQL {
	return GetOQLActuator().Update(s, data)
}

func (s *OQL) Delete() *OQL {
	return GetOQLActuator().Delete(s)
}
