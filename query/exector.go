package query

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/md"
	"github.com/ggoop/goutils/repositories"
	"github.com/jinzhu/gorm"
)

type IExector interface {
	PrepareQuery(mysql *repositories.MysqlRepo) (*gorm.DB, error)
	Select(query string, args ...interface{}) IExector
	Where(query string, args ...interface{}) IQWhere
	OrWhere(query string, args ...interface{}) IQWhere
	Joins(query string, args ...interface{}) IExector
	Group(query string, args ...interface{}) IExector
}
type oqlEntity struct {
	Alia     string
	IsMain   bool
	Entity   *md.MDEntity
	Path     string
	Sequence int
}
type oqlField struct {
	Entity *oqlEntity
	Field  *md.MDField
	Path   string
}
type oqlFrom struct {
	Query string
	Alia  string
	Expr  string
}
type oqlJoin struct {
	Query string
	Expr  string
	Args  []interface{}
}
type oqlSelect struct {
	Query string
	Expr  string
	Args  []interface{}
}
type oqlOrder struct {
	Query    string
	Sequence int
	Expr     string
	Args     []interface{}
}
type oqlGroup struct {
	Query string
	Expr  string
	Args  []interface{}
}

type exector struct {
	entities map[string]*oqlEntity
	fields   map[string]*oqlField
	froms    []*oqlFrom
	selects  []*oqlSelect
	joins    []*oqlJoin
	wheres   []*qWhere
	orders   []*oqlOrder
	groups   []*oqlGroup
}

func NewExector(query string) IExector {
	exec := &exector{
		entities: make(map[string]*oqlEntity),
		fields:   make(map[string]*oqlField),
		froms:    make([]*oqlFrom, 0),
		selects:  make([]*oqlSelect, 0),
		joins:    make([]*oqlJoin, 0),
		wheres:   make([]*qWhere, 0),
		orders:   make([]*oqlOrder, 0),
		groups:   make([]*oqlGroup, 0),
	}
	exec.From(query)
	return exec
}

func (m *exector) FormatEntity(entity *md.MDEntity) *oqlEntity {
	e := oqlEntity{Entity: entity}
	return &e
}
func (m *exector) FormatField(entity *oqlEntity, field *md.MDField) *oqlField {
	e := oqlField{Entity: entity, Field: field}
	return &e
}
func (m *exector) PrepareQuery(mysql *repositories.MysqlRepo) (*gorm.DB, error) {
	//parse
	for _, v := range m.froms {
		m.parseFromField(v)
	}
	for _, v := range m.selects {
		m.parseSelectField(v)
	}
	for _, v := range m.wheres {
		m.parseWhereField(v)
	}
	for _, v := range m.orders {
		m.parseOrderField(v)
	}
	for _, v := range m.groups {
		m.parseGroupField(v)
	}
	//build

	queryDB := m.buildFroms(mysql)
	queryDB = m.buildSelects(queryDB)
	queryDB = m.buildJoins(queryDB)

	if whereExpr, whereArgs, tag := m.buildWheres(m.wheres); tag > 0 {
		queryDB = queryDB.Where(whereExpr, whereArgs...)
	}
	queryDB = m.buildGroups(queryDB)
	queryDB = m.buildOrders(queryDB)
	return queryDB, nil
}

///=============== build
func (m *exector) buildFroms(mysql *repositories.MysqlRepo) *gorm.DB {
	parts := []string{}
	for _, v := range m.froms {
		parts = append(parts, v.Expr)
	}
	return mysql.Table(strings.Join(parts, ","))
}
func (m *exector) buildSelects(queryDB *gorm.DB) *gorm.DB {
	selects := make([]string, 0)
	for _, v := range m.selects {
		if v.Expr != "" {
			selects = append(selects, v.Expr)
		}
	}
	queryDB = queryDB.Select(selects)
	return queryDB
}
func (m *exector) buildWheres(wheres []*qWhere) (string, []interface{}, int) {
	if wheres == nil || len(wheres) == 0 {
		return "", nil, 0
	}
	tag := 0
	exprs := []string{}
	args := []interface{}{}
	if wheres != nil && len(wheres) > 0 {
		for _, v := range wheres {
			subExpr, subArgs, subTag := m.buildWheres(v.Children)
			tag += subTag
			if v.Expr != "" { //当前节点加上子节点
				if len(exprs) > 0 {
					exprs = append(exprs, " ", v.Logical, " ")
				}
				if subTag > 0 {
					// (a=b and (a=1 or a=2))
					exprs = append(exprs, "((", v.Expr, ") ", v.Logical, " ", subExpr, ")")
					args = append(args, v.Args...)
					args = append(args, subArgs...)
				} else {
					exprs = append(exprs, "(", v.Expr, ")")
					args = append(args, v.Args...)
				}
				tag += 1
			} else if subTag > 0 { //仅仅有子节点
				if len(exprs) > 0 {
					exprs = append(exprs, " ", v.Logical, " ")
				}
				if subTag > 1 {
					exprs = append(exprs, "(", subExpr, ")")
				} else {
					exprs = append(exprs, subExpr)
				}
				args = append(args, subArgs...)
			}
		}
	}
	return strings.Join(exprs, ""), args, tag
}
func (m *exector) buildGroups(queryDB *gorm.DB) *gorm.DB {
	selects := make([]string, 0)
	for _, v := range m.groups {
		if v.Expr != "" {
			selects = append(selects, v.Expr)
		}
	}
	queryDB = queryDB.Group(strings.Join(selects, ","))
	return queryDB
}
func (m *exector) buildOrders(queryDB *gorm.DB) *gorm.DB {
	selects := make([]string, 0)
	for _, v := range m.orders {
		if v.Expr != "" {
			selects = append(selects, v.Expr)
		}
	}
	queryDB = queryDB.Order(strings.Join(selects, ","))
	return queryDB
}
func (m *exector) buildJoins(queryDB *gorm.DB) *gorm.DB {
	tables := make([]*oqlEntity, 0)
	for _, v := range m.entities {
		tables = append(tables, v)
	}
	sort.Slice(tables, func(i, j int) bool {
		return tables[i].Sequence < tables[j].Sequence
	})
	for _, t := range tables {
		if t.Path == "" || t.IsMain {
			continue
		}
		relationship := m.parseField(t.Path)
		if relationship == nil {
			glog.Errorf("找不到关联字段")
			continue
		}
		if relationship.Field.Kind == "belongs_to" || relationship.Field.Kind == "has_one" {
			fkey := relationship.Entity.Entity.GetField(relationship.Field.ForeignKey)
			lkey := t.Entity.GetField(relationship.Field.AssociationKey)
			queryDB = queryDB.Joins(fmt.Sprintf("left join %v as %v on %v.%v=%v.%v", t.Entity.TableName, t.Alia, t.Alia, lkey.DBName, relationship.Entity.Alia, fkey.DBName))
		} else if relationship.Field.Kind == "has_many" {
			fkey := relationship.Entity.Entity.GetField(relationship.Field.ForeignKey)
			lkey := t.Entity.GetField(relationship.Field.AssociationKey)
			queryDB = queryDB.Joins(fmt.Sprintf("left join %v as %v on %v.%v=%v.%v", t.Entity.TableName, t.Alia, t.Alia, fkey.DBName, relationship.Entity.Alia, lkey.DBName))
		}
	}
	return queryDB
}

///=============== parse
func (m *exector) parseWhereField(value *qWhere) {
	if value.Query != "" {
		parts := strings.Split(strings.TrimSpace(value.Query), " ")
		field := m.parseField(parts[0])
		if field != nil {
			parts[0] = fmt.Sprintf("%s.%s", field.Entity.Alia, field.Field.DBName)
		}
		value.Expr = strings.Join(parts, " ")
	}
	if value.Children != nil && len(value.Children) > 0 {
		for _, v := range value.Children {
			m.parseWhereField(v)
		}
	}
}
func (m *exector) parseFromField(value *oqlFrom) {
	items := strings.Split(strings.TrimSpace(value.Query), ",")
	strs := []string{}
	for _, item := range items {
		parts := strings.Split(strings.TrimSpace(item), " ")
		if len(parts) == 1 {
			parts = append(parts, "")
		}
		form := m.parseEntity(parts[0], parts[len(parts)-1])
		form.IsMain = true
		strs = append(strs, fmt.Sprintf("%s as %s", form.Entity.TableName, form.Alia))
	}
	value.Expr = strings.Join(strs, ",")
}
func (m *exector) parseSelectField(value *oqlSelect) {
	items := strings.Split(strings.TrimSpace(value.Query), ",")
	strs := []string{}
	for _, item := range items {
		parts := strings.Split(strings.TrimSpace(item), " ")
		if len(parts) < 2 {
			parts = append(parts, "as", strings.ToLower(strings.Replace(parts[0], ".", "_", -1)))
		}
		field := m.parseField(parts[0])
		if field != nil {
			parts[0] = fmt.Sprintf("%s.%s", field.Entity.Alia, field.Field.DBName)
		}
		strs = append(strs, strings.Join(parts, " "))
	}
	value.Expr = strings.Join(strs, ",")
}

func (m *exector) parseGroupField(value *oqlGroup) {
	items := strings.Split(strings.TrimSpace(value.Query), ",")
	strs := []string{}
	for _, item := range items {
		field := m.parseField(strings.TrimSpace(item))
		if field != nil {
			strs = append(strs, fmt.Sprintf("%s.%s", field.Entity.Alia, field.Field.DBName))
		} else {
			strs = append(strs, item)
		}
	}
	value.Expr = strings.Join(strs, ",")
}
func (m *exector) parseOrderField(value *oqlOrder) {
	items := strings.Split(strings.TrimSpace(value.Query), ",")
	strs := []string{}
	for _, item := range items {
		parts := strings.Split(strings.TrimSpace(item), " ")
		field := m.parseField(parts[0])
		if field != nil {
			parts[0] = fmt.Sprintf("%s.%s", field.Entity.Alia, field.Field.DBName)
		}
		strs = append(strs, strings.Join(parts, " "))
	}
	value.Expr = strings.Join(strs, ",")
}

// 解析实体
func (m *exector) parseEntity(id, path string) *oqlEntity {
	path = strings.ToLower(strings.TrimSpace(path))
	if v, ok := m.entities[path]; ok {
		return v
	}
	entity := md.GetEntity(id)
	if entity == nil {
		glog.Errorf("找不到实体 %v", id)
		return nil
	}
	v := m.FormatEntity(entity)
	v.Sequence = len(m.entities) + 1
	v.Alia = fmt.Sprintf("a%v", v.Sequence)
	v.Path = path
	m.entities[path] = v
	return v
}

// 解析字段
func (m *exector) parseField(fieldPath string) *oqlField {
	fieldPath = strings.ToLower(strings.TrimSpace(fieldPath))
	if v, ok := m.fields[fieldPath]; ok {
		return v
	}
	start := 0
	parts := strings.Split(fieldPath, ".")
	var mainFrom *oqlFrom
	if len(parts) > 1 {
		for i, v := range m.froms {
			if v.Alia != "" && strings.ToLower(v.Alia) == parts[0] {
				mainFrom = m.froms[i]
				start = 1
				break
			}
		}
	}
	if mainFrom == nil {
		for i, v := range m.froms {
			if v.Alia == "" {
				mainFrom = m.froms[i]
				break
			}
		}
	}
	if mainFrom == nil {
		mainFrom = m.froms[0]
	}
	entity := m.entities[strings.ToLower(mainFrom.Alia)]

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
			return nil
		}
		field := m.FormatField(entity, mdField)
		field.Path = path
		m.fields[path] = field
		if i < len(parts)-1 && mdField.TypeType == md.TYPE_ENTITY {
			entity = m.parseEntity(mdField.TypeID, path)
		} else {
			return field
		}
	}
	return nil
}

// fieldA
// fieldA,fieldB,fieldC
// fieldA as a,fieldB
func (m *exector) Select(query string, args ...interface{}) IExector {
	item := &oqlSelect{Query: query, Args: args}
	m.selects = append(m.selects, item)
	return m
}

// fieldA =?
// fieldB =''
// FieldC is null
func (m *exector) Where(query string, args ...interface{}) IQWhere {
	item := &qWhere{Query: query, Args: args, Logical: "and"}
	m.wheres = append(m.wheres, item)
	return item
}
func (m *exector) OrWhere(query string, args ...interface{}) IQWhere {
	item := &qWhere{Query: query, Args: args, Logical: "or"}
	m.wheres = append(m.wheres, item)
	return item
}

// left join tableA on a=b
// left join tableA as c on c.a=b.b
// inner join tableA
func (m *exector) Joins(query string, args ...interface{}) IExector {
	item := &oqlJoin{Query: query, Args: args}
	m.joins = append(m.joins, item)
	return m
}

//fieldA
//fieldA,fieldB
func (m *exector) Group(query string, args ...interface{}) IExector {
	item := &oqlGroup{Query: query, Args: args}
	m.groups = append(m.groups, item)
	return m
}

//fieldA
//fieldA desc,fieldB
func (m *exector) Order(query string, args ...interface{}) IExector {
	item := &oqlOrder{Query: query, Args: args}
	m.orders = append(m.orders, item)
	return m
}

// tableA
// tableA a
// tableA as a
// tableA,tableB
func (m *exector) From(query string) IExector {
	items := strings.Split(strings.TrimSpace(query), ",")
	for _, item := range items {
		parts := strings.Split(strings.TrimSpace(item), " ")
		item := &oqlFrom{Query: item}
		if len(parts) > 1 {
			item.Alia = parts[len(parts)-1]
		}
		m.froms = append(m.froms, item)
	}
	return m
}
