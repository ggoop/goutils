package md

type OQLFrom interface {
	GetQuery() string
	GetAlias() string
	GetArgs() []interface{}
}
type oqlFrom struct {
	Origin interface{}
	Query  string
	Alias  string
	Args   []interface{}
}

func (m oqlFrom) GetQuery() string {
	return m.Query
}
func (m oqlFrom) GetAlias() string {
	return m.Alias
}
func (m oqlFrom) GetArgs() []interface{} {
	return m.Args
}
