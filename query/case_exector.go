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
	Run() (*gorm.DB, error)
	PrepareQuery(mysql *repositories.MysqlRepo) (*gorm.DB, error)
	Select(query string, args ...interface{}) IExector
	Where(query string, args ...interface{}) IExector
	OrWhere(query string, args ...interface{}) IExector
	Joins(query string, args ...interface{}) IExector
	Group(query string, args ...interface{}) IExector
}
type oqlEntity struct {
	Alia     string
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
	Name  string
	Expr  string
	Args  []interface{}
}
type oqlOrder struct {
	Query    string
	IsDesc   bool
	Sequence int
	Expr     string
	Args     []interface{}
}
type oqlGroup struct {
	Query string
	Expr  string
	Args  []interface{}
}

type oqlWhere struct {
	Query    string
	Operator string
	Value    []interface{}
	Sequence int
	Children []oqlWhere
	Expr     string
	Args     []interface{}
}

type exector struct {
	entities map[string]*oqlEntity
	fields   map[string]*oqlField
	selects  []oqlSelect
	from     *oqlFrom
	joins    []oqlJoin
	wheres   []oqlWhere
	orders   []oqlOrder
	groups   []oqlGroup
}

func NewExector(query string) IExector {
	parts := strings.Split(strings.TrimSpace(query), " ")
	from := oqlFrom{Query: parts[0]}
	if len(parts) > 1 {
		from.Alia = parts[len(parts)-1]
	}
	exec := &exector{
		entities: make(map[string]*oqlEntity),
		fields:   make(map[string]*oqlField),
		from:     &from,
		selects:  make([]oqlSelect, 0),
		joins:    make([]oqlJoin, 0),
		wheres:   make([]oqlWhere, 0),
		orders:   make([]oqlOrder, 0),
		groups:   make([]oqlGroup, 0),
	}
	return exec
}
func (m *exector) Run() (*gorm.DB, error) {
	return nil, nil
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
	mainMd := md.GetEntity(m.from.Query)
	if mainMd != nil {
		m.entities[""] = m.FormatEntity(mainMd)
	}
	//parse
	if m.selects != nil && len(m.selects) > 0 {
		for _, v := range m.selects {
			if v.Name == "" {
				v.Name = strings.Replace(v.Query, ".", "_", -1)
			}
			m.parseField(v.Query)
		}
	}
	if m.orders != nil && len(m.orders) > 0 {
		for _, v := range m.orders {
			m.parseField(v.Query)
		}
	}
	if m.wheres != nil && len(m.wheres) > 0 {
		for _, v := range m.wheres {
			m.parseWhereField(v)
		}
	}

	//build
	mainTable := m.entities[""]
	queryDB := mysql.Table(fmt.Sprintf("%v as %v", mainTable.Entity.TableName, mainTable.Alia))
	queryDB = m.buildSelects(queryDB)
	queryDB = m.buildJoins(queryDB)
	return queryDB, nil
}

///=============== build
func (m *exector) buildSelects(queryDB *gorm.DB) *gorm.DB {
	selects := make([]string, 0)
	if m.selects != nil && len(m.selects) > 0 {
		for _, v := range m.selects {
			field := m.parseField(v.Query)
			if field != nil {
				selects = append(selects, fmt.Sprintf("%v.%v as %v", field.Entity.Alia, field.Field.DBName, v.Name))
			}
		}
	}
	queryDB.Select(selects)
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
func (m *exector) parseWhereField(where oqlWhere) {
	if where.Children != nil && len(where.Children) > 0 {
		for _, v := range where.Children {
			m.parseWhereField(v)
		}
	} else {
		m.parseField(where.Query)
	}
}

// 解析实体
func (m *exector) parseEntity(id, path string) *oqlEntity {
	path = strings.ToLower(path)
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
	fieldPath = strings.ToLower(fieldPath)
	if v, ok := m.fields[fieldPath]; ok {
		return v
	}
	parts := strings.Split(fieldPath, ".")
	entity := m.entities[""]
	path := ""
	for i, part := range parts {
		if i > 0 {
			path += "."
		}
		path += part
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
	item := oqlSelect{Query: query, Args: args}
	m.selects = append(m.selects, item)
	return m
}

// fieldA =?
// fieldB =''
// FieldC is null
func (m *exector) Where(query string, args ...interface{}) IExector {
	item := oqlWhere{Query: query, Args: args, Operator: "and"}
	m.wheres = append(m.wheres, item)
	return m
}
func (m *exector) OrWhere(query string, args ...interface{}) IExector {
	item := oqlWhere{Query: query, Args: args, Operator: "or"}
	m.wheres = append(m.wheres, item)
	return m
}

// left join tableA on a=b
// left join tableA as c on c.a=b.b
// inner join tableA
func (m *exector) Joins(query string, args ...interface{}) IExector {
	item := oqlJoin{Query: query, Args: args}
	m.joins = append(m.joins, item)
	return m
}

//fieldA
//fieldA,fieldB
func (m *exector) Group(query string, args ...interface{}) IExector {
	item := oqlGroup{Query: query, Args: args}
	m.groups = append(m.groups, item)
	return m
}
