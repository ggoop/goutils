package md

type OQLField interface {
	GetQuery() string
	GetAlias() string
	GetArgs() []interface{}
}
type oqlField struct {
	Origin interface{}
	//表达式，如：sum(fieldA) as a
	Query string
	Alias string
	Expr  string
	Args  []interface{}
}

func (m oqlField) GetQuery() string {
	return m.Query
}
func (m oqlField) GetAlias() string {
	return m.Alias
}
func (m oqlField) GetArgs() []interface{} {
	return m.Args
}
