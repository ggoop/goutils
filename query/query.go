package query

import (
	"github.com/ggoop/goutils/md"
)

type Query struct {
	md.ModelUnscoped
	Code string `gorm:"size:100;name:编码"`
	Name string `gorm:"name:名称"`
	Memo string `gorm:"name:备注"`

	Fields []QueryField `gorm:"name:字段集合"`
}

func (s *Query) MD() *md.Mder {
	return &md.Mder{ID: "01e8f3067a691a50b46b697fa9f73d01", Name: "查询"}
}

type QueryField struct {
	md.ModelUnscoped
	Code    string `gorm:"size:100;name:编码"`
	Name    string `gorm:"name:名称"`
	Memo    string `gorm:"name:备注"`
	Query   *Query `gorm:"name:查询"`
	QueryID string `gorm:"name:查询ID"`
}

func (s *QueryField) MD() *md.Mder {
	return &md.Mder{ID: "01e8f30683641760aa9261a2b248c5f0", Name: "查询字段"}
}
