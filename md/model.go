package md

import (
	"time"
)

type Model struct {
	ID        string `gorm:"primary_key;size:100"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
type ModelUnscoped struct {
	ID        string `gorm:"primary_key;size:100"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
