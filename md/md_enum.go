package md

type MDEnumType struct {
	Model
	EntID    string `gorm:"size:100"`
	Name     string
	Memo     string
	IsSystem bool
}

func (s *MDEnumType) MD() *Mder {
	return &Mder{ID: "01e9125fe960a71bb1b47427ea1d5200", Name: "分类组"}
}

type MDEnum struct {
	Model
	EntID    string `gorm:"size:100"`
	GroupID  string `gorm:"size:100"`
	Code     string
	Name     string
	Memo     string
	Sequence int
	IsSystem bool
}

func (s *MDEnum) MD() *Mder {
	return &Mder{ID: "01e9125fe9611c4dc8d47427ea1d5200", Name: "分类"}
}
