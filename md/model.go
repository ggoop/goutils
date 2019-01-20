package md

import (
	"time"
)

type Model struct {
	ID        string     `gorm:"primary_key;size:100" json:"id"`
	CreatedAt Time       `gorm:"name:创建时间"  json:"created_at"`
	UpdatedAt Time       `gorm:"name:更新时间" json:"updated_at"`
	DeletedAt *time.Time `gorm:"name:删除时间" sql:"index" json:"deleted_at"`
}
type ModelUnscoped struct {
	ID        string `gorm:"primary_key;size:100" json:"id"`
	CreatedAt Time   `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt Time   `gorm:"name:更新时间" json:"updated_at"`
}
