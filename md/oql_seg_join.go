package md

type OQLJoin interface {
	GetQuery() string
	GetArgs() []interface{}
}
type oqlJoin struct {
	Origin interface{}
	Query  string
	Args   []interface{}
}

func (m oqlJoin) GetQuery() string {
	return m.Query
}
func (m oqlJoin) GetArgs() []interface{} {
	return m.Args
}
