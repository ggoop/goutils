package query

import (
	"github.com/ggoop/goutils/md"
	"github.com/ggoop/goutils/repositories"
	"github.com/jinzhu/gorm"
)

type IExector interface {
	Run() (*gorm.DB, error)
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
type oqlOrder struct {
	Field    string
	IsDesc   bool
	Sequence int
	Expr     string
	Args     []interface{}
}
type oqlGroup struct {
	Field string
	Expr  string
	Args  []interface{}
}
type oqlColumn struct {
	Field      string
	ColumnName string
	Expr       string
	Args       []interface{}
}
type oqlWhere struct {
	Field    string
	Operator string
	Value    []interface{}
	Sequence int
	Children []oqlWhere
	Expr     string
	Args     []interface{}
}

type exector struct {
	repo      *repositories.MysqlRepo
	mainEnity *oqlEntity
	entities  map[string]*oqlEntity
	fields    map[string]*oqlField
	orders    []oqlOrder
	columns   []oqlColumn
	wheres    []oqlWhere
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
