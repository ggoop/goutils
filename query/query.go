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

func (s *Query) MDID() string {
	return "01e8f0b45e12835fe7fd8cec4b7174de"
}

type QueryField struct {
	md.ModelUnscoped
	Code    string `gorm:"size:100;name:编码"`
	Name    string `gorm:"name:名称"`
	Memo    string `gorm:"name:备注"`
	Query   *Query `gorm:"name:查询"`
	QueryID string `gorm:"name:查询ID"`
}

func (s *QueryField) MDID() string {
	return "b6977390f18011e89e07b96946e7d763"
}
