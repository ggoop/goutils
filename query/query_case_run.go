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

func RunCase(repo *repositories.MysqlRepo, item QueryCase) (*gorm.DB, error) {
	if repo == nil || item.Query == nil {
		return nil, fmt.Errorf("参数不正确")
	}
	exec := &caseExector{repo: repo, qcase: &item, entities: make(map[string]*execEntity), fields: make(map[string]*execField)}
	return exec.Run()
}

type caseExector struct {
	repo      *repositories.MysqlRepo
	qcase     *QueryCase
	entities  map[string]*execEntity
	fields    map[string]*execField
	mainEnity *execEntity
}

type execEntity struct {
	Table    string
	Alia     string
	IsMain   bool
	Entity   *md.MDEntity
	Path     string
	Sequence int
}
type execField struct {
	Entity *execEntity
	Field  *md.MDField
	Path   string
}

func (m *caseExector) FormatQueryEntity(entity *md.MDEntity) *execEntity {
	e := execEntity{Table: entity.TableName, Entity: entity}
	return &e
}
func (m *caseExector) FormatQueryField(entity *execEntity, field *md.MDField) *execField {
	e := execField{Entity: entity, Field: field}
	return &e
}
func (m *caseExector) Run() (*gorm.DB, error) {
	entryEntity := md.GetEntity(m.qcase.Query.Entry)
	if entryEntity.ID == "" {
		return nil, fmt.Errorf("找不到实体 %v", m.qcase.Query.Entry)
	}
	m.mainEnity = m.FormatQueryEntity(entryEntity)
	m.mainEnity.Alia = "a0"
	m.mainEnity.IsMain = true
	m.mainEnity.Sequence = 0
	//parse
	if m.qcase.Columns != nil && len(m.qcase.Columns) > 0 {
		for _, v := range m.qcase.Columns {
			if v.ColumnName == "" {
				v.ColumnName = strings.Replace(v.Field, ".", "_", -1)
			}
			m.parseField(v.Field)
		}
	}
	if m.qcase.Orders != nil && len(m.qcase.Orders) > 0 {
		for _, v := range m.qcase.Orders {
			m.parseField(v.Field)
		}
	}
	if m.qcase.Wheres != nil && len(m.qcase.Wheres) > 0 {
		for _, v := range m.qcase.Wheres {
			m.parseWhereField(v)
		}
	}
	//build
	queryDB := m.repo.Table(fmt.Sprintf("%v as %v", m.mainEnity.Table, m.mainEnity.Alia))
	queryDB = m.buildColumns(queryDB)
	queryDB = m.buildJoins(queryDB)
	queryDB = m.buildWheres(queryDB)

	count := 0
	queryDB.Count(&count)
	return queryDB, nil
}
func (m *caseExector) buildColumns(queryDB *gorm.DB) *gorm.DB {
	selects := make([]string, 0)
	if m.qcase.Columns != nil && len(m.qcase.Columns) > 0 {
		for _, v := range m.qcase.Columns {
			field := m.parseField(v.Field)
			if field != nil {
				selects = append(selects, fmt.Sprintf("%v.%v as %v", field.Entity.Alia, field.Field.DBName, v.ColumnName))
			}
		}
	}
	queryDB.Select(selects)
	return queryDB
}
func (m *caseExector) buildJoins(queryDB *gorm.DB) *gorm.DB {
	tables := make([]*execEntity, 0)
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
func (m *caseExector) buildWheres(queryDB *gorm.DB) *gorm.DB {
	if m.qcase.Wheres != nil && len(m.qcase.Wheres) > 0 {
		for _, v := range m.qcase.Wheres {
			queryDB = m.buildWhereItem(queryDB, v)
		}
	}
	return queryDB
}
func (m *caseExector) buildWhereItem(queryDB *gorm.DB, item QueryWhere) *gorm.DB {

	return queryDB
}
func (m *caseExector) parseWhereField(where QueryWhere) {
	if where.Children != nil && len(where.Children) > 0 {
		for _, v := range where.Children {
			m.parseWhereField(v)
		}
	} else {
		m.parseField(where.Field)
	}
}

// 解析实体
func (m *caseExector) parseEntity(id, path string) *execEntity {
	path = strings.ToLower(path)
	if v, ok := m.entities[path]; ok {
		return v
	}
	entity := md.GetEntity(id)
	if entity == nil {
		glog.Errorf("找不到实体 %v", id)
		return nil
	}
	v := m.FormatQueryEntity(entity)
	v.Sequence = len(m.entities) + 1
	v.Alia = fmt.Sprintf("a%v", v.Sequence)
	v.Path = path
	m.entities[path] = v
	return v
}

// 解析字段
func (m *caseExector) parseField(fieldPath string) *execField {
	fieldPath = strings.ToLower(fieldPath)
	if v, ok := m.fields[fieldPath]; ok {
		return v
	}
	parts := strings.Split(fieldPath, ".")
	entity := m.mainEnity
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
		field := m.FormatQueryField(entity, mdField)
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
