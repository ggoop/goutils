package md

import (
	"github.com/ggoop/goutils/context"
	"regexp"
	"strings"
)

//公共查询
type OQL struct {
	Error    error
	errors   []error
	entities map[string]*oqlEntity
	fields   map[string]*oqlField
	froms    []*OQLFrom
	joins    []*OQLJoin
	selects  []*OQLSelect
	orders   []*OQLOrder
	wheres   []*OQLWhere
	groups   []*OQLGroup
	having   []*OQLWhere
	offset   int
	limit    int
	context  *context.Context
	actuator OQLActuator
}

func (s *OQL) GetMainFrom() *OQLFrom {
	if s.froms == nil || len(s.froms) == 0 {
		return nil
	}
	return s.froms[0]
}
func (s *OQL) SetContext(context *context.Context) *OQL {
	s.context = context
	return s
}
func (s *OQL) SetActuator(actuator OQLActuator) *OQL {
	s.actuator = actuator
	return s
}

// 设置 主 from ，示例：
//  tableA
//	tableA as a
//	tableA a
func (s *OQL) From(query interface{}, args ...interface{}) *OQL {
	if v, ok := query.(string); ok {
		seg := &OQLFrom{Query: v, Args: args}
		r := regexp.MustCompile(REGEXP_OQL_FROM)
		matches := r.FindStringSubmatch(v)
		if matches != nil && len(matches) == 4 {
			if matches[2] != "" {
				seg.Query = matches[1]
				seg.Alias = matches[2]
			} else {
				seg.Query = matches[3]
			}
		}
		s.froms = append(s.froms, seg)
	} else if v, ok := query.(OQLFrom); ok {
		s.froms = append(s.froms, &v)
	}
	return s
}
func (s *OQL) Join(joinType OQLJoinType, query string, condition string, args ...interface{}) *OQL {
	seg := &OQLJoin{Type: joinType, Query: query, Condition: condition, Args: args}
	r := regexp.MustCompile(REGEXP_OQL_FROM)
	matches := r.FindStringSubmatch(query)
	if matches != nil && len(matches) == 4 {
		if matches[2] != "" {
			seg.Query = matches[1]
			seg.Alias = matches[2]
		} else {
			seg.Query = matches[3]
		}
	}
	s.joins = append(s.joins, seg)
	return s
}

// 添加字段，示例：
//	单字段：fieldA ，fieldA as A
//	复合字段：sum(fieldA) AS A，fieldA+fieldB as c
//
func (s *OQL) Select(query interface{}, args ...interface{}) *OQL {
	if v, ok := query.(string); ok {
		seg := &OQLSelect{Query: v, Args: args}
		r := regexp.MustCompile(REGEXP_OQL_SELECT)
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
	} else if v, ok := query.(OQLSelect); ok {
		s.selects = append(s.selects, &v)
	}
	return s
}

//排序，示例：
// fieldA desc，fieldA + fieldB
func (s *OQL) Order(query interface{}, args ...interface{}) *OQL {
	if v, ok := query.(string); ok {
		seg := &OQLOrder{Query: v, Args: args}
		r := regexp.MustCompile(REGEXP_OQL_ORDER)
		matches := r.FindStringSubmatch(v)
		if matches != nil && len(matches) == 4 {
			if matches[2] != "" {
				seg.Query = matches[1]
				if strings.ToLower(matches[2]) == "desc" {
					seg.Order = OQL_ORDER_DESC
				} else {
					seg.Order = OQL_ORDER_ASC
				}
			} else {
				seg.Query = matches[3]
			}
		}
		s.orders = append(s.orders, seg)
	} else if v, ok := query.(OQLOrder); ok {
		s.orders = append(s.orders, &v)
	}
	return s
}
func (s *OQL) Group(query interface{}, args ...interface{}) *OQL {
	if v, ok := query.(string); ok {
		seg := &OQLGroup{Query: v, Args: args}
		s.groups = append(s.groups, seg)
	} else if v, ok := query.(OQLGroup); ok {
		s.groups = append(s.groups, &v)
	}
	return s
}
func (s *OQL) Where(query interface{}, args ...interface{}) *OQLWhere {
	var seg *OQLWhere
	if v, ok := query.(string); ok {
		seg = NewOQLWhere(v, args)
	} else if v, ok := query.(OQLWhere); ok {
		seg = &v
	} else {
		seg = NewOQLWhere("", args)
	}
	s.wheres = append(s.wheres, seg)
	return seg
}

func (s *OQL) Having(query interface{}, args ...interface{}) *OQLWhere {
	var seg *OQLWhere
	if v, ok := query.(string); ok {
		seg = NewOQLWhere(v, args)
	} else if v, ok := query.(OQLWhere); ok {
		seg = &v
	} else {
		seg = NewOQLWhere("", args)
	}
	s.having = append(s.having, seg)
	return seg
}
func (s *OQL) Count(value interface{}) *OQL {
	return s.actuator.Count(s, value)
}
func (s *OQL) Pluck(column string, value interface{}) *OQL {
	return s.actuator.Pluck(s, column, value)
}
func (s *OQL) Take(out interface{}) *OQL {
	return s.actuator.Take(s, out)
}
func (s *OQL) Find(out interface{}) *OQL {
	return s.actuator.Find(s, out)
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

	return s.actuator.Find(s, value)
}

//insert into table (aaa,aa,aa) values(aaa,aaa,aaa)
//field 从select 取， value 从 data 取
func (s *OQL) Create(data interface{}) *OQL {
	return s.actuator.Create(s, data)
}

//update table set aa=bb
//field 从select 取， value 从 data 取
func (s *OQL) Update(data interface{}) *OQL {
	return s.actuator.Update(s, data)
}
func (s *OQL) Delete() *OQL {
	return s.actuator.Delete(s)
}
func (s *OQL) AddErr(err error) *OQL {
	s.errors = append(s.errors, err)
	s.Error = err
	return s
}
