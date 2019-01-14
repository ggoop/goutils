package query

import (
	"github.com/ggoop/goutils/md"
	"github.com/ggoop/goutils/repositories"
	"github.com/jinzhu/gorm"
)

type IExector interface {
	Run() (*gorm.DB, error)
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
type execOrder struct {
	Field    string
	IsDesc   bool
	Sequence int
}
type execColumn struct {
	Field      string
	ColumnName string
	Title      string
}
type execWhere struct {
	Field    string
	Operator string
	Value    []interface{}
	Sequence int
	Children []execWhere
}

type exector struct {
	repo      *repositories.MysqlRepo
	mainEnity *execEntity
	entities  map[string]*execEntity
	fields    map[string]*execField
	orders    []execOrder
	columns   []execColumn
	wheres    []execWhere
}

func NewExector(repo *repositories.MysqlRepo, md md.MDEntity) IExector {
	main := execEntity{Entity: &md, Table: md.TableName, IsMain: true, Alia: "a0"}
	exec := &exector{
		repo:      repo,
		mainEnity: &main,
		entities:  make(map[string]*execEntity),
		fields:    make(map[string]*execField),
		orders:    make([]execOrder, 0),
		columns:   make([]execColumn, 0),
		wheres:    make([]execWhere, 0),
	}
	return exec
}
func (m *exector) Run() (*gorm.DB, error) {
	return nil, nil
}
func (m *exector) Where(field string, args ...interface{}) IExector {
	return m
}
