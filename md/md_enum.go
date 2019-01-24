package md

import "time"

type MDEnumType struct {
	ID        string    `gorm:"primary_key;size:50" json:"id"`
	CreatedAt time.Time `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"name:更新时间" json:"updated_at"`
	ScopeID   string    `gorm:"size:50;name:范围ID" json:"scope_id"`
	ScopeType string    `gorm:"size:50;name:范围类型" json:"scope_type"`
	Name      string
	Memo      string
	IsSystem  bool
}

func (s *MDEnumType) MD() *Mder {
	return &Mder{ID: "01e9125fe960a71bb1b47427ea1d5200", Name: "枚举"}
}

type MDEnum struct {
	ID        string    `gorm:"primary_key;size:50" json:"id"`
	CreatedAt time.Time `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"name:更新时间" json:"updated_at"`
	ScopeID   string    `gorm:"index:scope;size:50;name:范围ID" json:"scope_id"`
	ScopeType string    `gorm:"index:scope;size:50;name:范围类型" json:"scope_type"`
	Type      string    `gorm:"primary_key;size:50"`
	Name      string
	Memo      string
	Sequence  int
	IsSystem  bool
}

func (s *MDEnum) MD() *Mder {
	return &Mder{ID: "01e9125fe9611c4dc8d47427ea1d5200", Name: "枚举值"}
}
