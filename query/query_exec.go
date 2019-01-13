package query

import (
	"fmt"
	"strings"

	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/md"
	"github.com/ggoop/goutils/repositories"
	"github.com/jinzhu/gorm"
)

func Run(repo *repositories.MysqlRepo, item Query) (*gorm.DB, error) {
	if repo == nil || item.Code == "" {
		return nil, fmt.Errorf("参数不正确")
	}
	exec := &queryExec{repo: repo, query: &item}
	return exec.Run()
}

type queryExec struct {
	repo      *repositories.MysqlRepo
	query     *Query
	entities  map[string]*execEntity
	fields    map[string]*execEntityField
	mainEnity *execEntity
}

type execEntity struct {
	Table  string
	Alia   string
	IsMain bool
	Entity *md.MDEntity
}
type execEntityField struct {
	Table       *execEntity
	Alia        string
	DBFieldName string
	Path        string
}

func (m *queryExec) FormatQueryEntity(entity *md.MDEntity) *execEntity {
	e := execEntity{Table: entity.TableName, Entity: entity}
	return &e
}
func (m *queryExec) Run() (*gorm.DB, error) {
	entryEntity := &md.MDEntity{}
	m.repo.Where("code=?", m.query.Entry).Preload("Fields").Take(entryEntity)
	if entryEntity.ID == "" {
		return nil, fmt.Errorf("找不到实体 %v", m.query.Entry)
	}
	m.mainEnity = m.FormatQueryEntity(entryEntity)
	m.mainEnity.Alia = "a0"
	return nil, nil
}
func (m *queryExec) parseEntity(id, path string) *execEntity {
	path = strings.ToLower(path)
	if v, ok := m.entities[path]; ok {
		return v
	}
	entity := &md.MDEntity{}
	m.repo.Where("id=?", id).Preload("Fields").Take(entity)
	v := m.FormatQueryEntity(entity)
	if entity.ID == "" {
		glog.Errorf("找不到实体 %v", id)
		return nil
	}
	v.Alia = fmt.Sprintf("a%v", len(m.entities)+1)
	m.entities[path] = v
	return v
}
func (m *queryExec) parseField(fieldPath string) {
	parts := strings.Split(fieldPath, ".")
	entity := m.mainEnity
	path := ""
	for _, part := range parts {
		field := entity.Entity.GetField(part)
		if field == nil {
			break
		}
		entity = m.parseEntity(field.TypeID, path)
	}
}
