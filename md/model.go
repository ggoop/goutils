package md

import (
	"time"
)

type Model struct {
	ID        string     `gorm:"primary_key;size:100"`
	CreatedAt time.Time  `gorm:"name:创建时间"`
	UpdatedAt time.Time  `gorm:"name:更新时间"`
	DeletedAt *time.Time `gorm:"name:删除时间" sql:"index"`
}
type ModelUnscoped struct {
	ID        string    `gorm:"primary_key;size:100"`
	CreatedAt time.Time `gorm:"name:创建时间"`
	UpdatedAt time.Time `gorm:"name:更新时间"`
}
