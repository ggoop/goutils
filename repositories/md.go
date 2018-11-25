package repositories

import (
	"github.com/jinzhu/gorm"
)

type mder interface {
	MDID() string
}

type MD struct {
	Value interface{}
	db    *gorm.DB
}

func (m *MD) MDID() string {
	if mder, ok := m.Value.(mder); ok {
		return mder.MDID()
	}
	return ""
}
func (m *MD) Migrate() {
	if m.MDID() == "" {
		return
	}
}
