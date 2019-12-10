package md

type Model struct {
	ID        string `gorm:"primary_key;size:50" json:"id"`
	CreatedAt Time   `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt Time   `gorm:"name:更新时间" json:"updated_at"`
}
